package handlers

import (
	"net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("metrics\n"))
}

