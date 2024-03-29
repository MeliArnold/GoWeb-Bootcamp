package main

import (
	"API2/ejercicio3/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	handlers.ChargeProducts()

	r := chi.NewRouter()

	r.Get("/parametros/{id}", handlers.ListByID)
	r.Get("/precio/{price}", handlers.ListByPrice)

	http.ListenAndServe(":8080", r)
}
