package gophchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

const subsidy = 50

// TXOutput represents a transaction output
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (txo *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return txo.ScriptPubKey == unlockingData
}

// TxInput represents a transaction input
type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (txi *TxInput) CanUnlockOuputWith(unlockingData string) bool {
	return txi.ScriptSig == unlockingData
}

// Transaction represents a Bitcoin transaction
type Transaction struct {
	id   []byte
	Vin  []TxInput
	Vout []TxOutput
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data != "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := []TxInput{{Txid: []byte{}, Vout: -1, ScriptSig: data}}
	txOut := []TxOutput{{Value: subsidy, ScriptPubKey: to}}
	tx := Transaction{id: nil, Vin: txIn, Vout: txOut}
	tx.SetID()

	return &tx
}

// IsCoinbase checks whether the transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := json.NewEncoder(&encoded)
	if err := enc.Encode(tx); err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.id = hash[:]
	return nil
}
