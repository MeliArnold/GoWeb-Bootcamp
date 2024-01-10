package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

func conParametro(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Get the parameter from the route using chi
	param := chi.URLParam(r, "id")

	json.NewEncoder(w).Encode("el parametro es: " + param)
	return

}
