// UpdateGuard Agent is the only component mounted with the Docker socket.
// It exposes a small authenticated API, never a generic Docker proxy.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/updateguard/updateguard/internal/dockerapi"
)

func main() {
	token := os.Getenv("UPDATEGUARD_AGENT_TOKEN")
	if token == "" {
		log.Fatal("UPDATEGUARD_AGENT_TOKEN must be set")
	}

	engine := dockerapi.New(os.Getenv("UPDATEGUARD_DOCKER_SOCKET"))
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := engine.Ping(r.Context()); err != nil {
			http.Error(w, "docker engine unavailable", http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/v1/services", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		services, err := engine.Discover(r.Context(), "local")
		if err != nil {
			http.Error(w, "container discovery failed", http.StatusBadGateway)
			return
		}
		writeJSON(w, services)
	})

	addr := os.Getenv("UPDATEGUARD_AGENT_LISTEN_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	log.Printf("UpdateGuard restricted agent listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, requireToken(token, mux)))
}

func requireToken(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") ||
			strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") != token {
			http.Error(w, "agent authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
