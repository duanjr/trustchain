package blockchain

import (
	"bytes"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/trie"
	"log"
	"math"
	"strconv"
)

type Blockchain struct {
	Blocks          []*Block
	memPool         []string
	PkiTrie         *trie.Trie
	DirectTrustTrie *trie.Trie
	CompTrustTrie   *trie.Trie
	Id2DT           map[string]map[string]float64
	c               float64
	AddressList     *[]string
}

const memPoolCapacity = 30001

func NewBlockchain() *Blockchain {
	pkiMemdb := memorydb.New()
	pkiTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(pkiMemdb))
	directTrustMemdb := memorydb.New()
	directTrustTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(directTrustMemdb))
	compTrustMemdb := memorydb.New()
	compTrustTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(compTrustMemdb))
	return &Blockchain{
		Blocks:          []*Block{newGenesisBlock()},
		memPool:         []string{},
		PkiTrie:         pkiTrie,
		DirectTrustTrie: directTrustTrie,
		CompTrustTrie:   compTrustTrie,
		Id2DT:           make(map[string]map[string]float64),
		c:               1,
		AddressList:     new([]string),
	}
}

func NewBlockchainWithBlocks(newBlocks []*Block) *Blockchain {
	pkiMemdb := memorydb.New()
	pkiTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(pkiMemdb))
	directTrustMemdb := memorydb.New()
	directTrustTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(directTrustMemdb))
	compTrustMemdb := memorydb.New()
	compTrustTrie, _ := trie.New(common.Hash{}, trie.NewDatabase(compTrustMemdb))
	return &Blockchain{
		Blocks:          newBlocks,
		memPool:         []string{},
		PkiTrie:         pkiTrie,
		DirectTrustTrie: directTrustTrie,
		CompTrustTrie:   compTrustTrie,
		Id2DT:           make(map[string]map[string]float64),
		c:               1,
		AddressList:     new([]string),
	}
}

func newGenesisBlock() *Block {
	return newBlock([]string{"Genesis Block"}, []byte{}, common.Hash{}, common.Hash{}, common.Hash{})
}

func (bc *Blockchain) AddBlock(records []string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(records, prevBlock.Hash, bc.PkiTrie.Hash(), bc.DirectTrustTrie.Hash(), bc.CompTrustTrie.Hash())
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) AddRecord(record string) {
	bc.memPool = append(bc.memPool, record)

	if len(bc.memPool) >= memPoolCapacity {
		bc.minePendingRecords()
	}
}

func (bc *Blockchain) minePendingRecords() {
	bc.CalculateAllCompTrust()
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := newBlock(bc.memPool, prevBlock.Hash, bc.PkiTrie.Hash(), bc.DirectTrustTrie.Hash(), bc.CompTrustTrie.Hash())
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

func (bc *Blockchain) CompTrust(addressI string, addressJ string) float64 {
	DT, ok := bc.Id2DT[addressI][addressJ]
	if !ok {
		return 0
	}

	iT := 0.0
	m := 0
	DtSum := 0.0

	for k, DtIk := range bc.Id2DT[addressI] {
		if k != addressJ && DtIk > 0 {
			m++
			DtSum += DtIk
			DtKj, ok := bc.Id2DT[k][addressJ]
			if ok {
				iT += DtIk * DtKj
			}
		}
	}
	if m > 0 {
		iT /= DtSum
	}

	alpha := 0.0
	mu := 0.0
	sigma := 0.0

	for k, DtIk := range bc.Id2DT[addressI] {
		if k != addressJ && DtIk > 0 {
			m++
			DtKj, ok := bc.Id2DT[k][addressJ]
			if !ok {
				DtKj = 0
			}
			sigma += (DtKj - iT) * (DtKj - iT)
		}
	}

	mu = float64(m) / (float64(m) + bc.c)
	if m > 0 {
		sigma = math.Sqrt(sigma / float64(m))
		sigma = 1 / (1 + sigma)
	}

	alpha = (mu + sigma) / 4

	result := (1-alpha)*DT + alpha*iT

	return result
}

func (bc *Blockchain) CalculateAllCompTrust() {
	for _, i := range *bc.AddressList {
		for _, j := range *bc.AddressList {
			if i != j {
				compTrust := bc.CompTrust(i, j)
				err := bc.CompTrustTrie.TryUpdate([]byte(i+"&"+j), []byte(strconv.FormatFloat(compTrust, 'f', 6, 64)))
				if err != nil {
					log.Fatalf("Error updating CompTrustTrie: %v", err)
				}
			}
		}
	}
}
