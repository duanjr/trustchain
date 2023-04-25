package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
)

type Account struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Address    string `json:"address"`
}

func main() {
	accounts := make([]Account, 3000)

	// 生成3000个账户
	for i := 0; i < 3000; i++ {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Error generating key: %v", err)
		}

		privateKeyBytes := crypto.FromECDSA(privateKey)
		privateKeyHex := hex.EncodeToString(privateKeyBytes)

		publicKey := privateKey.Public().(*ecdsa.PublicKey)
		publicKeyBytes := crypto.FromECDSAPub(publicKey)
		publicKeyHex := hex.EncodeToString(publicKeyBytes)

		address := fmt.Sprintf("192.168.0.%d", i+1)

		accounts[i] = Account{
			PrivateKey: privateKeyHex,
			PublicKey:  publicKeyHex,
			Address:    address,
		}
	}

	// 将账户持久化为JSON文件
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling accounts: %v", err)
	}

	err = ioutil.WriteFile("accounts.json", data, 0644)
	if err != nil {
		log.Fatalf("Error writing accounts to file: %v", err)
	}

	fmt.Println("Accounts saved to accounts.json")
}
