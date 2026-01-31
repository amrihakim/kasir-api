package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"kasir-api/database"
	"log"

	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// =====================
	// ENV SETUP
	// =====================
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	log.Println("PORT =", config.Port)
	log.Println("DB_CONN set =", config.DBConn != "")

	// =====================
	// DB INIT (NON-FATAL)
	// =====================
	var (
		productHandler  *handlers.ProductHandler
		categoryHandler *handlers.CategoryHandler
	)

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Println("DB NOT READY:", err)
	} else {
		log.Println("DB READY")

		defer db.Close()

		productRepo := repositories.NewProductRepository(db)
		productService := services.NewProductService(productRepo)
		productHandler = handlers.NewProductHandler(productService)

		categoryRepo := repositories.NewCategoryRepository(db)
		categoryService := services.NewCategoryService(categoryRepo)
		categoryHandler = handlers.NewCategoryHandler(categoryService)
	}

	// =====================
	// ROUTES
	// =====================
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
	})

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if productHandler == nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}
		productHandler.HandleProducts(w, r)
	})

	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		if productHandler == nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}
		productHandler.HandleProductByID(w, r)
	})

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if categoryHandler == nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}
		categoryHandler.HandleCategories(w, r)
	})

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if categoryHandler == nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}
		categoryHandler.HandleCategoryByID(w, r)
	})

	// =====================
	// SERVER
	// =====================
	addr := ":" + config.Port
	log.Println("Server running at", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
