package server

import (
	"fmt"
	"github.com/duanjr/trustchain/node"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer() {
	node := node.NewNode()

	router := mux.NewRouter()
	router.HandleFunc("/add-record", node.AddRecord).Methods("POST")
	router.HandleFunc("/blocks", node.GetBlockchain).Methods("GET")
	router.HandleFunc("/add-peer", node.AddPeerHandler).Methods("POST")
	router.HandleFunc("/pki/register", node.AddPKIRecord).Methods("POST")
	fmt.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}