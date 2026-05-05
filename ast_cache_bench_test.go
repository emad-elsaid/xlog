package xlog

import (
	"sync/atomic"
	"testing"
)

// BenchmarkASTQueryCache_ColdAccess measures performance on first access (cache miss).
func BenchmarkASTQueryCache_ColdAccess(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache := NewASTQueryCache()
		cache.Get("test-key", func() any {
			return "computed-value"
		})
	}
}

// BenchmarkASTQueryCache_WarmAccess measures performance on subsequent accesses (cache hits).
func BenchmarkASTQueryCache_WarmAccess(b *testing.B) {
	cache := NewASTQueryCache()
	cache.Get("test-key", func() any {
		return "computed-value"
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache.Get("test-key", func() any {
			return "computed-value"
		})
	}
}

// BenchmarkASTQueryCache_ExpensiveComputation simulates caching expensive operations.
func BenchmarkASTQueryCache_ExpensiveComputation(b *testing.B) {
	cache := NewASTQueryCache()

	expensiveCompute := func() any {
		// Simulate expensive computation
		result := 0
		for i := 0; i < 1000; i++ {
			result += i
		}
		return result
	}

	b.Run("Uncached", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			expensiveCompute()
		}
	})

	b.Run("Cached", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			cache.Get("expensive", expensiveCompute)
		}
	})
}

// BenchmarkASTQueryCache_ConcurrentReads measures performance under concurrent read load.
func BenchmarkASTQueryCache_ConcurrentReads(b *testing.B) {
	cache := NewASTQueryCache()
	cache.Get("test-key", func() any {
		return "computed-value"
	})

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("test-key", func() any {
				return "computed-value"
			})
		}
	})
}

// BenchmarkASTQueryCache_ConcurrentWrites measures performance under concurrent invalidation.
func BenchmarkASTQueryCache_ConcurrentWrites(b *testing.B) {
	cache := NewASTQueryCache()
	var counter int32

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "test-key"
			cache.Get(key, func() any {
				return atomic.AddInt32(&counter, 1)
			})
			cache.Invalidate(key)
		}
	})
}

// BenchmarkASTQueryCache_MultipleKeys measures performance with multiple independent cache keys.
func BenchmarkASTQueryCache_MultipleKeys(b *testing.B) {
	cache := NewASTQueryCache()

	// Pre-populate cache
	for i := 0; i < 10; i++ {
		cache.Get(string(rune('a'+i)), func() any {
			return i
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cache.Get(string(rune('a'+(i%10))), func() any {
			return i
		})
	}
}
