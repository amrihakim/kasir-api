package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func extractIDFromPath(r *http.Request, prefix string) (int, error) {
    idStr := strings.TrimPrefix(r.URL.Path, prefix)
    idStr = strings.Trim(idStr, "/")               // hilangkan trailing slash jika ada
    if idStr == "" {
        return 0, fmt.Errorf("missing id in path")
    }

	if idx := strings.Index(idStr, "/"); idx != -1 {
        idStr = idStr[:idx]
    }
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, fmt.Errorf("invalid id: %w", err)
    }
    return id, nil
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/products/")
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func updateProductById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/products/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedProduct Product

	err = json.NewDecoder(r.Body).Decode(&updatedProduct)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if 	product.ID == id {
			updatedProduct.ID = product.ID
			products[i] = updatedProduct
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduct)
			return
		}
	}
	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func deleteProductById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/products/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "product deleted successfully",
			})
			return
		}
	}
	http.Error(w, "Product Not Found", http.StatusNotFound)
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/categories/")
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func updateCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/categories/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedCategory Category

	err = json.NewDecoder(r.Body).Decode(&updatedCategory)

	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, category := range categories {
		if 	category.ID == id {
			updatedCategory.ID = category.ID
			categories[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r, "/categories/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, category := range categories {
		if category.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "category deleted successfully",
			})
			return
		}
	}
	http.Error(w, "Category Not Found", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getProductById(w, r)
		} else if r.Method == http.MethodPut {
			updateProductById(w, r)
		} else if r.Method == http.MethodDelete {
			deleteProductById(w, r)
		}
	})


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
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getCategoryById(w, r)
		} else if r.Method == http.MethodPut {
			updateCategoryById(w, r)
		} else if r.Method == http.MethodDelete {
			deleteCategoryById(w, r)
		}
	})

	
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running in localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}