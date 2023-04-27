package node

import (
	"encoding/json"
	"fmt"
	"github.com/duanjr/trustchain/blockchain"
	"github.com/duanjr/trustchain/pki"
	"github.com/duanjr/trustchain/trust"
	"net/http"
	"strconv"
)

type Node struct {
	Blockchain *blockchain.Blockchain
	Peers      []string
}

func NewNode() *Node {
	res := &Node{blockchain.NewBlockchain(), []string{}}
	pki.Initialize(res.Blockchain.PkiTrie)
	trust.Initialize(res.Blockchain.DirectTrustTrie, res.Blockchain.PkiTrie, res.Blockchain.CompTrustTrie,
		res.Blockchain.Id2DT, res.Blockchain.AddressList)
	return res
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
	var pkiReq pki.RegisterRequest
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

func (n *Node) UpdatePKIRecord(w http.ResponseWriter, r *http.Request) {
	var updateReq pki.UpdateRequest
	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	success, msg := pki.Update(updateReq.PublicKey1, updateReq.Signature1, updateReq.PublicKey2, updateReq.Signature2, updateReq.Address)
	if success {
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Updated successfully"))
	} else {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

func (n *Node) QueryPKIRecord(w http.ResponseWriter, r *http.Request) {
	var pkiReq pki.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&pkiReq)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	publicKey, err := pki.Query(pkiReq.Address)
	if err != nil {
		http.Error(w, "No such identity", http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(publicKey))
	}
}
func (n *Node) RevokePKIRecord(w http.ResponseWriter, r *http.Request) {
	var pkiReq pki.RevokeRequest
	err := json.NewDecoder(r.Body).Decode(&pkiReq)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	success, record := pki.Revoke(pkiReq.PublicKey, pkiReq.Signature, pkiReq.Address)
	if success {
		n.Blockchain.AddRecord(record)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Revoked successfully"))
	} else {
		http.Error(w, record, http.StatusBadRequest)
	}
}

func (n *Node) TrustSubmitRecord(w http.ResponseWriter, r *http.Request) {
	var req trust.SubmitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	err = trust.Submit(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (n *Node) DirectTrustQueryRecord(w http.ResponseWriter, r *http.Request) {
	var req trust.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	trustValue, err := trust.QueryDirect(req)
	if err != nil {
		http.Error(w, "No such identity", http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(trustValue))
	}
}

func (n *Node) CompTrustQuery(w http.ResponseWriter, r *http.Request) {
	var req trust.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	trustValue, err := trust.QueryComp(req)
	if err != nil {
		http.Error(w, "No such identity", http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(trustValue))
	}
}

func (n *Node) CalcCompTrustQuery(w http.ResponseWriter, r *http.Request) {
	var req trust.QueryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	trustValue := n.Blockchain.CompTrust(req.AddressI, req.AddressJ)

	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(strconv.FormatFloat(trustValue, 'f', -1, 64)))
}
