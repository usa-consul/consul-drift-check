// Package paginate splits large api.KVPairs result sets into fixed-size pages
// for memory-efficient processing when comparing Consul KV namespaces that
// contain thousands of keys.
//
// Usage:
//
//	pages := paginate.Split(pairs, paginate.Options{PageSize: 250})
//	for _, p := range pages {
//		// process p.Pairs
//	}
package paginate
