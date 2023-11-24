package httpreverse

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

func useRoundRobinLoadBalanceHandler(ff *ReverseFroxy) *ReverseFroxy {
	hpm := ff.HostProxyMap
	ff.handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		outreq := req.Clone(ctx)
		if req.ContentLength == 0 {
			// Issue 16036: https://github.com/golang/go/issues/16036
			// briefly, it is about retry mechanism of http.Transport.
			// When reused connection is broken, the body must be nil
			// for the request to be retried.
			outreq.Body = nil
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
			// Issue 33142: https://github.com/golang/go/issues/33142
			// Clone doesn't set header if it isn't set.
			outreq.Header = make(http.Header)
		}

		host, _, splErr := net.SplitHostPort(req.Host)
		matcher, ok := hpm.MatchHost(host)
		if !ok {
			http.Error(w, "host not found", http.StatusNotFound)
			return
		}
		proxyTarget, path, ok := matcher.Match(req.URL.Path)
		if !ok {
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}

		targetURL := proxyTarget.NextTargetURL(path)
		rewriteReqURLToTarget(outreq, targetURL)

		outreq.Host = targetURL.Host
		outreq.Header.Set("Host", targetURL.Host)

		if outreq.Form != nil {
			outreq.URL.RawQuery = cleanQueryParams(outreq.URL.RawQuery)
		}
		outreq.Close = false

		reqUpgType := upgradeType(outreq.Header)
		if !isASCIIPrintable(reqUpgType) {
			http.Error(w, "client tried switching to invalid protocol", http.StatusServiceUnavailable)
			return
		}

		removeHopByHopHeaders(outreq.Header)

		// Issue 21096: tell backend applications that care about trailer support
		// that we support trailers. (We do, but we don't go out of our way to
		// advertise that unless the incoming client request thought it was worth
		// mentioning.) Note that we look at req.Header, not outreq.Header, since
		// the latter has passed through removeHopByHopHeaders.
		if HeaderValuesContainsToken(req.Header["Te"], "trailers") {
			outreq.Header.Add("Trailer", "trailers")
		}

		// After stripping all the hop-by-hop connection headers above, add back any
		// necessary for protocol upgrades, such as for websockets.
		if reqUpgType != "" {
			outreq.Header.Set("Connection", "Upgrade")
			outreq.Header.Set("Upgrade", reqUpgType)
		}

		if splErr == nil {
			// If we aren't the first proxy retain prior
			// X-Forwarded-For information as a comma+space
			// separated list and fold multiple headers into one.
			prior, ok := outreq.Header["X-Forwarded-For"]
			omit := ok && prior == nil // Issue 38079: nil now means don't populate the header
			if len(prior) > 0 {
				host = strings.Join(prior, ", ") + ", " + host
			}
			if !omit {
				outreq.Header.Set("X-Forwarded-For", host)
			}
		}

		if _, ok := outreq.Header["User-Agent"]; !ok {
			// If the outbound request doesn't have a User-Agent header set,
			// don't send the default Go HTTP client User-Agent.
			outreq.Header.Set("User-Agent", "")
		}

		trace := &httptrace.ClientTrace{
			Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
				h := w.Header()
				copyHeader(h, http.Header(header))
				w.WriteHeader(code)

				// Clear headers, it's not automatically done by ResponseWriter.WriteHeader() for 1xx response
				for k := range h {
					delete(h, k)
				}
				return nil
			},
		}
		outreq = outreq.WithContext(httptrace.WithClientTrace(outreq.Context(), trace))

		transport := http.Transport{
			TLSClientConfig: &tls.Config{
				// Golang uses the OS certificate store.
				// By setting this to true, it accepts any certificate from the backend.
				InsecureSkipVerify: true,
			},
		}
		res, err := transport.RoundTrip(outreq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if res.StatusCode == http.StatusSwitchingProtocols {
			handleUpgradeResponse(w, outreq, res)
			return
		}

		removeHopByHopHeaders(res.Header)

		copyHeader(w.Header(), res.Header)

		// The "Trailer" header isn't included in the Transport's response,
		// at least for *http.Transport. Build it up from Trailer.
		announcedTrailers := len(res.Trailer)
		if announcedTrailers > 0 {
			trailerKeys := make([]string, 0, len(res.Trailer))
			for k := range res.Trailer {
				trailerKeys = append(trailerKeys, k)
			}
			w.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
		}

		w.WriteHeader(res.StatusCode)

		err = copyResponse(w, res.Body, flushInterval(res))
		if err != nil {
			defer res.Body.Close()
			// Since we're streaming the response, if we run into an error all we can do
			// is abort the request. Issue 23643: ReverseProxy should use ErrAbortHandler
			// on read error while copying body.
			if !shouldPanicOnCopyError(req) {
				log.Printf("suppressing panic for copyResponse error in test; copy error: %v", err)
				return
			}
			panic(http.ErrAbortHandler)
		}
		res.Body.Close() // close now, instead of defer, to populate res.Trailer

		if len(res.Trailer) > 0 {
			// Force chunking if we saw a response trailer.
			// This prevents net/http from calculating the length for short
			// bodies and adding a Content-Length.
			if fl, ok := w.(http.Flusher); ok {
				fl.Flush()
			}
		}

		if len(res.Trailer) == announcedTrailers {
			copyHeader(w.Header(), res.Trailer)
			return
		}

		for k, vv := range res.Trailer {
			k = http.TrailerPrefix + k
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
	})
	return ff
}

