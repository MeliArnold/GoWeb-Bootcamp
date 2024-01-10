package main

import "net/http"

func Handler1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
