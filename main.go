package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"./blockchain"
)

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//blockchain.Init()

	go func() {

	}()
	log.Fatal(run())

}

func getPort () string {
	if (os.Getenv("PORT") != "") {
		return os.Getenv("PORT")
	}

	defaultPort := "8080"

	log.Println("Missing HTTP Server Listening port using", defaultPort, "as default port")

	return defaultPort
}

// web server
func run() error {
	muxRouter := makeMuxRouter()
	httpPort := getPort()
	log.Println("HTTP Server Listening on port :", httpPort)

	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        muxRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockchain).Methods("GET")
	router.HandleFunc("/", handleWriteBlock).Methods("POST")
	return router
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := blockchain.GetBlockChain()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newBlock := blockchain.AddBlock(m.BPM)

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}