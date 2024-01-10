package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	// agrupa las rutas de productos
	r.Route("/productos", func(r chi.Router) {

		r.Get("/items", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, this is a GET with items"))
		})
		r.Get("/precios", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, this is a GET with precios!"))
		})
	})

	// agrupa las rutas de usuarios
	r.Route("/users", func(r chi.Router) {

		r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, this is a GET with posts"))
		})
		r.Get("/nosotros/{id}", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, this is a GET with nosotros!"))
		})

		r.Get("/parametros/{id}", conParametro)
	})

	http.ListenAndServe(":8080", r)
}
