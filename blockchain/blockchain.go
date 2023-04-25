package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Blockchain struct {
	Blocks  []*Block
	memPool []string
}

const memPoolCapacity = 10

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{newGenesisBlock()}, []string{}}
}

func NewBlockchainWithBlocks(newBlocks []*Block) *Blockchain {
	return &Blockchain{newBlocks, []string{}}
}

func newGenesisBlock() *Block {
	return newBlock([]string{"Genesis Block"}, []byte{})
}

func (bc *Blockchain) AddBlock(records []string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(records, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) AddRecord(record string) {
	bc.memPool = append(bc.memPool, record)

	if len(bc.memPool) >= memPoolCapacity {
		bc.minePendingRecords()
	}
}

func (bc *Blockchain) minePendingRecords() {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(bc.memPool, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)

	bc.memPool = []string{}
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		hash := sha256.Sum256(currentBlock.prepareData())
		if !bytes.Equal(currentBlock.Hash, hash[:]) {
			return false
		}

		if !bytes.Equal(currentBlock.PrevBlockHash, prevBlock.Hash) {
			return false
		}
	}
	return true
}
