package merge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const (
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
)

var (
	httpReferenceClient = http.Client{
		Timeout: 5 * time.Second,
	}

	schemaCache = map[string]*gojsonschema.Schema{}
)

// ImportReference load a generic reference into memory.
// Supports file based references through "file://". See ImportFile for more details
func ImportReference[T any](reference string) (T, error) {
	if strings.HasPrefix(reference, FilePrefix) {
		return ImportFile[T](reference[len(FilePrefix):])
	}
	if strings.HasPrefix(reference, HTTPPrefix) || strings.HasPrefix(reference, HTTPSPrefix) {
		return ImportHTTP[T](reference)
	}
	if strings.HasPrefix(reference, JSONPrefix) {
		var document T
		return document, json.Unmarshal([]byte(reference[len(JSONPrefix):]), &document)
	}
	if strings.HasPrefix(reference, YAMLPrefix) {
		var document T
		return document, yaml.Unmarshal([]byte(reference[len(JSONPrefix):]), &document)
	}
	var empty T
	return empty, fmt.Errorf("unknown reference '%s'", reference)
}

// ImportFile unmarshals data from file into the generic type.
// Supports YAML (.yaml/.yml) and JSON (.json) formats
func ImportFile[T any](path string) (T, error) {
	var document T
	raw, err := os.ReadFile(path)
	if err != nil {
		return document, fmt.Errorf("failed to import file: %w", err)
	}
	if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(raw, &document)
		if err != nil {
			return document, fmt.Errorf("failed to import json file: %w", err)
		}
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		err = yaml.Unmarshal(raw, &document)
		if err != nil {
			return document, fmt.Errorf("failed to import yaml file: %w", err)
		}
	} else {
		return document, fmt.Errorf("failed to import file, unknown file type '%s'", path)
	}
	return document, nil
}

// ImportHTTP unmarshals data from http address into the generic type.
// Supports Content-Type application/json
func ImportHTTP[T any](url string) (T, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	var document T
	if err != nil {
		return document, fmt.Errorf("failed to build import http reference request from '%s': %w", url, err)
	}
	response, err := httpReferenceClient.Do(request)
	if err != nil {
		return document, fmt.Errorf("failed to do import http reference from '%s': %w", url, err)
	}
	contentType := response.Header.Get("Content-Type")
	if contentType == "application/json" {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return document, fmt.Errorf("failed to read import http reference body: %w", err)
		}
		return document, json.Unmarshal(body, &document)
	}
	return document, fmt.Errorf("unsupported import http reference content-type '%s'", contentType)
}
