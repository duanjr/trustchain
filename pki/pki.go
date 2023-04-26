package pki

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
)

var Trie *trie.Trie

func Initialize(t *trie.Trie) {
	Trie = t
}

type RegisterRequest struct {
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

type UpdateRequest struct {
	PublicKey1 string `json:"publicKey1"`
	Signature1 string `json:"signature1"`
	PublicKey2 string `json:"publicKey2"`
	Signature2 string `json:"signature2"`
	Address    string `json:"address"`
}

func Update(publicKey1, signature1, publicKey2, signature2, address string) (bool, string) {
	if publicKey1 == "" || signature1 == "" || publicKey2 == "" || signature2 == "" || address == "" {
		return false, "Missing values"
	}

	msg1 := publicKey2
	hashedMessage1 := crypto.Keccak256Hash([]byte(msg1))
	sig1, err := base64.StdEncoding.DecodeString(signature1)
	if err != nil {
		return false, "Invalid signature1 format"
	}

	recoveredPubkey1, err := crypto.SigToPub(hashedMessage1.Bytes(), sig1)
	if err != nil {
		return false, "Unable to recover public key from signature1"
	}

	recoveredAddress1 := crypto.PubkeyToAddress(*recoveredPubkey1)
	pubkeyBytes1, err := hex.DecodeString(publicKey1)
	if err != nil {
		return false, "Invalid publicKey1 format"
	}

	pub1, err := crypto.UnmarshalPubkey(pubkeyBytes1)
	computedAddress1 := crypto.PubkeyToAddress(*pub1)

	msg2 := "register" + address
	hashedMessage2 := crypto.Keccak256Hash([]byte(msg2))
	sig2, err := base64.StdEncoding.DecodeString(signature2)
	if err != nil {
		return false, "Invalid signature2 format"
	}

	recoveredPubkey2, err := crypto.SigToPub(hashedMessage2.Bytes(), sig2)
	if err != nil {
		return false, "Unable to recover public key from signature2"
	}

	recoveredAddress2 := crypto.PubkeyToAddress(*recoveredPubkey2)
	pubkeyBytes2, err := hex.DecodeString(publicKey2)
	if err != nil {
		return false, "Invalid publicKey2 format"
	}

	pub2, err := crypto.UnmarshalPubkey(pubkeyBytes2)
	computedAddress2 := crypto.PubkeyToAddress(*pub2)

	if computedAddress1 == recoveredAddress1 && computedAddress2 == recoveredAddress2 {
		oldPubKey, _ := Trie.TryGet([]byte(address))
		if !bytes.Equal(oldPubKey, pubkeyBytes1) {
			return false, "Wrong old public key"
		}

		record := fmt.Sprintf("PKI:Update:PublicKey:%s:Address:%s", publicKey2, address)
		Trie.Update([]byte(address), pubkeyBytes2)
		return true, record
	} else {
		return false, "Wrong signatures"
	}
}

type QueryRequest struct {
	Address string `json:"address"`
}

func Query(address string) (string, error) {
	if address == "" {
		return "", errors.New("missing address")
	}

	pubkeyBytes, err := Trie.TryGet([]byte(address))
	if err != nil {
		return "", errors.New("no such identity")
	}

	return hex.EncodeToString(pubkeyBytes), nil
}

type RevokeRequest struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func Revoke(publicKey, signature, address string) (bool, string) {
	if publicKey == "" || signature == "" || address == "" {
		return false, "Missing values"
	}

	msg := "revoke" + address
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
		storedPublicKey, err := Query(address)
		if err != nil || storedPublicKey != publicKey {
			return false, "Wrong public key to revoke"
		}

		record := fmt.Sprintf("PKI:Revoke:PublicKey:%s:Address:%s", publicKey, address)
		Trie.Delete([]byte(address))
		return true, record
	} else {
		return false, "Wrong argument"
	}
}
