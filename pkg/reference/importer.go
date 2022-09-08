package reference

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// Prefixes

	// FilePrefix reference prefix that resolves the remainder as a file path
	FilePrefix = "file://"
	// HTTPPrefix reference prefix that resolves to HTTP data
	HTTPPrefix = "http://"
	// HTTPSPrefix reference prefix that resolves to HTTPS data
	HTTPSPrefix = "https://"
	// JSONPrefix reference prefix that resolves the remainder as raw json string data
	JSONPrefix = "json:"
	// YAMLPrefix reference prefix that resolves the remainder as raw yaml string data
	YAMLPrefix = "yaml:"
	// SchemaPrefix reference prefix that resolves a file from a schema directory
	SchemaPrefix = "schema://"

	// Formats

	JSONFormat = "json"
	YAMLFormat = "yaml"
)

// Info details about your reference.
type Info struct {
	Type      string
	Format    string
	Reference string
}

// Resolution is the final result of a reference.
type Resolution[T any] struct {
	Info Info
	Data T
}

// Resolver is a map of prefixes to clients to resolve the remainder of the import process.
type Resolver map[string]ReferenceClient

// NewDefaultResolver creates an Resolver with default prefixes applied
func NewDefaultResolver(fileRoot, schemaRoot string, timeout time.Duration, userAgent string) Resolver {
	client := &http.Client{
		Timeout: timeout,
	}
	headers := map[string]string{
		"User-Agent": userAgent,
	}
	return Resolver{
		FilePrefix:   (&FileClient{Root: fileRoot}).Import,
		SchemaPrefix: (&FileClient{Root: schemaRoot}).Import,
		HTTPPrefix:   (&HTTPClient{Protocol: "http", Client: client, Headers: headers}).Import,
		HTTPSPrefix:  (&HTTPClient{Protocol: "https", Client: client, Headers: headers}).Import,
		JSONPrefix:   (&EmbeddedClient{Format: JSONFormat}).Import,
		YAMLFormat:   (&EmbeddedClient{Format: YAMLFormat}).Import,
	}
}

// ImportInterface import the reference
func (r Resolver) ImportInterface(reference string) (Resolution[interface{}], error) {
	return Import[interface{}](r, reference)
}

func (r Resolver) ImportMap(reference string) (Resolution[map[string]interface{}], error) {
	return Import[map[string]interface{}](r, reference)
}

// Import the provided reference using the importer then un-marshaling to the expected type
func Import[T any](importer Resolver, reference string) (Resolution[T], error) {
	for prefix, client := range importer {
		if strings.HasPrefix(reference, prefix) {
			// Load
			data, info, err := client(reference[len(prefix):], Info{
				Reference: reference,
			})
			if err != nil {
				return Resolution[T]{}, fmt.Errorf("failed to get client data for reference '%s': %w", reference, err)
			}
			// Un-marshaling
			if info.Format == JSONFormat {
				var document T
				err := json.Unmarshal(data, &document)
				return Resolution[T]{
					Info: info,
					Data: document,
				}, err
			} else if info.Format == YAMLFormat {
				var document T
				err := yaml.Unmarshal(data, &document)
				return Resolution[T]{
					Info: info,
					Data: document,
				}, err
			} else {
				return Resolution[T]{}, fmt.Errorf("unknown format '%s' for reference '%s'", info.Format, reference)
			}
		}
	}
	return Resolution[T]{}, fmt.Errorf("unknown reference '%s'", reference)
}
