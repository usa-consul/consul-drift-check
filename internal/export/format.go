package export

import "fmt"

// ParseFormat converts a raw string into a Format, returning an error
// if the value is not recognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatCSV:
		return FormatCSV, nil
	case FormatNDJSON:
		return FormatNDJSON, nil
	default:
		return "", fmt.Errorf("unknown export format %q: must be one of csv, ndjson", s)
	}
}

// Formats returns all supported Format values.
func Formats() []Format {
	return []Format{FormatCSV, FormatNDJSON}
}
