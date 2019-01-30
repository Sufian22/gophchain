package gophchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const targetBits = 16
const maxNonce = math.MaxInt64

// ProofOfWork Structure
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork Returns a new ProofOfWork structure with the block provided
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{block: b, target: target}
}

// Run It's the core function of our PoW algorithm
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

// Validate Returns if the proof of work is valid
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.Data,
			pow.block.PrevBlockHash,
			intToHex(pow.block.Timestamp),
			intToHex(int64(targetBits)),
			intToHex(nonce),
		},
		[]byte{},
	)
	return data
}

func intToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}
