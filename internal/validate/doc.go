// Package validate applies user-defined rules to Consul KV pairs and surfaces
// violations with configurable severity levels.
//
// Rules are loaded from a YAML file and matched against KV keys by prefix.
// Each rule may enforce a maximum value length or require a non-empty value.
package validate
