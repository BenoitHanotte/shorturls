package handlers

import (
	"net/http"
)

func RetrieveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("retrieve\n"))
}

