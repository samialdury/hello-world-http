package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

var (
	VersionSHA = "unknown-sha"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Incoming request", "method", r.Method, "host", r.Host, "path", r.URL.Path)

		res := Response{
			VersionSHA:    VersionSHA,
			RemoteAddr:    r.RemoteAddr,
			XRealIP:       r.Header.Get("X-Real-Ip"),
			XForwardedFor: r.Header.Get("X-Forwarded-For"),
		}

		writeJSON(w, http.StatusOK, res)
	})

	port := getEnv("PORT", "8080")

	slog.Info("Starting app...", "port", port)

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), mux)
}

type Response struct {
	VersionSHA    string `json:"version_sha"`
	RemoteAddr    string `json:"remote_addr"`
	XRealIP       string `json:"x_real_ip"`
	XForwardedFor string `json:"x_forwarded_for"`
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
