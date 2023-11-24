package dashboard

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/root"
	"github.com/lifthus/froxy/internal/dashboard/session"
)

func loadSessionMux() error {
	HandleGET("/api/session", func(w http.ResponseWriter, r *http.Request) {
		getSessionInfo(w, r)
	})
	HandlePOST("/api/session/root", func(w http.ResponseWriter, r *http.Request) {
		rootSignIn(w, r)
	})
	return nil
}

func getSessionInfo(w http.ResponseWriter, r *http.Request) {
	cinfo := &session.ClientInfo{
		IPAddr: getIPAddr(r),
	}

	cinfob, err := json.Marshal(cinfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cinfob)
}

func getIPAddr(r *http.Request) string {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func rootSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uname := r.PostForm.Get("username")
	pw := r.PostForm.Get("password")

	ok := root.Validate(uname, pw)

	// TODO: Sign Root

	fmt.Println(uname, pw, ok)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
