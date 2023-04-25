package tests

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/duanjr/trustchain/blockchain"
	"github.com/duanjr/trustchain/node"
)

func TestGetBlockchainHandler(t *testing.T) {
	testNode := node.NewNode()
	testNode.Blockchain.AddRecord("record 1")
	testNode.Blockchain.AddRecord("record 2")
	testNode.Blockchain.AddRecord("record 3")
	testNode.Blockchain.AddRecord("record 4")
	testNode.Blockchain.AddRecord("record 5")
	testNode.Blockchain.AddRecord("record 6")
	testNode.Blockchain.AddRecord("record 7")
	testNode.Blockchain.AddRecord("record 8")
	testNode.Blockchain.AddRecord("record 9")
	testNode.Blockchain.AddRecord("record 10")

	req := httptest.NewRequest("GET", "/blocks", nil)
	w := httptest.NewRecorder()

	testNode.GetBlockchain(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var blocks []*blockchain.Block
	json.Unmarshal(body, &blocks)

	if len(blocks) != 2 {
		t.Errorf("Expected 2 blocks in the chain, got %d", len(blocks))
	}
}

func TestAddPeerHandler(t *testing.T) {
	testNode := node.NewNode()

	req := httptest.NewRequest("POST", "/add-peer", strings.NewReader("http://localhost:8081"))
	w := httptest.NewRecorder()

	testNode.AddPeerHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	// There should be only 1 peer in the list.
	if len(testNode.Peers) != 1 {
		t.Errorf("Expected 1 peer in the list, got %d", len(testNode.Peers))
	}
}

// TestSynchronizeBlockchain will test the SynchronizeBlockchain function.
func TestSynchronizeBlockchain(t *testing.T) {
	// Start the first node.
	node1 := node.NewNode()
	router1 := mux.NewRouter()
	router1.HandleFunc("/blocks", node1.GetBlockchain).Methods("GET")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router1))
	}()

	// Start the second node.
	node2 := node.NewNode()
	router2 := mux.NewRouter()
	router2.HandleFunc("/blocks", node2.GetBlockchain).Methods("GET")
	go func() {
		log.Fatal(http.ListenAndServe(":8081", router2))
	}()

	// Add a record to the first node.
	node1.Blockchain.AddRecord("record 1")

	// Add the first node as a peer to the second node.
	node2.AddPeer("http://localhost:8080")

	// Synchronize the second node's blockchain.
	node2.SynchronizeBlockchain()

	if len(node2.Blockchain.Blocks) != len(node1.Blockchain.Blocks) {
		t.Errorf("Expected %d blocks in node2's chain, got %d", len(node1.Blockchain.Blocks), len(node2.Blockchain.Blocks))
	}
}
