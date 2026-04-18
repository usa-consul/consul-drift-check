// Package overlap detects KV keys that are reachable under more than one
// namespace prefix. This can indicate misconfigured service registrations or
// accidental duplication across environment boundaries.
//
// Usage:
//
//	results := overlap.Find(pairs, []string{"prod/", "staging/", "shared/"})
//	for _, r := range results {
//		fmt.Printf("key %s appears in: %v\n", r.Key, r.Prefixes)
//	}
package overlap
