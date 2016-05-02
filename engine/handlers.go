package engine

import (
	"appengine"
	"appengine/user"

	"net/http"
	"regexp"
)

func Handle(fn func(http.ResponseWriter, *http.Request, appengine.Context, *user.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)

		fn(w, r, c, u)
	}
}

func HandleUserReq(fn func(http.ResponseWriter, *http.Request, appengine.Context, *user.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
		fn(w, r, c, u)
	}
}

var re_mail = regexp.MustCompile("^([^@]*)@([^@]*)$")

func HandleDomainReq(fn func(http.ResponseWriter, *http.Request, appengine.Context, *user.User), domain string) http.HandlerFunc {
	return HandleUserReq(func(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
		var m = re_mail.FindStringSubmatch(u.Email)
		if m == nil || m[2] != domain {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fn(w, r, c, u)
	})
}

func HandleAdminReq(fn func(http.ResponseWriter, *http.Request, appengine.Context, *user.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)
		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
		if !user.IsAdmin(c) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		fn(w, r, c, u)
	}
}
