package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	DocumentAPI  bool `json:"document-api" yaml:"document-api"`
	GraphQLAPI   bool `json:"graphql-api" yaml:"graphql-api"`
	SandboxAPI   bool `json:"sandbox-api" yaml:"sandbox-api"`
	ReferenceAPI bool `json:"reference-api" yaml:"reference-api"`
}

// GetNewConfig parses the current config. The config is rebuilt each time
func GetNewConfig() *Config {
	restAPI := flag.Bool("document-api", envBool("ENABLE_DOCUMENT_API", false), "enables the document endpoints")
	graphQL := flag.Bool("graphql-api", envBool("ENABLE_GRAPHQL_API", false), "enables the graphql endpoints")
	sandboxAPI := flag.Bool("sandbox-api", envBool("ENABLE_SANDBOX_API", false), "enables the sandbox endpoints")
	referenceAPI := flag.Bool("reference-api", envBool("ENABLE_REFERENCE_API", false), "enables the reference endpoints")
	flag.Parse()
	return &Config{
		DocumentAPI:  *restAPI,
		GraphQLAPI:   *graphQL,
		SandboxAPI:   *sandboxAPI,
		ReferenceAPI: *referenceAPI,
	}
}

// envBool is os.Getenv with basic boolean validation and a default fallback
func envBool(key string, def bool) bool {
	if value := strings.ToLower(os.Getenv(key)); value != "" {
		return value == "t" || value == "true" || value == "yes" || value == "y"
	}
	return def
}
