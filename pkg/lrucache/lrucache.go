package lrucache

import "errors"

// LRUCache implements a Least Recently Used (LRU) cache.
// It uses a doubly linked list to maintain the order of usage and a map for O(1) access.
// The cache evicts the least recently used item when it exceeds its capacity.
// It provides methods to get and put items in the cache.
type Node struct {
	Key   string
	Value string
	Prev  *Node
	Next  *Node
}

type LRUCache struct {
	Capacity int
	Head     *Node
	Tail     *Node
	Cache    map[string]*Node
}

// NewLRUCache creates a new LRUCache Instance with the specified capacity.
func NewLRUCache(capacity int) (*LRUCache, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid capacity: must be greater than 0")
	}

	return &LRUCache{
		Capacity: capacity,
		Head:     nil,
		Tail:     nil,
		Cache:    make(map[string]*Node),
	} , nil
}

// Get retrieves the value for a given key from the cache.
func (c *LRUCache) Get(key string) (string, bool) {
	if node, ok := c.Cache[key]; ok {
		// Move the accessed node to the head of the list
		c.moveToHead(node)
		return node.Value, true
	}
	return "", false
}

func (c *LRUCache) moveToHead(node *Node) {
	if c.Head == node {
		return
	}

	// Detach the node from its current position
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if c.Tail == node {
		c.Tail = node.Prev
	}

	// Move the node to the head
	node.Next = c.Head
	node.Prev = nil
	if c.Head != nil {
		c.Head.Prev = node
	}
	c.Head = node

	// Adjust tail if necessary
	if c.Tail == nil {
		c.Tail = c.Head
	}
	if c.Tail != nil {
		c.Tail.Next = nil
	}

	// Update cache reference
	c.Cache[node.Key] = node
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves the node to the head.
func (c *LRUCache) Put(key string, value string) {
	// If the key already exists, update the value and move to head
	if node, ok := c.Cache[key]; ok {
		node.Value = value
		c.moveToHead(node)
	}

	// Create a new node
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	// Add the new node to the cache
	c.Cache[key] = newNode

	// If the cache is at capacity, remove the least recently used item
	if len(c.Cache) > c.Capacity {
		// Remove the least recently used item
		delete(c.Cache, c.Tail.Key)
		c.Tail = c.Tail.Prev
		if c.Tail != nil {
			c.Tail.Next = nil
		} else {
			c.Head = nil // If the list is now empty, reset head
		}
	}

	// Insert the new node at the head of the list
	if c.Head == nil {
		c.Head = newNode
		c.Tail = newNode
	} else {
		newNode.Next = c.Head
		c.Head.Prev = newNode
		c.Head = newNode
	}
}

// Clear removes all items from the cache.
func (c *LRUCache) Clear() {
	c.Head = nil
	c.Tail = nil
	c.Cache = make(map[string]*Node)
}

// Size returns the current number of items in the cache.
func (c *LRUCache) Size() int {
	return len(c.Cache)
}

// IsEmpty checks if the cache is empty.
func (c *LRUCache) IsEmpty() bool {
	return len(c.Cache) == 0
}

// Contains checks if the cache contains a specific key.
func (c *LRUCache) Contains(key string) bool {
	_, ok := c.Cache[key]
	return ok
}

func (c *LRUCache) BatchPut(items map[string]string) {
	for key, value := range items {
		c.Put(key, value)
	}
}

// BatchGet retrieves multiple values from the cache.
func (c *LRUCache) BatchGet(keys []string) (map[string]string, bool) {
	results := make(map[string]string)
	for _, key := range keys {
		if value, ok := c.Get(key); ok {
			results[key] = value
		} else {
			results[key] = "" // or handle missing keys differently
		}
	}

	if len(results) == 0 {
		return nil, false
	}

	return results, true
}
