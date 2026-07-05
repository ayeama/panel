package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/ayeama/panel/internal"
	"github.com/ayeama/panel/internal/middleware"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	db, err := sql.Open("sqlite3", "panel.db?_busy_timeout=100&_foreign_keys=true&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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
			cpu := rand.Float64() * 100
			memory := rand.Float64() * 100
			disk := rand.Float64() * 100

			msg := []byte(fmt.Sprintf("{\"cpu\":%.2f,\"memory\":%.2f,\"disk\":%.2f}", cpu, memory, disk))
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
			time.Sleep(time.Second)
		}
	})

	addr := "0.0.0.0:8000"
	cert, key, err := internal.GetCertificate()
	if err != nil {
		log.Fatal(err)
	}

	handler := middleware.Log(mux)
	handler = middleware.Cors(handler)

	log.Printf("starting https://%s\n", addr)

	err = http.ListenAndServeTLS(addr, cert, key, handler)
	if err != nil {
		log.Fatal(err)
	}
}
