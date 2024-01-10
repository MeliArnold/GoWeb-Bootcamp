package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/ping", Handler1)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
		return
	}
}
