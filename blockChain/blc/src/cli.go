package src

import (
	"fmt"
	"os"
	"strconv"
)

const Usage = `
	create "create blockchain"
	delete "delete blockchain DB"
	print "print block Chain"
	balance <address> "get address balance"
	send <from> <to> <amount> "transaction"
	createWallet "create wallet"
`

// Run 启动
func (cli *CLI) Run() {
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Println(Usage)
		return
	}

	switch cmds[1] {
	case "create":
		cli.createBlockChain()
		fmt.Println(Usage)
	case "delete":
		cli.deleteBlockChainDB()
	case "print":
		cli.printChain()
	case "balance":
		if len(cmds) != 3 {
			fmt.Println("input args error!")
			fmt.Println(Usage)
			return
		}
		address := cmds[2]
		fmt.Println("address:", address)
		if address == "" {
			fmt.Printf("address not be empty")
			fmt.Println(Usage)
			return
		}
		cli.getBalance(address)
	case "send":
		fmt.Println("len(cmds):", len(cmds))
		fmt.Println("cmds:", cmds)
		if len(cmds) != 5 {
			fmt.Println("input args error!")
			fmt.Println(Usage)
			return
		}
		from := cmds[2]
		to := cmds[3]
		amount, err := strconv.ParseFloat(cmds[4], 64)
		if err != nil {
			fmt.Printf("amount is not a number")
			os.Exit(1)
		}
		cli.send(from, to, amount)

	case "createWallet":
		cli.createWallet()
	default:
		fmt.Println(Usage)
	}
}
