package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/CHIRANTAN-001/lrucache/pkg/lrucache"
	"github.com/gorilla/mux"
)

func main() {
	// Create a new LRU cache with a capacity of 5 items
	cache, err := lrucache.NewLRUCache(5)
	if err != nil {
		fmt.Printf("Error creating cache: %v\n", err)
		return
	}

	// Create a new router
	router := mux.NewRouter()

	// HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Define handler for /cache/{id}
	router.HandleFunc("/cache/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]
		if key == "" {
			http.Error(w, "Key parameter is required", http.StatusBadRequest)
			return
		}

		// Check the cache first
		if value, ok := cache.Get(key); ok {
			fmt.Printf("Cache hit for key: %s\n", key)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(value))
			return
		}

		fmt.Printf("Cache miss for key: %s\n", key)

		url := "https://dummyjson.com/products/" + key
		res, err := client.Get(url)
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, "Failed to read response body", http.StatusInternalServerError)
			return
		}

		// Put the data in the cache
		cache.Put(key, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}).Methods("GET")

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
