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
	router.HandleFunc("/pki/update", node.UpdatePKIRecord).Methods("POST")
	router.HandleFunc("/pki/query", node.QueryPKIRecord).Methods("POST")
	router.HandleFunc("/trust/submit", node.TrustSubmitRecord).Methods("POST")
	router.HandleFunc("/trust/query-direct", node.DirectTrustQueryRecord).Methods("POST")
	router.HandleFunc("/trust/query-comp", node.CompTrustQuery).Methods("POST")
	router.HandleFunc("/trust/query-comp-calc", node.CalcCompTrustQuery).Methods("POST")
	fmt.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
