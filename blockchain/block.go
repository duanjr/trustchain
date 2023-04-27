package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math"
	"math/big"
	"strings"
	"time"
)

const targetBits = 16
const maxNonce = math.MaxInt64

type Block struct {
	Timestamp           int64
	Records             []string
	PrevBlockHash       []byte
	Hash                []byte
	Nonce               int
	PkiRootHash         common.Hash
	DirectTrustRootHash common.Hash
	CompTrustRootHash   common.Hash
}

func newBlock(records []string, prevBlockHash []byte, pkiRootHash common.Hash,
	directTrustRootHash common.Hash, compTrustRootHash common.Hash) *Block {
	block := &Block{time.Now().Unix(), records, prevBlockHash,
		[]byte{}, 0, pkiRootHash, directTrustRootHash, compTrustRootHash}
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
