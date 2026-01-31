package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"log"

	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

// var products = []models.Product{
// 	{ID: 1, Name: "Laptop", Price: 1500, Stock: 10, CategoryID: 1},
// 	{ID: 2, Name: "Smartphone", Price: 800, Stock: 25, CategoryID: 1},
// 	{ID: 3, Name: "Tablet", Price: 600, Stock: 15, CategoryID: 1},
// }

// var categories = []models.Category{
// 	{ID: 1, Name: "Electronics", Description: "Devices and gadgets"},
// 	{ID: 2, Name: "Home Appliances", Description: "Appliances for home use"},
// 	{ID: 3, Name: "Books", Description: "Various genres of books"},
// }

// func updateProductById(w http.ResponseWriter, r *http.Request) {
// 	id, err := extractIDFromPath(r, "/products/")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var updatedProduct Product

// 	err = json.NewDecoder(r.Body).Decode(&updatedProduct)

// 	if err != nil {
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	for i, product := range products {
// 		if 	product.ID == id {
// 			updatedProduct.ID = product.ID
// 			products[i] = updatedProduct
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(updatedProduct)
// 			return
// 		}
// 	}
// 	http.Error(w, "Product Not Found", http.StatusNotFound)
// }

// func deleteProductById(w http.ResponseWriter, r *http.Request) {
// 	id, err := extractIDFromPath(r, "/products/")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	for i, product := range products {
// 		if product.ID == id {
// 			products = append(products[:i], products[i+1:]...)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "product deleted successfully",
// 			})
// 			return
// 		}
// 	}
// 	http.Error(w, "Product Not Found", http.StatusNotFound)
// }

// func getCategoryById(w http.ResponseWriter, r *http.Request) {
// 	id, err := extractIDFromPath(r, "/categories/")

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	for _, category := range categories {
// 		if category.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(category)
// 			return
// 		}
// 	}
// 	http.Error(w, "Category Not Found", http.StatusNotFound)
// }

// func updateCategoryById(w http.ResponseWriter, r *http.Request) {
// 	id, err := extractIDFromPath(r, "/categories/")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var updatedCategory Category

// 	err = json.NewDecoder(r.Body).Decode(&updatedCategory)

// 	if err != nil {
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	for i, category := range categories {
// 		if 	category.ID == id {
// 			updatedCategory.ID = category.ID
// 			categories[i] = updatedCategory
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(updatedCategory)
// 			return
// 		}
// 	}
// 	http.Error(w, "Category Not Found", http.StatusNotFound)
// }

// func deleteCategoryById(w http.ResponseWriter, r *http.Request) {
// 	id, err := extractIDFromPath(r, "/categories/")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	for i, category := range categories {
// 		if category.ID == id {
// 			categories = append(categories[:i], categories[i+1:]...)
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "category deleted successfully",
// 			})
// 			return
// 		}
// 	}
// 	http.Error(w, "Category Not Found", http.StatusNotFound)
// }

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	// http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == http.MethodGet {
	// 		getProductById(w, r)
	// 	} else if r.Method == http.MethodPut {
	// 		updateProductById(w, r)
	// 	} else if r.Method == http.MethodDelete {
	// 		deleteProductById(w, r)
	// 	}
	// })

	// http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(products)
	// } else if r.Method == http.MethodPost {
	// 	var newProduct Product

	// 	err := json.NewDecoder(r.Body).Decode(&newProduct)

	// 	if err != nil {
	// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
	// 	}

	// 	newProduct.ID = len(products) + 1
	// 	products = append(products, newProduct)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusCreated)
	// 	json.NewEncoder(w).Encode(newProduct)

	// }
	// })

	// endpoint categories
	// http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == http.MethodGet {
	// 		getCategoryById(w, r)
	// 	} else if r.Method == http.MethodPut {
	// 		updateCategoryById(w, r)
	// 	} else if r.Method == http.MethodDelete {
	// 		deleteCategoryById(w, r)
	// 	}
	// })

	// http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method == http.MethodGet {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		json.NewEncoder(w).Encode(categories)
	// 	} else if r.Method == http.MethodPost {
	// 		var newCategory Category

	// 		err := json.NewDecoder(r.Body).Decode(&newCategory)

	// 		if err != nil {
	// 			http.Error(w, "Invalid Request", http.StatusBadRequest)
	// 		}

	// 		newCategory.ID = len(categories) + 1
	// 		categories = append(categories, newCategory)

	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusCreated)
	// 		json.NewEncoder(w).Encode(newCategory)

	// 	}
	// })

	//  endpoint root
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Failed to read config:", err)
		}
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Println("DB_CONN:", config.DBConn)

	if config.Port == "" {
		config.Port = "8080"
	}

	// DB INIT
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// DEPENDENCY INJECTION
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// ROUTES
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	http.HandleFunc("/products", productHandler.HandleProducts)
	http.HandleFunc("/products/", productHandler.HandleProductByID)

	http.HandleFunc("/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/categories/", categoryHandler.HandleCategoryByID)

	// SERVER
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running at", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
