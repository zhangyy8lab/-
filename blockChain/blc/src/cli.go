package v1

import (
	"fmt"
	"os"
)

const Usage = `
	addBlock --data DATA "add a block"
	printChain "print block Chain"
`

type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) addBlock(data []byte) error {
	cli.Bc = GetBlockchainInstance()
	return cli.Bc.AddBlock(data)
}

func (cli *CLI) printChain() {
	cli.Bc = GetBlockchainInstance()
	it := cli.Bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("cur_Hash: %x\n", block.Hash)
		fmt.Println("blockData:", string(block.Data))

		if len(it.currentHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println(Usage)
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "addBlock":
		if len(os.Args) > 3 && os.Args[2] == "--data" {
			data := os.Args[3]
			if data == "" {
				fmt.Printf("data not be empty")
				os.Exit(1)
			}

			if err := cli.addBlock([]byte(data)); err != nil {
				fmt.Printf("add block failed, err:%v\n", err)
			}
		}
	case "printChain":
		cli.printChain()
	default:
		fmt.Println(Usage)
	}
}
