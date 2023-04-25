package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

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

	// 打印第一个账户
	if len(accounts) > 0 {
		fmt.Printf("First account:\nPrivate Key: %s\nPublic Key: %s\nAddress: %s\n",
			accounts[0].PrivateKey, accounts[0].PublicKey, accounts[0].Address)
	}
}
