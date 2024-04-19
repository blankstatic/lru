package lru

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestLRU(t *testing.T) {
	cache := NewLRUCache(128)
	for i := 0; i < 256; i++ {
		cache.Add(fmt.Sprint(i), fmt.Sprint(i))
	}
	if cache.Len() != 128 {
		t.Fatalf("bad len: %v", cache.Len())
	}
	for i := 0; i < 128; i++ {
		_, ok := cache.Get(fmt.Sprint(i))
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := cache.Get(fmt.Sprint(i))
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		ok := cache.Remove(fmt.Sprint(i))
		if !ok {
			t.Fatalf("should be contained")
		}
		ok = cache.Remove(fmt.Sprint(i))
		if ok {
			t.Fatalf("should not be contained")
		}
		_, ok = cache.Get(fmt.Sprint(i))
		if ok {
			t.Fatalf("should be deleted")
		}
	}
}

func getRand(tb testing.TB) int64 {
	out, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		tb.Fatal(err)
	}
	return out.Int64()
}

func BenchmarkLRU_Rand(b *testing.B) {
	l := NewLRUCache(8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = getRand(b) % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Add(fmt.Sprint(trace[i]), fmt.Sprint(trace[i]))
		} else {
			_, ok := l.Get(fmt.Sprint(trace[i]))
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkLRU_Freq(b *testing.B) {
	l := NewLRUCache(8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = getRand(b) % 16384
		} else {
			trace[i] = getRand(b) % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(fmt.Sprint(trace[i]), fmt.Sprint(trace[i]))
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		_, ok := l.Get(fmt.Sprint(trace[i]))
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}
