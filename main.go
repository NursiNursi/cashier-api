package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/models"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn 	string `mapstructure:"DB_CONN"`
}

var products = []models.Product{
	{ID: 1, Name: "Laptop", Price: 1000, Stock: 10},
	{ID: 2, Name: "Smartphone", Price: 500, Stock: 20},
	{ID: 3, Name: "Tablet", Price: 300, Stock: 15},
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/product", productHandler.HandleProduct)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Starting server on :" + config.Port)

	err = http.ListenAndServe(":" + config.Port, nil)

	if err != nil {
		fmt.Println("Server failed to start:", err)
	}	
}