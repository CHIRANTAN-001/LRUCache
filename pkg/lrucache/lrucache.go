package lrucache

import (
	"errors"
	"sync"
)

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
	mutex    sync.RWMutex
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
		mutex:    sync.RWMutex{},
	}, nil
}

// Get retrieves the value for a given key from the cache.
// Returns the value and true if found, empty string and false otherwise.
func (c *LRUCache) Get(key string) (string, bool) {
	c.mutex.Lock() // Use write lock since we modify the list order
	defer c.mutex.Unlock()
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

	// Remove the node from its current position
	c.removeNode(node)

	// Add the node to the head of the list
	c.addToHead(node)
}

// removeNode removes a node from the doubly linked list.
func (c *LRUCache) removeNode(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		c.Head = node.Next // If it's the head, move head to next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		c.Tail = node.Prev // If it's the tail, move tail to prev
	}
}

// addToHead adds a node to the head of the doubly linked list.
func (c *LRUCache) addToHead(node *Node) {
	node.Prev = nil
	node.Next = c.Head

	if c.Head != nil {
		c.Head.Prev = node
	}
	c.Head = node

	if c.Tail == nil {
		c.Tail = node
	}
}

// removeTail removes the least recently used item (tail) from the cache.
func (c *LRUCache) removeTail() *Node {
	if c.Tail == nil {
		return nil
	}

	tailNode := c.Tail
	c.removeNode(tailNode)
	return tailNode
}

// Put adds a key-value pair to the cache.
// If the key already exists, it updates the value and moves the node to the head.
func (c *LRUCache) Put(key string, value string) {
	// Lock the cache for writing to ensure thread safety
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If the key already exists, update the value and move to head
	if node, ok := c.Cache[key]; ok {
		node.Value = value
		// Move the node to the head of the list
		c.moveToHead(node)
		return
	}

	// Create a new node
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	// If the cache is at capacity, remove the least recently used item
	if len(c.Cache) >= c.Capacity {
		tail := c.removeTail()
		if tail != nil {
			delete(c.Cache, tail.Key)
		}
	}
	
	// Add the new node to the cache
	c.Cache[key] = newNode
	c.addToHead(newNode)
}

// Clear removes all items from the cache.
func (c *LRUCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Head = nil
	c.Tail = nil
	c.Cache = make(map[string]*Node)
}

// Size returns the current number of items in the cache.
func (c *LRUCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.Cache)
}

// IsEmpty checks if the cache is empty.
func (c *LRUCache) IsEmpty() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.Cache) == 0
}

// Contains checks if the cache contains a specific key.
func (c *LRUCache) Has(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.Cache[key]
	return ok
}

