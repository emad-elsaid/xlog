package xlog

import "sync"

// ASTQueryCache provides thread-safe caching for expensive AST-based queries.
// Extensions can use this to avoid redundant AST parsing when computing
// aggregate data across multiple pages.
type ASTQueryCache struct {
	mu      sync.RWMutex
	queries map[string]any
}

// NewASTQueryCache creates a new cache instance.
func NewASTQueryCache() *ASTQueryCache {
	return &ASTQueryCache{
		queries: make(map[string]any),
	}
}

// Get retrieves a cached value for the given key, or computes it using the
// provided function if not cached. The compute function is only called once
// even under concurrent access for the same key.
func (c *ASTQueryCache) Get(key string, compute func() any) any {
	// Fast path: check with read lock
	c.mu.RLock()
	cached, ok := c.queries[key]
	c.mu.RUnlock()

	if ok {
		return cached
	}

	// Slow path: compute with write lock
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock (another goroutine may have computed)
	if cached, ok := c.queries[key]; ok {
		return cached
	}

	result := compute()
	c.queries[key] = result
	return result
}

// Invalidate removes a cached value for the given key.
// Subsequent Get calls will recompute the value.
func (c *ASTQueryCache) Invalidate(key string) {
	c.mu.Lock()
	delete(c.queries, key)
	c.mu.Unlock()
}

// InvalidateAll clears all cached values.
// Useful when multiple pages have changed and selective invalidation is impractical.
func (c *ASTQueryCache) InvalidateAll() {
	c.mu.Lock()
	c.queries = make(map[string]any)
	c.mu.Unlock()
}
