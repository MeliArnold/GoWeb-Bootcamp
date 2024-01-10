package handler

import (
	"API2/ejercicio2/modelo"
	"encoding/json"
	"net/http"
)

func Greetings(response http.ResponseWriter, request *http.Request) {

	var persona modelo.Person

	// Intenta decodificar el cuerpo de la solicitud JSON en la variable 'categoria'
	if err := json.NewDecoder(request.Body).Decode(&persona); err != nil {
		respuesta := map[string]string{
			"estado":  "error",
			"mensaje": "Ocurrio un error inesperado!",
		}
		json.NewEncoder(response).Encode(respuesta)
		return
	}

	datos := modelo.Person{FirstName: persona.FirstName, LastName: persona.LastName}
	// Codifica la respuesta de éxito en formato JSON y la envía al cliente
	json.NewEncoder(response).Encode(datos)
}
