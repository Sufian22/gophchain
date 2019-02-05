package gophchain

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// CLI process the command line arguments
type CLI struct{}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() error {
	cli.validateArgs()

	// addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Blockchain address")
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "Blockchain address")

	var err error
	switch os.Args[1] {
	case "createblockchain":
		err = createBlockchainCmd.Parse(os.Args[2:])
	case "getbalance":
		err = getBalanceCmd.Parse(os.Args[2:])
	case "printchain":
		err = printChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if err != nil {
		return err
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	return nil
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address creates a blockchain with the specified address")
	fmt.Println("  getbalance -address gets the balance for the blockchain with the specified address")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) createBlockchain(address string) {
	bc, err := CreateBlockchain(address)
	if err != nil {
		fmt.Printf("Could not create blockchain with address: %s", address)
		os.Exit(1)
	}
	bc.Db.Close()
	fmt.Println("Blockchain created")
}

func (cli *CLI) getBalance(address string) {
	bc, err := NewBlockchain(address)
	if err != nil {
		fmt.Printf("Could not create blockchain with address: %s", address)
		os.Exit(1)
	}
	bc.Db.Close()

	balance := 0
	UTXOs, err := bc.FindUTXO(address)
	if err != nil {
		fmt.Printf("Could not get transactions from the blockchain with address: %s", address)
		os.Exit(1)
	}

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) printChain() error {
	bc, err := CreateBlockchain("")
	if err != nil {
		fmt.Printf("Could not create blockchain with address: %s", "")
		os.Exit(1)
	}
	bc.Db.Close()

	iterator := bc.Iterator()
	for {
		block, err := iterator.Next()
		if err != nil {
			return err
		}

		pow := NewProofOfWork(block)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil
}
