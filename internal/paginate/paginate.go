// Package paginate provides helpers for iterating over large Consul KV
// result sets in fixed-size pages.
package paginate

import "github.com/hashicorp/consul/api"

// Page holds a slice of KV pairs representing a single page of results.
type Page struct {
	Index uint64
	Pairs api.KVPairs
}

// Options controls pagination behaviour.
type Options struct {
	// PageSize is the maximum number of entries returned per page.
	// Defaults to 100 if zero.
	PageSize int
}

func (o Options) pageSize() int {
	if o.PageSize <= 0 {
		return 100
	}
	return o.PageSize
}

// Split divides pairs into a slice of Pages, each containing at most
// Options.PageSize entries. The Index field is set to the 1-based page number.
func Split(pairs api.KVPairs, opts Options) []Page {
	size := opts.pageSize()
	if len(pairs) == 0 {
		return nil
	}

	var pages []Page
	for i := 0; i < len(pairs); i += size {
		end := i + size
		if end > len(pairs) {
			end = len(pairs)
		}
		pages = append(pages, Page{
			Index: uint64(len(pages) + 1),
			Pairs: pairs[i:end],
		})
	}
	return pages
}

// Keys returns all keys contained in a Page.
func Keys(p Page) []string {
	out := make([]string, len(p.Pairs))
	for i, kv := range p.Pairs {
		out[i] = kv.Key
	}
	return out
}
