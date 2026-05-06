package xlog

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestASTQueryCache_GetComputesOnFirstCall(t *testing.T) {
	cache := NewASTQueryCache()
	var computeCount int32

	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "computed-value"
	}

	result := cache.Get("test-key", compute)

	if result != "computed-value" {
		t.Errorf("expected 'computed-value', got %v", result)
	}

	if atomic.LoadInt32(&computeCount) != 1 {
		t.Errorf("expected compute called once, got %d", computeCount)
	}
}

func TestASTQueryCache_GetReturnsCachedValue(t *testing.T) {
	cache := NewASTQueryCache()
	var computeCount int32

	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "computed-value"
	}

	// First call - should compute
	result1 := cache.Get("test-key", compute)
	// Second call - should return cached
	result2 := cache.Get("test-key", compute)

	if result1 != result2 {
		t.Errorf("expected same result, got %v and %v", result1, result2)
	}

	if atomic.LoadInt32(&computeCount) != 1 {
		t.Errorf("expected compute called once, got %d", computeCount)
	}
}

func TestASTQueryCache_InvalidateClearsCache(t *testing.T) {
	cache := NewASTQueryCache()
	var computeCount int32

	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "computed-value"
	}

	// First call
	cache.Get("test-key", compute)

	// Invalidate
	cache.Invalidate("test-key")

	// Second call should recompute
	cache.Get("test-key", compute)

	if atomic.LoadInt32(&computeCount) != 2 {
		t.Errorf("expected compute called twice, got %d", computeCount)
	}
}

func TestASTQueryCache_ConcurrentAccess(t *testing.T) {
	cache := NewASTQueryCache()
	var computeCount int32

	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "computed-value"
	}

	var wg sync.WaitGroup
	numWorkers := 100

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Get("test-key", compute)
		}()
	}

	wg.Wait()

	// Despite 100 concurrent calls, compute should only be called once
	// (or a very small number of times due to race on initial check)
	count := atomic.LoadInt32(&computeCount)
	if count > 5 {
		t.Errorf("expected compute called ~1 time, got %d", count)
	}
}

func TestASTQueryCache_MultipleKeys(t *testing.T) {
	cache := NewASTQueryCache()

	result1 := cache.Get("key1", func() any { return "value1" })
	result2 := cache.Get("key2", func() any { return "value2" })

	if result1 != "value1" {
		t.Errorf("expected 'value1', got %v", result1)
	}

	if result2 != "value2" {
		t.Errorf("expected 'value2', got %v", result2)
	}

	// Invalidating one key shouldn't affect the other
	cache.Invalidate("key1")

	var computeCount int32
	cache.Get("key1", func() any {
		atomic.AddInt32(&computeCount, 1)
		return "new-value1"
	})

	cache.Get("key2", func() any {
		atomic.AddInt32(&computeCount, 1)
		return "value2"
	})

	// key1 should have recomputed, key2 should use cache
	if atomic.LoadInt32(&computeCount) != 1 {
		t.Errorf("expected compute called once (for key1), got %d", computeCount)
	}
}

func TestASTQueryCache_InvalidateAll(t *testing.T) {
	cache := NewASTQueryCache()

	cache.Get("key1", func() any { return "value1" })
	cache.Get("key2", func() any { return "value2" })
	cache.Get("key3", func() any { return "value3" })

	cache.InvalidateAll()

	var computeCount int32
	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "recomputed"
	}

	cache.Get("key1", compute)
	cache.Get("key2", compute)
	cache.Get("key3", compute)

	if atomic.LoadInt32(&computeCount) != 3 {
		t.Errorf("expected all 3 keys to recompute, got %d", computeCount)
	}
}

func TestASTQueryCache_IntegrationWithPageEvents(t *testing.T) {
	cache := NewASTQueryCache()
	var computeCount int32

	compute := func() any {
		atomic.AddInt32(&computeCount, 1)
		return "computed-value"
	}

	// First call - computes
	cache.Get("test-key", compute)

	// Setup cache invalidation on page change
	Listen(PageChanged, func(p Page) error {
		cache.Invalidate("test-key")
		return nil
	})

	// Trigger page change event
	Trigger(PageChanged, NewPage("test"))

	// Second call - should recompute due to invalidation
	cache.Get("test-key", compute)

	if atomic.LoadInt32(&computeCount) != 2 {
		t.Errorf("expected compute called twice (once before, once after event), got %d", computeCount)
	}
}
