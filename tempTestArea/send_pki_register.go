package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"net/http"
)

func sendRegisterRequest(publicKey, signature, address string) {
	data := map[string]string{
		"publicKey": publicKey,
		"signature": signature,
		"address":   address,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/pki/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", body)
}

type Account struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Address    string `json:"address"`
}

func main() {
	data, err := ioutil.ReadFile("accounts.json")
	if err != nil {
		log.Fatalf("Error reading accounts from file: %v", err)
	}

	var accounts []Account
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		log.Fatalf("Error unmarshalling accounts: %v", err)
	}

	fmt.Printf("Read %d accounts from accounts.json\n", len(accounts))
	privateKeyBytes, err := hex.DecodeString(accounts[0].PrivateKey)
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	publicKey := accounts[0].PublicKey
	address := accounts[0].Address

	message := "register" + address
	hashedMessage := crypto.Keccak256Hash([]byte(message))
	signature, err := crypto.Sign(hashedMessage.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("Error signing message: %v", err)
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	fmt.Printf("PublicKey: %s\n", publicKey)
	fmt.Printf("Signature: %s\n", signatureBase64)
	fmt.Printf("Address: %s\n", address)

	sendRegisterRequest(publicKey, signatureBase64, address)
	sendRegisterRequest(publicKey, signatureBase64, address)
}
