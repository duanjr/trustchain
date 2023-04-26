package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
)

type Blockchain struct {
	Blocks  []*Block
	memPool []string
	PkiTrie *trie.Trie
}

const memPoolCapacity = 3001

func NewBlockchain() *Blockchain {
	pkiMemdb := memorydb.New()
	pkiTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(pkiMemdb))
	return &Blockchain{[]*Block{newGenesisBlock()}, []string{}, pkiTrie}
}

func NewBlockchainWithBlocks(newBlocks []*Block) *Blockchain {
	pkiMemdb := memorydb.New()
	pkiTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(pkiMemdb))
	return &Blockchain{newBlocks, []string{}, pkiTrie}
}

func newGenesisBlock() *Block {
	return newBlock([]string{"Genesis Block"}, []byte{}, common.Hash{})
}

func (bc *Blockchain) AddBlock(records []string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(records, prevBlock.Hash, bc.PkiTrie.Hash())
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
	newBlock := newBlock(bc.memPool, prevBlock.Hash, bc.PkiTrie.Hash())
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
