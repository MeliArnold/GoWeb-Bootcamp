package handlers

import (
	"API2/ejercicio3/modelos"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// "ChargeProducts lee el contenido de un archivo JSON con datos de producto y lo devuelve como un slice de modelos.Product.
// Abre el archivo productos.json, lee su contenido y lo interpreta en un slice de modelos.Product utilizando la función json.Unmarshal."
func ChargeProducts() []modelos.Product {
	jsonFile, err := os.Open("./productos.json")
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var products []modelos.Product

	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		log.Fatalf("unable to parse json: %v", err)
	}

	return products
}

// getProductByID busca un producto en base a su ID dentro de un slice de modelos.Product.
func getProductByID(products []modelos.Product, id int) (modelos.Product, error) {
	for _, p := range products {
		if p.Id == id {
			return p, nil
		}
	}
	return modelos.Product{}, fmt.Errorf("product with ID %d not found", id)
}

// ListByID es una función que maneja una solicitud HTTP GET para obtener un producto por ID.
func ListByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	products := ChargeProducts()
	// Get the parameter from the route using chi
	param := chi.URLParam(request, "id")
	if param == "" {
		fmt.Println("ID not provided")
		return
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		fmt.Printf("Error parsing ID: %v\n", err)
		return
	}

	product, err := getProductByID(products, id)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(response).Encode(product); err != nil {
		// handle error
		resp := map[string]string{
			"error": "An error occurred  2",
		}
		json.NewEncoder(response).Encode(resp)
	}
}

// getProductsByPrice filtra el slice de productos por precio y devuelve los productos cuyo precio sea mayor que el valor proporcionado.
// Recibe un slice de modelos.Product y un valor de tipo float32 representando el precio.
// Itera sobre cada producto en el slice y verifica si su precio es mayor que el valor proporcionado.
// Si el precio es mayor, agrega el producto al resultado.
// Si el resultado está vacío, devuelve un error indicando que no se encontraron productos con un precio mayor al valor proporcionado.
// De lo contrario, devuelve el resultado y nil como error.
func getProductsByPrice(products []modelos.Product, price float32) ([]modelos.Product, error) {
	var result []modelos.Product
	for _, p := range products {
		if float32(p.Price) > price {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no products found with price larger than %.2f", price)
	}
	return result, nil
}

// ListByPrice funcion que filtra la lista de productos por precio y devuelve los resultados como JSON.
func ListByPrice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	products := ChargeProducts()
	// Get the parameter from the route using chi
	param := chi.URLParam(request, "price")
	if param == "" {
		fmt.Println("Price not provided")
		return
	}

	price, err := strconv.ParseFloat(param, 32)
	if err != nil {
		fmt.Printf("Error parsing price: %v\n", err)
		return
	}

	filteredProducts, err := getProductsByPrice(products, float32(price))
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(response).Encode(filteredProducts); err != nil {
		// handle error
		resp := map[string]string{
			"error": "An error occurred 2",
		}
		json.NewEncoder(response).Encode(resp)
	}
}
