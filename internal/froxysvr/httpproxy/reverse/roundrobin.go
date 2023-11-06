package reverse

import (
	"context"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/lifthus/froxy/pkg/helper"
)

func useRoundRobinLoadBalanceHandler(ff *ReverseFroxy) *ReverseFroxy {
	hpm := ff.HostProxyMap
	ff.handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if ctx.Done() != nil {
			// CloseNotifier predates context.Context, and has been
			// entirely superseded by it. If the request contains
			// a Context that carries a cancellation signal, don't
			// bother spinning up a goroutine to watch the CloseNotify
			// channel (if any).
			//
			// If the request Context has a nil Done channel (which
			// means it is either context.Background, or a custom
			// Context implementation with no cancellation signal),
			// then consult the CloseNotifier if available.
		} else if cn, ok := w.(http.CloseNotifier); ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()
			notifyChan := cn.CloseNotify()
			go func() {
				select {
				case <-notifyChan:
					cancel()
				case <-ctx.Done():
				}
			}()
		}

		matcher, ok := hpm.MatchHost(req.Host)
		if !ok {
			http.Error(w, "host not found", http.StatusNotFound)
			return
		}

		outreq := req.Clone(ctx)
		if req.ContentLength == 0 {
			outreq.Body = nil // Issue 16036: https://github.com/golang/go/issues/16036
		}
		if outreq.Body != nil {
			// Reading from the request body after returning from a handler is not
			// allowed, and the RoundTrip goroutine that reads the Body can outlive
			// this handler. This can lead to a crash if the handler panics (see
			// Issue 46866). Although calling Close doesn't guarantee there isn't
			// any Read in flight after the handle returns, in practice it's safe to
			// read after closing it.
			defer outreq.Body.Close()
		}
		if outreq.Header == nil {
			outreq.Header = make(http.Header) // Issue 33142: https://github.com/golang/go/issues/33142
		}

		proxyTarget, path, ok := matcher.Match(req.URL.Path)
		if !ok {
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}
		targetURL := proxyTarget.NextTargetURL()
		if targetURL == nil {
			http.Error(w, "no base path matched", http.StatusNotFound)
			return
		}

		targetURL = targetURL.JoinPath(path)

		targetURL = rewriteURL(req.URL, outreq.URL)
		outreq.URL = targetURL
		if outreq.Form != nil {
			outreq.URL.RawQuery = cleanQueryParams(outreq.URL.RawQuery)
		}
		outreq.Close = false

		reqUpType := upgradeType(outreq.Header)
		if !IsPrint(reqUpType) {
			http.Error(w, "client tried switching to invalid protocol", http.StatusServiceUnavailable)
			return
		}
		removeHopByHopHeaders(outreq.Header)
	})
	return ff
}

// IsPrint returns whether s is ASCII and printable according to
// https://tools.ietf.org/html/rfc20#section-4.2.
func IsPrint(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' || s[i] > '~' {
			return false
		}
	}
	return true
}

func rewriteURL(reqURL *url.URL, target *url.URL) *url.URL {
	reqURL.Scheme = target.Scheme
	reqURL.Host = target.Host
	reqURL.Path, reqURL.RawPath = helper.JoinURLPath(target, reqURL)
	targetQuery := target.RawQuery
	if targetQuery == "" || reqURL.RawQuery == "" {
		reqURL.RawQuery = targetQuery + reqURL.RawQuery
	} else {
		reqURL.RawQuery = targetQuery + "&" + reqURL.RawQuery
	}
	return reqURL
}

func cleanQueryParams(s string) string {
	reencode := func(s string) string {
		v, _ := url.ParseQuery(s)
		return v.Encode()
	}
	for i := 0; i < len(s); {
		switch s[i] {
		case ';':
			return reencode(s)
		case '%':
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				return reencode(s)
			}
			i += 3
		default:
			i++
		}
	}
	return s
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func upgradeType(h http.Header) string {
	if !HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	return h.Get("Upgrade")
}

// HeaderValuesContainsToken reports whether any string in values
// contains the provided token, ASCII case-insensitively.
func HeaderValuesContainsToken(values []string, token string) bool {
	for _, v := range values {
		if headerValueContainsToken(v, token) {
			return true
		}
	}
	return false
}

// headerValueContainsToken reports whether v (assumed to be a
// 0#element, in the ABNF extension described in RFC 7230 section 7)
// contains token amongst its comma-separated tokens, ASCII
// case-insensitively.
func headerValueContainsToken(v string, token string) bool {
	for comma := strings.IndexByte(v, ','); comma != -1; comma = strings.IndexByte(v, ',') {
		if tokenEqual(trimOWS(v[:comma]), token) {
			return true
		}
		v = v[comma+1:]
	}
	return tokenEqual(trimOWS(v), token)
}

// trimOWS returns x with all optional whitespace removes from the
// beginning and end.
func trimOWS(x string) string {
	// TODO: consider using strings.Trim(x, " \t") instead,
	// if and when it's fast enough. See issue 10292.
	// But this ASCII-only code will probably always beat UTF-8
	// aware code.
	for len(x) > 0 && isOWS(x[0]) {
		x = x[1:]
	}
	for len(x) > 0 && isOWS(x[len(x)-1]) {
		x = x[:len(x)-1]
	}
	return x
}

// isOWS reports whether b is an optional whitespace byte, as defined
// by RFC 7230 section 3.2.3.
func isOWS(b byte) bool { return b == ' ' || b == '\t' }

// tokenEqual reports whether t1 and t2 are equal, ASCII case-insensitively.
func tokenEqual(t1, t2 string) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i, b := range t1 {
		if b >= utf8.RuneSelf {
			// No UTF-8 or non-ASCII allowed in tokens.
			return false
		}
		if lowerASCII(byte(b)) != lowerASCII(t2[i]) {
			return false
		}
	}
	return true
}

// lowerASCII returns the ASCII lowercase version of b.
func lowerASCII(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

// removeHopByHopHeaders removes hop-by-hop headers.
func removeHopByHopHeaders(h http.Header) {
	// RFC 7230, section 6.1: Remove headers listed in the "Connection" header.
	for _, f := range h["Connection"] {
		for _, sf := range strings.Split(f, ",") {
			if sf = textproto.TrimString(sf); sf != "" {
				h.Del(sf)
			}
		}
	}
	// RFC 2616, section 13.5.1: Remove a set of known hop-by-hop headers.
	// This behavior is superseded by the RFC 7230 Connection header, but
	// preserve it for backwards compatibility.
	for _, f := range hopHeaders {
		h.Del(f)
	}
}
