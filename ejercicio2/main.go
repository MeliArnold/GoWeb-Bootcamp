package main

import (
	"API2/ejercicio2/modelo/handler"
	"fmt"
	"net/http"
)

func main() {

	//  probando greetings
	http.HandleFunc("/greetings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		handler.Greetings(w, r)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
