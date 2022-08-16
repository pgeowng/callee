package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pgeowng/goanki/notes"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(404)
			return
		}

		path := r.URL.Path

		if path == "/notes" {
			notes, err := notes.QueryAllNotes()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(notes)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			return
		}

		w.WriteHeader(404)
	})

	addr := ":2313"
	fmt.Println("listen on", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type H = map[string]interface{}

func EndPointHandler(func(input interface{})) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&reqBody)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			err = json.Encoder(w).Encode(H{"error": err.Error()})
			return
		}

		// ...

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			err = json.Encoder(w).Encode(H{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.Encoder(w).Encode(H{"result": result})
		return

	}
}

func responseError() {

}
