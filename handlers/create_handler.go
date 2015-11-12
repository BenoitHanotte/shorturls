package handlers

import (
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create\n"))
}

