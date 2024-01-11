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

func validateProduct(product modelos.Product) error {
	if len(product.Name) == 0 || len(product.Expiration) == 0 || product.Price <= 0 || len(product.CodeValue) == 0 || product.Quantity <= 0 {
		return fmt.Errorf("invalid product data")
	}

	var products []modelos.Product
	// Check for unique code_value
	for _, existingProduct := range products {
		if existingProduct.CodeValue == product.CodeValue {
			return fmt.Errorf("code_value %s already used", product.CodeValue)
		}
	}

	return nil
}

// modified addProduct function to use validateProduct with products array
func addProduct(newproduct modelos.Product) error {
	products := ChargeProducts()

	// Loop through all products and check if code_value already exists
	for _, product := range products {
		if product.CodeValue == newproduct.CodeValue {
			return fmt.Errorf("code value %s already exists", newproduct.CodeValue)
		}
	}

	// Assign new id
	if len(products) > 0 {
		newproduct.Id = products[len(products)-1].Id + 1
	}

	// Validation moved here
	err := validateProduct(newproduct)
	if err != nil {
		return err
	}

	products = append(products, newproduct)
	jsonFile, err := os.OpenFile("./productos.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer jsonFile.Close()

	byteValue, err := json.Marshal(products)
	if err != nil {
		log.Fatalf("unable to marshal products: %v", err)
	}
	_, err = jsonFile.Write(byteValue)
	if err != nil {
		log.Fatalf("unable to write to file: %v", err)
	}
	return nil
}

func AddProductHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/properties")
	if request.Method != http.MethodPost {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(request.Body)
	var newproduct modelos.Product
	err := decoder.Decode(&newproduct)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validateProduct(newproduct)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		//response.WriteHeader(http.StatusBadRequest)
		respuesta := map[string]string{
			"mensaje": "error de validacion",
		}
		json.NewEncoder(response).Encode(respuesta)

	}
	err = addProduct(newproduct)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(newproduct)
}

// Editar un producto en el slice

// EditarProducto es una función que se utiliza para editar un producto existente en la base de datos.
// Recibe una respuesta http.ResponseWriter y una solicitud http.Request como parámetros.
// Establece el encabezado de la respuesta a "application/properties".
// Si el método de solicitud no es PUT, establece el estado de la respuesta a Method Not Allowed (405) y retorna.
// Decodifica el cuerpo de la solicitud JSON utilizando json.NewDecoder y lo asigna a la variable updatedProduct de tipo modelos.Product.
// Si hay un error durante la decodificación, establece el estado de la respuesta a Bad Request (400) y retorna.
// Llama a la función validateProduct con updatedProduct para validar los datos del producto.
// Si hay un error de validación, establece el encabezado de la respuesta a "application/json".
// Crea una respuesta en formato JSON con un mensaje de error y lo envía en la respuesta.
// Llama a la función UpdateProduct pasando updatedProduct para actualizar el producto en la base de datos.
// Establece el estado de la respuesta a OK (200).
// Codifica updatedProduct en formato JSON y lo envía en la respuesta.
func EditarProducto(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/properties")
	if request.Method != http.MethodPut {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(request.Body)
	var updatedProduct modelos.Product
	err := decoder.Decode(&updatedProduct)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validateProduct(updatedProduct)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		respuesta := map[string]string{
			"mensaje": "error de validacion",
		}
		json.NewEncoder(response).Encode(respuesta)
		return
	}
	UpdateProduct(updatedProduct)
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(updatedProduct)
}

// UpdateProduct actualiza un producto existente en la lista de productos.
// Recibe un nuevo producto (newProduct) y busca un producto en la lista de productos cargada con la función ChargeProducts
// que tenga el mismo id que el nuevo producto.
// Si encuentra un producto con el mismo id, lo actualiza con el nuevo producto y guarda la lista actualizada utilizando la función SaveProducts.
// Si no encuentra un producto con el mismo id, devuelve un error indicando que el producto no fue encontrado.
func UpdateProduct(newProduct modelos.Product) error {
	var productList []modelos.Product = ChargeProducts()
	for index, product := range productList {
		if product.Id == newProduct.Id {
			productList[index] = newProduct

			return SaveProducts(productList)
		}
	}
	return fmt.Errorf("Producto no encontrado")
}

// SaveProducts toma un slice de modelos.Product y lo guarda en un archivo JSON llamado "productos.json".
// Convierte los datos en formato JSON utilizando la función json.Marshal y luego los guarda en el archivo utilizando ioutil.WriteFile.
// Devuelve un error si hubo algún problema al convertir o guardar los datos en el archivo.
func SaveProducts(products []modelos.Product) error {
	jsonData, err := json.Marshal(products)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./productos.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// fin editar producto en el slice

// PatchPrice actualiza el precio de un producto existente en función de los datos proporcionados en la solicitud HTTP PATCH.
//
// Establece el encabezado "Content-Type" de la respuesta a "application/json".
// Si el método de la solicitud no es PATCH, establece el código de estado HTTP en 405 (Method Not Allowed) y finaliza la ejecución.
//
// Decodifica los parámetros de la solicitud HTTP PATCH utilizando json.Decoder.
// Deserializa los datos decodificados en una variable "updatedProduct" de tipo modelos.Product.
// Si ocurre un error durante la decodificación, establece el código de estado HTTP en 400 (Bad Request) y finaliza la ejecución.
//
// Obtiene una lista de productos utilizando la función ChargeProducts().
// Recorre la lista de productos y actualiza el precio del producto con el ID proporcionado en "updatedProduct".
// Si encuentra el producto y actualiza su precio, guarda la lista actualizada de productos utilizando la función SaveProducts().
// Si ocurre un error al guardar los productos, establece el código de estado HTTP en 500 (Internal Server Error) y finaliza la ejecución.
// Establece el código de estado HTTP en 200 (OK), codifica "updatedProduct" en JSON y lo escribe en la respuesta.
// Finaliza la ejecución.
//
// Si no encuentra ningún producto con el ID proporcionado, establece el código de estado HTTP en 404 (Not Found) y finaliza la ejecución.
func PatchPrice(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	if request.Method != http.MethodPatch {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Decode Parameters
	decoder := json.NewDecoder(request.Body)
	var updatedProduct modelos.Product
	err := decoder.Decode(&updatedProduct)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	// find product with provided id and update its Price
	var productList []modelos.Product = ChargeProducts()
	for index, product := range productList {
		if product.Id == updatedProduct.Id {
			productList[index].Price = updatedProduct.Price
			err := SaveProducts(productList)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			response.WriteHeader(http.StatusOK)
			json.NewEncoder(response).Encode(updatedProduct)
			return
		}
	}
	response.WriteHeader(http.StatusNotFound)

}
