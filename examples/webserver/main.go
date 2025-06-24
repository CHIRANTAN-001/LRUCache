package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/CHIRANTAN-001/lrucache/pkg/lrucache"
)

func main() {

	cache, err := lrucache.NewLRUCache(5)
	if err != nil {
		fmt.Printf("Error creating cache: %v\n", err)
		return
	}

	// Define handler for /cache endpoint
	http.HandleFunc("/cache/{id}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("id")
		if key == "" {
			http.Error(w, "Key parameter is required", http.StatusBadRequest)
			return
		}

		if value, ok := cache.Get(key); ok {
			fmt.Printf("Cache hit for key: %s\n", key)
			w.Write([]byte(value))
			return
		}

		url := "https://dummyjson.com/products/" + key
		res, err := http.Get(url)
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

		cache.Put(key, string(body))
		w.Write(body)
	})

	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