func rewriteReqURLToTarget(req *http.Request, target *url.URL) {
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host

	req.URL.Path = target.Path

	targetQuery := target.RawQuery
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// cleanQueryParams removes invalid query params and returns all valid query params.
func cleanQueryParams(s string) string {
	for i := 0; i < len(s); {
		switch s[i] {
		case ';':
			return prunInvalidQueryParams(s)
		case '%':
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				return prunInvalidQueryParams(s)
			}
			i += 3
		default:
			i++
		}
	}
	return s
}

func prunInvalidQueryParams(s string) string {
	v, _ := url.ParseQuery(s) // always returns all valid query params
	return v.Encode()
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

func lowerASCII(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

// isPrint returns whether s is ASCII and printable according to
// https://tools.ietf.org/html/rfc20#section-4.2.
func isASCIIPrintable(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' || s[i] > '~' {
			return false
		}
	}
	return true
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

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func handleUpgradeResponse(rw http.ResponseWriter, req *http.Request, res *http.Response) {
	reqUpgType := upgradeType(req.Header)
	resUpgType := upgradeType(res.Header)
	if !isASCIIPrintable(resUpgType) { // we already checked reqUpgType is ASCII from the caller ServeHTTP.
		http.Error(rw, "backend tried switching to invalid protocol", http.StatusBadGateway)
	}
	if !equalFold(reqUpgType, resUpgType) {
		http.Error(rw, "backend tried switching to different protocol", http.StatusBadGateway)
	}

	hj, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "hijacking not supported", http.StatusInternalServerError)
		return
	}
	backConn, ok := res.Body.(io.ReadWriteCloser)
	if !ok {
		http.Error(rw, "internal error: 101 switching protocols response with non-writable body", http.StatusInternalServerError)
		return
	}

	backConnCloseCh := make(chan bool)
	go func() {
		// Ensure that the cancellation of a request closes the backend.
		// https://golang.org/issue/35559
		select {
		case <-req.Context().Done():
		case <-backConnCloseCh:
		}
		backConn.Close()
	}()

	conn, brw, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "hijacking failed on protocol switch", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	copyHeader(rw.Header(), res.Header)

	res.Header = rw.Header()
	res.Body = nil // so res.Write only writes the headers; we have res.Body in backConn above.
	if err := res.Write(brw); err != nil {
		http.Error(rw, "response write failed", http.StatusInternalServerError)
		return
	}
	if err := brw.Flush(); err != nil {
		http.Error(rw, "response flush failed", http.StatusInternalServerError)
		return
	}
	errc := make(chan error, 1)
	spc := switchProtocolCopier{user: conn, backend: backConn}
	go spc.copyToBackend(errc)
	go spc.copyFromBackend(errc)
	<-errc
}

