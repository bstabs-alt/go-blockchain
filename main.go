package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

type Block struct {
	Index     int
	Timestamp string
	Coins     int
	PrevHash  string
	Hash      string
}

type Message struct {
	Coins int
}

var Blockchain []Block

func main() {
	if err := godotenv.Load(); err != nil {
		panic(".env file failed to load")
	}

	go func() {
		t := time.Now().UTC()
		baseBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(baseBlock)
		Blockchain = append(Blockchain, baseBlock)
	}()

	if err := run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}

func run() error {
	mux := makeMuxRouter()
	s := http.Server{
		Addr:           "localhost:" + os.Getenv("PORT"),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening on ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("Failed to serve: %w", err)
	}
	return nil
}
