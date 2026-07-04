package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

	mux.HandleFunc("GET /instances/{id}", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{\"id\":\"1\",\"name\":\"another\",\"image\":\"minecraft\",\"status\":\"stopped\"}")
		w.Header().Add("Content-Type", "application/json")
	})

	mux.HandleFunc("GET /instances", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "[{\"id\":\"1\",\"name\":\"another\",\"image\":\"minecraft\",\"status\":\"stopped\"},{\"id\":\"2\",\"name\":\"one\",\"image\":\"valheim\",\"status\":\"running\"}]")
		w.Header().Add("Content-Type", "application/json")
	})

	mux.HandleFunc("GET /instances/{id}/attach", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				return
			}

			if err := c.WriteMessage(mt, msg); err != nil {
				return
			}
		}
	})

	mux.HandleFunc("GET /instances/{id}/stats", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			if err := c.WriteMessage(websocket.TextMessage, []byte("{\"cpu\":10,\"memory\":70,\"disk\":35}")); err != nil {
				return
			}
			time.Sleep(time.Second)
		}
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
