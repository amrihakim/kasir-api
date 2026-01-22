package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    int `json:"id"`
	Name  string `json:"name"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 1500, Stock: 10},
	{ID: 2, Name: "Smartphone", Price: 800, Stock: 25},
	{ID: 3, Name: "Tablet", Price: 600, Stock: 15},
}

var categories = []Category{
	{ID: 1, Name: "Electronics", Description: "Devices and gadgets"},
	{ID: 2, Name: "Home Appliances", Description: "Appliances for home use"},
	{ID: 3, Name: "Books", Description: "Various genres of books"},
}

func main() {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)
		} else if r.Method == http.MethodPost {
			var newProduct Product

			err := json.NewDecoder(r.Body).Decode(&newProduct)

			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			}

			newProduct.ID = len(products) + 1
			products = append(products, newProduct)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newProduct)

		}
	})

	

	// endpoint categories
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == http.MethodPost {
			var newCategory Category

			err := json.NewDecoder(r.Body).Decode(&newCategory)

			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			}

			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)

		}
	})



	//  endpoint root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}