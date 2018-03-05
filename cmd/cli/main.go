package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alanyang/blockchain/core"
)

func main() {

	addChainCmd := flag.NewFlagSet("createblock", flag.ExitOnError)
	eachChainCmd := flag.NewFlagSet("eachblock", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getkeyhash := flag.NewFlagSet("getkeyhash", flag.ExitOnError)

	createBlockAddress := addChainCmd.String("address", "", "")
	getBalanceAddress := getBalanceCmd.String("address", "", "")
	getkeyhashAddress := getkeyhash.String("address", "", "")

	if len(os.Args) < 2 {
		os.Exit(0)
	}

	switch os.Args[1] {
	case "createblock":
		addChainCmd.Parse(os.Args[2:])
	case "eachblock":
		eachChainCmd.Parse(os.Args[2:])
	case "createwallet":
		createWalletCmd.Parse(os.Args[2:])
	case "getbalance":
		getBalanceCmd.Parse(os.Args[2:])
	case "getkeyhash":
		getkeyhash.Parse(os.Args[2:])
	default:
		os.Exit(0)
	}

	if addChainCmd.Parsed() {
		c, err := core.NewBlockChain()
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		b, ok := c.AddBlock(core.NewBlock(c.Last, []*core.Transaction{core.NewCoinbaseTransaction(*createBlockAddress, "")}))
		if !ok {
			fmt.Println("Invalid block, rejected!")
			os.Exit(0)
		}
		fmt.Println(b)
	}

	if eachChainCmd.Parsed() {
		c, err := core.NewBlockChain()
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		fmt.Println(c.String())
	}

	if createWalletCmd.Parsed() {
		fmt.Println(core.NewWallet())
	}

	if getBalanceCmd.Parsed() {
		c, err := core.NewBlockChain()
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()
		fmt.Println(c.GetBalance(*getBalanceAddress))
	}

	if getkeyhash.Parsed() {
		fmt.Printf("%X\n", core.PubKeyHashFromAddress(*getkeyhashAddress))
	}
}
