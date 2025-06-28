
# Go LRU Cache

A thread-safe, high-performance Least Recently Used (LRU) cache implementation in Go. This package uses a doubly-linked list and hashmap to achieve O(1) time complexity for both Get and Put operations. Perfect for caching scenarios where you need fast access and automatic eviction of least recently used items.

## Features

- **O(1) Operations:** Efficient Get and Put operations using a doubly-linked list and hashmap.
- **Thread-Safe:** Safe for concurrent use with built-in synchronization.
- **Customizable Capacity:** Set the cache size to fit your needs.
- **Runnable Example:** Includes a demo to test the cache in action.

## Installation
```
go get github.com/CHIRANTAN-001/lrucache/pkg/lrucache

```

## Example
The repository includes a runnable example in example/main.go. To try it:
```
go run example/main.go

```

## Thread Safety

The cache is designed to be thread-safe, using Go’s sync.RWMutex to handle concurrent Get and Put operations. You can safely use it in multi-goroutine environments without additional synchronization.

## Star the Repo

If you find this package useful, please give it a ⭐ on GitHub to show your support!
