package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/CHIRANTAN-001/lrucache/pkg/lrucache"
)

// CacheStats tracks cache hits and misses using atomic operations for thread safety.
type CacheStats struct {
	hits   int64
	misses int64
}

func (cs *CacheStats) RecordHit() {
	atomic.AddInt64(&cs.hits, 1)
}

func (cs *CacheStats) RecordMiss() {
	atomic.AddInt64(&cs.misses, 1)
}

func (cs *CacheStats) GetStats() (hits, misses int64, hitRate float64) {
	hits = atomic.LoadInt64(&cs.hits)
	misses = atomic.LoadInt64(&cs.misses)
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}
	return hits, misses, hitRate
}

func (cs *CacheStats) Reset() {
	atomic.StoreInt64(&cs.hits, 0)
	atomic.StoreInt64(&cs.misses, 0)
}

var stats = &CacheStats{}

// getProductDetailsFromAPI fetches product details from an external API.
func getProductDetailsFromAPI(id int) (string, error) {
	url := fmt.Sprintf("https://dummyjson.com/products/%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch product details: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch product details: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

// getProduct retrieves a product from the cache or API, updating global stats.
func getProduct(id int, cache *lrucache.LRUCache) (string, error) {
	key := fmt.Sprintf("product_%d", id)

	if value, ok := cache.Get(key); ok {
		stats.RecordHit()
		return value, nil
	}

	stats.RecordMiss()

	product, err := getProductDetailsFromAPI(id)
	if err != nil {
		return "", err
	}

	cache.Put(key, product)
	return product, nil
}

// benchmarkCacheHit simulates concurrent users requesting products and returns benchmark stats.
func benchmarkCacheHit(cache *lrucache.LRUCache, users, productRange int) (int64, int64, float64) {
	localStats := &CacheStats{} // Local stats for this benchmark run
	var wg sync.WaitGroup

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := r.Intn(productRange) + 1
			_, _ = getProductWithStats(id, cache, localStats)
		}()
	}

	wg.Wait()

	return localStats.GetStats()
}

// getProductWithStats is used by the benchmark to track hits/misses in local stats.
func getProductWithStats(id int, cache *lrucache.LRUCache, stats *CacheStats) (string, error) {
	key := fmt.Sprintf("product_%d", id)

	if value, ok := cache.Get(key); ok {
		stats.RecordHit()
		return value, nil
	}

	stats.RecordMiss()

	product, err := getProductDetailsFromAPI(id)
	if err != nil {
		return "", err
	}

	cache.Put(key, product)
	return product, nil
}

func main() {
	cache, err := lrucache.NewLRUCache(5)
	if err != nil {
		log.Fatal("Failed to create LRUCache:", err)
	}

	config := fiber.Config{
		Prefork: true,
	}

	app := fiber.New(config)

	// Hello World endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Product details endpoint
	app.Get("/product/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid product ID",
			})
		}

		if id < 1 || id > 100 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Product ID must be between 1 and 100",
			})
		}

		product, err := getProduct(id, cache)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var productDetails map[string]interface{}
		err = json.Unmarshal([]byte(product), &productDetails)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": productDetails,
		})
	})

	// Benchmark endpoint: http://localhost:8080/stats?users=20&range=3
	app.Get("/stats", func(c *fiber.Ctx) error {
		users, err := strconv.Atoi(c.Query("users", "20"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid users parameter"})
		}
		productRange, err := strconv.Atoi(c.Query("range", "3"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid range parameter"})
		}

		hits, misses, hitRate := benchmarkCacheHit(cache, users, productRange)
		return c.JSON(fiber.Map{
			"hits":     hits,
			"misses":   misses,
			"hit_rate": fmt.Sprintf("%.2f", hitRate),
			"total":    hits + misses,
		})
	})

	err = app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
