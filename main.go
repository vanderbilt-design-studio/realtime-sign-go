package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/olahol/melody"
)

func main() {
	var currStatus interface{}

	router := chi.NewRouter()
	m := melody.New()
	m.HandleConnect(func(sock *melody.Session) {
		if json, err := json.Marshal(&currStatus); err == nil {
			sock.Write(json)
		}
	})

	apiKey := os.Getenv("RSIGN_API_KEY")
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != apiKey {
			render.Status(r, 401)
			render.PlainText(w, r, "Invalid API key")
			return
		}
		var newStatus interface{}
		if err := render.Decode(r, &newStatus); err != nil {
			render.Status(r, 400)
			render.PlainText(w, r, "Invalid JSON")
			return
		}
		if reflect.DeepEqual(currStatus, newStatus) {
			render.JSON(w, r, render.M{"updated": false})
			return
		}
		currStatus = newStatus
		if json, err := json.Marshal(&currStatus); err == nil {
			m.Broadcast(json)
			fmt.Printf("%s\n", json)
		}
		render.JSON(w, r, render.M{"updated": true})
	})

	addr := ":3000"
	if value, ok := os.LookupEnv("RSIGN_ADDR"); ok {
		addr = value
	}
	http.ListenAndServe(addr, router)
}
