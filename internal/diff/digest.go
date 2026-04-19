package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/hashicorp/consul/api"
)

// DigestEntry holds the SHA-256 digest of a KV namespace snapshot.
type DigestEntry struct {
	Prefix string `json:"prefix"`
	Digest string `json:"digest"`
	KeyCount int    `json:"key_count"`
}

// Digest computes a deterministic SHA-256 fingerprint over a set of KV pairs.
// Pairs are sorted by key before hashing so the result is order-independent.
func Digest(prefix string, pairs []*api.KVPair) DigestEntry {
	sorted := make([]*api.KVPair, len(pairs))
	copy(sorted, pairs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	h := sha256.New()
	for _, p := range sorted {
		h.Write([]byte(p.Key))
		h.Write([]byte{0})
		h.Write(p.Value)
		h.Write([]byte{0})
	}

	return DigestEntry{
		Prefix:   prefix,
		Digest:   hex.EncodeToString(h.Sum(nil)),
		KeyCount: len(pairs),
	}
}

// DigestsEqual returns true when two DigestEntry values represent identical
// namespace content.
func DigestsEqual(a, b DigestEntry) bool {
	return a.Digest == b.Digest
}
