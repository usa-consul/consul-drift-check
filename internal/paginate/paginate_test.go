package paginate_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nicholasgasior/consul-drift-check/internal/paginate"
)

func makePairs(n int) api.KVPairs {
	pairs := make(api.KVPairs, n)
	for i := 0; i < n; i++ {
		pairs[i] = &api.KVPair{Key: fmt.Sprintf("key/%d", i), Value: []byte("v")}
	}
	return pairs
}

func TestSplit_EmptyInput(t *testing.T) {
	pages := paginate.Split(nil, paginate.Options{})
	if len(pages) != 0 {
		t.Fatalf("expected 0 pages, got %d", len(pages))
	}
}

func TestSplit_SinglePage(t *testing.T) {
	pages := paginate.Split(makePairs(5), paginate.Options{PageSize: 10})
	if len(pages) != 1 {
		t.Fatalf("expected 1 page, got %d", len(pages))
	}
	if len(pages[0].Pairs) != 5 {
		t.Fatalf("expected 5 pairs, got %d", len(pages[0].Pairs))
	}
}

func TestSplit_MultiplePages(t *testing.T) {
	pages := paginate.Split(makePairs(25), paginate.Options{PageSize: 10})
	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(pages))
	}
	if len(pages[2].Pairs) != 5 {
		t.Fatalf("last page: expected 5 pairs, got %d", len(pages[2].Pairs))
	}
}

func TestSplit_IndexIsOneBased(t *testing.T) {
	pages := paginate.Split(makePairs(10), paginate.Options{PageSize: 5})
	for i, p := range pages {
		if p.Index != uint64(i+1) {
			t.Fatalf("page %d: expected index %d, got %d", i, i+1, p.Index)
		}
	}
}

func TestSplit_DefaultPageSize(t *testing.T) {
	pages := paginate.Split(makePairs(100), paginate.Options{})
	if len(pages) != 1 {
		t.Fatalf("expected 1 page with default size 100, got %d", len(pages))
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	pairs := makePairs(3)
	p := paginate.Page{Index: 1, Pairs: pairs}
	keys := paginate.Keys(p)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "key/0" {
		t.Fatalf("unexpected key: %s", keys[0])
	}
}
