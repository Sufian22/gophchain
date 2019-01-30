package gophchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Block structure
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

// NewGenesisBlock Returns the initial blockchain Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlock Creates a new Block with the proof of work needed
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
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

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func Deserialize(d []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(d))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}
