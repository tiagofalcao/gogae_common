package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
}
