package gophchain

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "gophchain.db"
const blocksBucket = "blocksBucket"

// Blockchain structure
type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(address string) (*Blockchain, error) {
	if dbExists() == false {
		return nil, fmt.Errorf("No existing blockchain found. Create one first.")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("1"))
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Blockchain{tip, db}, nil
}

// CreateBlockchain Returns a initial blockchain structure
func CreateBlockchain(address string) (*Blockchain, error) {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			coinbaseTx := NewCoinbaseTX(address, "Initial data")
			genesis := NewGenesisBlock(coinbaseTx)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			genesisSerialized, err := genesis.Serialize()
			if err != nil {
				return err
			}

			if err = b.Put(genesis.Hash, genesisSerialized); err != nil {
				return err
			}

			if err = b.Put([]byte("1"), genesis.Hash); err != nil {
				return err
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil
	})

	return &Blockchain{tip, db}, nil
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	var lastHash []byte
	if err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("1"))
		return nil
	}); err != nil {
		return err
	}

	newBlock := NewBlock(transactions, lastHash)
	err := bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serializedBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		if err = b.Put(newBlock.Hash, serializedBlock); err != nil {
			return err
		}

		if err = b.Put([]byte("1"), newBlock.Hash); err != nil {
			return err
		}

		bc.tip = newBlock.Hash
		return nil
	})

	return err
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) ([]Transaction, error) {
	var unspentTX []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block, err := bci.Next()
		if err != nil {
			return nil, err
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.id)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTX = append(unspentTX, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOuputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTX, nil
}

func (bc *Blockchain) FindUTXO(address string) ([]TxOutput, error) {
	var UTXO []TxOutput
	unspentTransactions, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXO = append(UTXO, out)
			}
		}
	}

	return UTXO, nil
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
