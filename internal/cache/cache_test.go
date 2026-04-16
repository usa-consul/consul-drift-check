package cache_test

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/your-org/consul-drift-check/internal/cache"
)

func samplePairs() api.KVPairs {
	return api.KVPairs{
		{Key: "app/db", Value: []byte("postgres")},
		{Key: "app/port", Value: []byte("5432")},
	}
}

func TestGet_MissOnEmpty(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("app/")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSetAndGet_ReturnsValue(t *testing.T) {
	c := cache.New(time.Minute)
	pairs := samplePairs()
	c.Set("app/", pairs)

	got, ok := c.Get("app/")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != len(pairs) {
		t.Fatalf("expected %d pairs, got %d", len(pairs), len(got))
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("app/", samplePairs())

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("app/")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("app/", samplePairs())
	c.Invalidate("app/")

	_, ok := c.Get("app/")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("app/", samplePairs())
	c.Set("svc/", samplePairs())
	c.Flush()

	if _, ok := c.Get("app/"); ok {
		t.Error("expected app/ to be flushed")
	}
	if _, ok := c.Get("svc/"); ok {
		t.Error("expected svc/ to be flushed")
	}
}
