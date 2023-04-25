package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
	"strings"
	"time"
)

const targetBits = 20
const maxNonce = math.MaxInt64

type Block struct {
	Timestamp     int64
	Records       []string
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func newBlock(records []string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), records, prevBlockHash, []byte{}, 0}
	block.mine()

	return block
}

func (b *Block) mine() {
	var hashInt big.Int
	var hash [32]byte

	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	for b.Nonce = 0; b.Nonce < maxNonce; b.Nonce++ {
		data := b.prepareData()
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			b.Hash = hash[:]
			break
		}
	}
}

func (b *Block) prepareData() []byte {
	data := bytes.Join(
		[][]byte{
			b.PrevBlockHash,
			[]byte(strings.Join(b.Records, ";")),
			intToHex(b.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(int64(b.Nonce)),
		},
		[]byte{},
	)

	return data
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
