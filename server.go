package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func makeMuxRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleGet)
	mux.HandleFunc("POST /", handlePost)
	return mux
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	bytes, err := json.MarshalIndent(Blockchain, "", "	")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(bytes))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var m Message

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&m); err != nil {
		respondWithJSON(w, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	prevBlock := Blockchain[len(Blockchain)-1]

	newBlock, err := genBlock(Blockchain[len(Blockchain)-1], m.Coins)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, r.Body)
		return
	}

	if err := isValidBlock(prevBlock, newBlock); err != nil {
		respondWithJSON(w, http.StatusInternalServerError, r.Body)
		return
	}
	newBlockchain := append(Blockchain, newBlock)
	replaceChain(newBlockchain)
	spew.Dump(Blockchain)

	respondWithJSON(w, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, status int, body any) {
	response, err := json.MarshalIndent(body, "", "	")
	if err != nil {
		w.WriteHeader(status)
		_, _ = w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(response)
}
