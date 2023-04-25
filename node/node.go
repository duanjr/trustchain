package node

import (
	"encoding/json"
	"fmt"
	"github.com/duanjr/trustchain/blockchain"
	"github.com/duanjr/trustchain/pki"
	"net/http"
)

type Node struct {
	Blockchain *blockchain.Blockchain
	Peers      []string
}

func NewNode() *Node {
	return &Node{blockchain.NewBlockchain(), []string{}}
}

func (n *Node) AddRecord(w http.ResponseWriter, r *http.Request) {
	var record string

	if _, err := fmt.Fscanf(r.Body, "%s", &record); err != nil {
		http.Error(w, "Error reading record", http.StatusBadRequest)
		return
	}

	n.Blockchain.AddRecord(record)
	w.WriteHeader(http.StatusNoContent)
}

func (n *Node) GetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n.Blockchain.Blocks)
}

func (n *Node) AddPeer(peer string) {
	n.Peers = append(n.Peers, peer)
}

func (n *Node) AddPeerHandler(w http.ResponseWriter, r *http.Request) {
	var peer string

	if _, err := fmt.Fscanf(r.Body, "%s", &peer); err != nil {
		http.Error(w, "Error reading peer", http.StatusBadRequest)
		return
	}

	n.AddPeer(peer)
	w.WriteHeader(http.StatusNoContent)
}

func (n *Node) replaceBlockchain(newBlocks []*blockchain.Block) {
	newChain := blockchain.NewBlockchainWithBlocks(newBlocks)

	if len(newChain.Blocks) > len(n.Blockchain.Blocks) && newChain.IsValid() {
		n.Blockchain = newChain
	}
}

func (n *Node) SynchronizeBlockchain() {
	for _, peer := range n.Peers {
		resp, err := http.Get(fmt.Sprintf("http://%s/blocks", peer))
		if err != nil {
			continue
		}

		var blocks []*blockchain.Block
		err = json.NewDecoder(resp.Body).Decode(&blocks)
		resp.Body.Close()
		if err != nil {
			continue
		}

		n.replaceBlockchain(blocks)
	}
}

func (n *Node) AddPKIRecord(w http.ResponseWriter, r *http.Request) {
	var pkiReq pki.Request
	err := json.NewDecoder(r.Body).Decode(&pkiReq)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	success, record := pki.Register(pkiReq.PublicKey, pkiReq.Signature, pkiReq.Address)
	if success {
		n.Blockchain.AddRecord(record)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Registered successfully"))
	} else {
		http.Error(w, record, http.StatusBadRequest)
	}
}
