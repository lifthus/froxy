package service

import (
	"encoding/json"
	"net/http"

	"github.com/lifthus/froxy/internal/dashboard/httphelper"
	"github.com/lifthus/froxy/internal/dashboard/root"
)

func GetSessionInfo(w http.ResponseWriter, r *http.Request) {
	cinfo := httphelper.ClientInfo(r)

	cinfob, err := json.Marshal(cinfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(cinfob)
}

func RootSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uname := r.PostForm.Get("username")
	pw := r.PostForm.Get("password")

	ok := root.Validate(uname, pw)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cinfo := httphelper.ClientInfo(r)
	cinfo.Root = true
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	cinfo := httphelper.ClientInfo(r)
	cinfo.Root = false
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
