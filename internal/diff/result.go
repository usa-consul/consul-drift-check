package diff

// Result represents a single key comparison outcome between two KV namespaces.
type Result struct {
	// Key is the Consul KV key being compared.
	Key string `json:"key"`
	// Status is one of: match, modified, source-only, destination-only.
	Status string `json:"status"`
	// SourceValue is the raw value from the source namespace.
	SourceValue []byte `json:"source_value,omitempty"`
	// DestValue is the raw value from the destination namespace.
	DestValue []byte `json:"dest_value,omitempty"`
	// Weight is an optional score assigned by diff.Apply.
	Weight float64 `json:"weight,omitempty"`
	// Labels holds arbitrary metadata attached by labelmap or classify.
	Labels map[string]string `json:"labels,omitempty"`
	// Tags holds free-form tags attached by the tag package.
	Tags []string `json:"tags,omitempty"`
	// Annotation is a human-readable note attached by diff.Annotate.
	Annotation string `json:"annotation,omitempty"`
}
