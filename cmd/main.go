package main

import (
	"log"

	"github.com/sufian22/gophchain"
)

func main() {
	bc, err := gophchain.NewBlockchain()
	if err != nil {
		log.Fatalf("Could not initialize blockchain: %s", err.Error())
	}
	defer bc.Db.Close()

	cli := gophchain.CLI{Bc: bc}
	cli.Run()
}
