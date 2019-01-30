package gophchain

import (
	"github.com/boltdb/bolt"
)

const dbFile = "gophchain.db"
const blocksBucket = "blocksBucket"

// Blockchain structure
type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}

// NewBlockchain Returns a initial blockchain structure
func NewBlockchain() (*Blockchain, error) {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			genesisSerialized, err := genesis.Serialize()
			if err != nil {
				return err
			}

			err = b.Put(genesis.Hash, genesisSerialized)
			if err != nil {
				return err
			}

			err = b.Put([]byte("1"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil
	})

	return &Blockchain{tip, db}, nil
}

// AddBlock Adds new block to the blockchain
func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte
	if err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))

		return nil
	}); err != nil {
		return err
	}

	newBlock := NewBlock(data, lastHash)
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serializedBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		err = b.Put(newBlock.Hash, serializedBlock)
		if err != nil {
			return err
		}

		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			return err
		}

		bc.tip = newBlock.Hash
		return nil
	})

	return err
}

// BlockchainIterator for BoltDB
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.Db}
}

func (i *BlockchainIterator) Next() (*Block, error) {
	var block *Block
	if err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		var err error
		block, err = Deserialize(encodedBlock)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	i.currentHash = block.PrevBlockHash
	return block, nil
}
