// Package transform provides a composable pipeline for applying sequential
// transformations to Consul KV pairs before comparison or export.
//
// Stages are registered by name and executed in insertion order. Built-in
// helpers such as StripPrefix and LowerKeys cover the most common cases;
// callers may supply arbitrary functions for custom behaviour.
package transform
