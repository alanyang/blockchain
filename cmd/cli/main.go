package main

import (
	"blockchain/core"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	addChainCmd := flag.NewFlagSet("add", flag.ExitOnError)
	eachChainCmd := flag.NewFlagSet("each", flag.ExitOnError)

	addChainData := addChainCmd.String("name", "", "block name")

	if len(os.Args) < 2 {
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		addChainCmd.Parse(os.Args[2:])
	case "each":
		eachChainCmd.Parse(os.Args[2:])
	default:
		os.Exit(0)
	}

	c, err := core.NewBlockChain()
	if err != nil {
		log.Fatal(err)
	}
	defer c.DB.Close()

	if addChainCmd.Parsed() {
		fmt.Println(c.AddBlock(core.NewBlock(*addChainData)))
	}

	if eachChainCmd.Parsed() {
		fmt.Println(c.String())
	}
}
