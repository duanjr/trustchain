package pki

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
)

var memDB = memorydb.New()
var Trie, _ = trie.New(common.Hash{}, trie.NewDatabase(memDB))

type Request struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func Register(publicKey, signature, address string) (bool, string) {
	if publicKey == "" || signature == "" || address == "" {
		return false, "Missing values"
	}

	msg := "register" + address
	hashedMessage := crypto.Keccak256Hash([]byte(msg))
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, "Invalid signature format"
	}

	recoveredPubkey, err := crypto.SigToPub(hashedMessage.Bytes(), sig)
	if err != nil {
		return false, "Unable to recover public key from signature"
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	pubkeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, "Invalid public key format"
	}

	pub, err := crypto.UnmarshalPubkey(pubkeyBytes)
	computedAddress := crypto.PubkeyToAddress(*pub)

	if recoveredAddress == computedAddress {
		if val, _ := Trie.TryGet([]byte(address)); val != nil {
			return false, "Address registered"
		}

		record := fmt.Sprintf("PKI:Register:PublicKey:%s:Address:%s", publicKey, address)
		Trie.Update([]byte(address), pubkeyBytes)
		return true, record
	} else {
		return false, "Wrong argument"
	}
}
