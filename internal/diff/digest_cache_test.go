package diff

import (
	"testing"
	"time"
)

func sampleDigest(prefix string) DigestEntry {
	return DigestEntry{Prefix: prefix, Digest: "abc123", KeyCount: 3}
}

func TestDigestCache_MissOnEmpty(t *testing.T) {
	c := NewDigestCache(time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestDigestCache_SetAndGet(t *testing.T) {
	c := NewDigestCache(time.Minute)
	e := sampleDigest("app/")
	c.Set("app/", e)
	got, ok := c.Get("app/")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got.Digest != e.Digest {
		t.Fatalf("digest mismatch: got %s", got.Digest)
	}
}

func TestDigestCache_ExpiredEntry(t *testing.T) {
	c := NewDigestCache(time.Millisecond)
	c.Set("k", sampleDigest("k"))
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestDigestCache_Invalidate(t *testing.T) {
	c := NewDigestCache(time.Minute)
	c.Set("k", sampleDigest("k"))
	c.Invalidate("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestDigestCache_ZeroTTL_DefaultsToThirtySeconds(t *testing.T) {
	c := NewDigestCache(0)
	if c.ttl != 30*time.Second {
		t.Fatalf("expected 30s TTL, got %s", c.ttl)
	}
}
