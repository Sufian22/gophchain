package gophchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

// Block structure
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

// NewGenesisBlock Returns the initial blockchain Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// NewBlock Creates a new Block with the proof of work needed
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.id)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func Deserialize(d []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}

	return &block, nil
}
