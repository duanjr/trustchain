package blockchain

type Blockchain struct {
	Blocks  []*Block
	memPool []string
}

const memPoolCapacity = 10

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{newGenesisBlock()}, []string{}}
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
