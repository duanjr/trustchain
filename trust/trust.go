package trust

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"math"
	"strconv"
	"strings"
	"time"
)

var Trie *trie.Trie
var PKITrie *trie.Trie
var id2DT map[string]map[string]float64
var addressList *[]string

func Initialize(t1 *trie.Trie, t2 *trie.Trie, m map[string]map[string]float64, a *[]string) {
	Trie = t1
	PKITrie = t2
	id2DT = m
	addressList = a
}

type SubmitRequest struct {
	AddressI   string  `json:"addressI"`
	AddressJ   string  `json:"addressJ"`
	TrustValue float64 `json:"trustValue"`
	Timestamp  int64   `json:"timestamp"`
	Signature  string  `json:"signature"`
}

func Submit(req SubmitRequest) error {
	if req.TrustValue > 1 || req.TrustValue < -1 {
		return errors.New("expected trustValue between 1 and -1")
	}

	currentTime := time.Now().Unix()
	if math.Abs(float64(req.Timestamp-currentTime)) > 40 {
		return errors.New("invalid timestamp")
	}

	msg := fmt.Sprintf("submit%s.%s.%f.%d", req.AddressI, req.AddressJ, req.TrustValue, req.Timestamp)
	hashedMessage := crypto.Keccak256Hash([]byte(msg))
	signature, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return errors.New("invalid signature")
	}

	addressRecover, err := crypto.Ecrecover(hashedMessage.Bytes(), signature)
	if err != nil {
		return errors.New("unable to recover address")
	}

	addressBytes := crypto.Keccak256(addressRecover[1:])[12:]
	address := hex.EncodeToString(addressBytes)
	if strings.ToLower("0x"+address) != strings.ToLower(req.AddressI) {
		return errors.New("wrong signature")
	}

	if _, err := PKITrie.TryGet([]byte(req.AddressI)); err != nil {
		return errors.New("Address " + req.AddressI + " is not registered")
	}

	id2DT[req.AddressI][req.AddressJ] = req.TrustValue
	Trie.Update([]byte(req.AddressI+req.AddressJ), []byte(strconv.FormatFloat(req.TrustValue, 'f', -1, 64)))

	if !contains(*addressList, req.AddressI) {
		*addressList = append(*addressList, req.AddressI)
	}
	if !contains(*addressList, req.AddressJ) {
		*addressList = append(*addressList, req.AddressJ)
	}

	return nil
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
