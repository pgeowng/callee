package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/review", handleReview)
	mux.HandleFunc("/answer", handleAnswer)

	fmt.Println(":4532")
	http.ListenAndServe(":4532", mux)
}

type Card struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Status struct {
	Status string `json:"status"`
	// ReviewEnd bool `json:"review_end"`
}

var StatusReviewEnd Status = Status{Status: "review_end"}

var cardsIdx = 0
var cards []Card = []Card{
	{"how big is earth?", "very big"},
	{"how big is sun?", "bigger than earth"},
	{"４月に学校が始まります。", "始まる"},
	{"駅前のビルに郵便局が移ります。", "移[うつ,うつる;k2]る"},
}

func handleReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
		return
	}

	var err error

	if cardsIdx >= len(cards) {
		err = json.NewEncoder(w).Encode(StatusReviewEnd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
				panic(err)
			}
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cards[cardsIdx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			panic(err)
		}
		return
	}
}

type AnswerParams struct {
	Button string `json:"button"`
}

func handleAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "not found"}`))
		return
	}

	var err error
	var params AnswerParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			panic(err)
		}
		return
	}

	fmt.Println(cardsIdx, len(cards))
	if cardsIdx+1 >= len(cards) {
		err = json.NewEncoder(w).Encode(StatusReviewEnd)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
				panic(err)
			}
			return
		}
		return
	}

	cardsIdx++

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cards[cardsIdx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			panic(err)
		}
		return
	}

}
