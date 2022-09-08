package reference

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	HTTPProtocol  = "http"
	HTTPSProtocol = "https"
)

// ReferenceClient a reference to a []byte for marshaling
// TODO: pass importer to client to create recursive prefix based imports. This would allow
// for prefixes like "type:rules:json:{}" that could embed info into the actual reference.
type ReferenceClient = func(data string, info Info) ([]byte, Info, error)

// EmbeddedClient uses the reference as raw data to be parsed
type EmbeddedClient struct {
	Format string
}

func (ec *EmbeddedClient) Import(data string, info Info) ([]byte, Info, error) {
	info.Format = ec.Format
	return []byte(data), info, nil
}

// FileClient collects trimmed as a path relative to the root dir.
type FileClient struct {
	Root string
}

func (fc *FileClient) Import(path string, info Info) ([]byte, Info, error) {
	// Import
	if fc.Root == "" {
		return []byte{}, info, fmt.Errorf("importer contained no root location")
	}
	raw, err := os.ReadFile(filepath.Join(fc.Root, path))
	// Info Updates
	if strings.HasSuffix(path, ".json") {
		info.Format = JSONFormat
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		info.Format = YAMLFormat
	}
	return raw, info, err
}

// HTTPClient collects based on url for the network
// In order to make a reference like https://not.a.real.com to work with https:// being the prefix
// It requires the Protocol be specified to be able to reattach the correct http/https prefix
type HTTPClient struct {
	Protocol string
	Client   *http.Client
	Headers  map[string]string
}

func (hc *HTTPClient) Import(data string, info Info) ([]byte, Info, error) {
	url := fmt.Sprintf("%s://%s", hc.Protocol, data)
	// Import
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, info, fmt.Errorf("failed to build http reference request from '%s': %w", url, err)
	}
	if hc.Headers != nil {
		for header, value := range hc.Headers {
			request.Header.Add(header, value)
		}
	}
	response, err := hc.Client.Do(request)
	if err != nil {
		return []byte{}, info, fmt.Errorf("failed to do http reference from '%s': %w", url, err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, info, fmt.Errorf("failed to read body of http reference from '%s': %w", url, err)
	}
	// Info Updates
	contentType := response.Header.Get("Content-Type")
	// By default we assume json
	if contentType == "" || contentType == "application/json" {
		info.Format = JSONFormat
	}
	return body, info, nil
}
