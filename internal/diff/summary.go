package diff

import "github.com/your-org/consul-drift-check/internal/diff"

// StatusAdded indicates a key present only in source.
const StatusAdded = "added"

// StatusRemoved indicates a key present only in destination.
const StatusRemoved = "removed"

// StatusModified indicates a key whose value differs between source and destination.
const StatusModified = "modified"

// StatusMatch indicates identical values in both source and destination.
const StatusMatch = "match"
