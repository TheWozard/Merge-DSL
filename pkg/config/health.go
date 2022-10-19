package config

import (
	"encoding/json"
	"net/http"
)

type HealthEndpoint struct {
	Config *Config
}

// ServeHealth serves health data. Current this is just the raw config
func (h HealthEndpoint) ServeHealth(w http.ResponseWriter, r *http.Request) {
	raw, err := json.Marshal(h.Config)
	if err != nil {
		http.Error(w, "Failed to marshal health data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(raw)
}
