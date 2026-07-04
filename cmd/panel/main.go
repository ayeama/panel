package main

import (
	"io"
	"log"
	"net/http"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO add configuration

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /instances", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "[{\"id\":1,\"name\":\"another\",\"status\":\"stopped\"},{\"id\":2,\"name\":\"one\",\"status\":\"running\"}]")
		w.Header().Add("Content-Type", "application/json")
	})

	addr := "localhost:8000"
	cert := "server.crt"
	key := "server.key"
	handler := cors(mux)

	log.Printf("starting https://%s\n", addr)

	err := http.ListenAndServeTLS(addr, cert, key, handler)
	if err != nil {
		log.Fatal(err)
	}
}
