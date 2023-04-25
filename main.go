package main

import (
	"fmt"
	"github.com/duanjr/trustchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddRecord("Record 1")
	bc.AddRecord("Record 2")
	bc.AddRecord("Record 3")
	bc.AddRecord("Record 4")
	bc.AddRecord("Record 5")
	bc.AddRecord("Record 6")
	bc.AddRecord("Record 7")
	bc.AddRecord("Record 8")
	bc.AddRecord("Record 9")
	bc.AddRecord("Record 10")
	bc.AddRecord("Record 11")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Records)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