// switchProtocolCopier exists so goroutines proxying data back and
// forth have nice names in stacks.
type switchProtocolCopier struct {
	user, backend io.ReadWriter
}

func (c switchProtocolCopier) copyFromBackend(errc chan<- error) {
	_, err := io.Copy(c.user, c.backend)
	errc <- err
}

func (c switchProtocolCopier) copyToBackend(errc chan<- error) {
	_, err := io.Copy(c.backend, c.user)
	errc <- err
}

// EqualFold is strings.EqualFold, ASCII only. It reports whether s and t
// are equal, ASCII-case-insensitively.
func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

// lower returns the ASCII lowercase version of b.
func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

func flushInterval(res *http.Response) time.Duration {
	resCT := res.Header.Get("Content-Type")

	// For Server-Sent Events responses, flush immediately.
	// The MIME type is defined in https://www.w3.org/TR/eventsource/#text-event-stream
	if baseCT, _, _ := mime.ParseMediaType(resCT); baseCT == "text/event-stream" {
		return -1 // negative means immediately
	}

	// We might have the case of streaming for which Content-Length might be unset.
	if res.ContentLength == -1 {
		return -1
	}

	return 0
}

func copyResponse(dst io.Writer, src io.Reader, flushInterval time.Duration) error {
	if flushInterval != 0 {
		if wf, ok := dst.(writeFlusher); ok {
			mlw := &maxLatencyWriter{
				dst:     wf,
				latency: flushInterval,
			}
			defer mlw.stop()

			// set up initial timer so headers get flushed even if body writes are delayed
			mlw.flushPending = true
			mlw.t = time.AfterFunc(flushInterval, mlw.delayedFlush)

			dst = mlw
		}
	}

	var buf []byte
	_, err := copyBuffer(dst, src, buf)
	return err
}

type writeFlusher interface {
	io.Writer
	http.Flusher
}

type maxLatencyWriter struct {
	dst     writeFlusher
	latency time.Duration // non-zero; negative means to flush immediately

	mu           sync.Mutex // protects t, flushPending, and dst.Flush
	t            *time.Timer
	flushPending bool
}

func (m *maxLatencyWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, err = m.dst.Write(p)
	if m.latency < 0 {
		m.dst.Flush()
		return
	}
	if m.flushPending {
		return
	}
	if m.t == nil {
		m.t = time.AfterFunc(m.latency, m.delayedFlush)
	} else {
		m.t.Reset(m.latency)
	}
	m.flushPending = true
	return
}

func (m *maxLatencyWriter) delayedFlush() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.flushPending { // if stop was called but AfterFunc already started this goroutine
		return
	}
	m.dst.Flush()
	m.flushPending = false
}

func (m *maxLatencyWriter) stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushPending = false
	if m.t != nil {
		m.t.Stop()
	}
}

// copyBuffer returns any write errors or non-EOF read errors, and the amount
// of bytes written.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	if len(buf) == 0 {
		buf = make([]byte, 32*1024)
	}
	var written int64
	for {
		nr, rerr := src.Read(buf)
		if rerr != nil && rerr != io.EOF && rerr != context.Canceled {
			log.Printf("httputil: ReverseProxy read error during body copy: %v", rerr)
		}
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if werr != nil {
				return written, werr
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if rerr != nil {
			if rerr == io.EOF {
				rerr = nil
			}
			return written, rerr
		}
	}
}

// shouldPanicOnCopyError reports whether the reverse proxy should
// panic with http.ErrAbortHandler. This is the right thing to do by
// default, but Go 1.10 and earlier did not, so existing unit tests
// weren't expecting panics. Only panic in our own tests, or when
// running under the HTTP server.
func shouldPanicOnCopyError(req *http.Request) bool {
	// We seem to be running under an HTTP server, so
	// it'll recover the panic.
	return req.Context().Value(http.ServerContextKey) != nil
	// Otherwise act like Go 1.10 and earlier to not break
	// existing tests.
}
