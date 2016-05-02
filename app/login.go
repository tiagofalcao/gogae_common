package app

import (
	"appengine"
	"appengine/user"

	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, "/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
	} else {
		target := r.FormValue("target")
		if target == "" {
			target = "/"
		}
		w.Header().Set("Location", target)
	}
	w.WriteHeader(http.StatusFound)
	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		url, err := user.LogoutURL(c, "/")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
	} else {
		w.Header().Set("Location", "/")
	}
	w.WriteHeader(http.StatusFound)
	return
}
