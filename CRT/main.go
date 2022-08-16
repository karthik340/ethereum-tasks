package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/e5e9354dc8f142ed8d0cacb4702a0ecc")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	var previousHash string = ""
	var previousBlock *big.Int = big.NewInt(0)
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			BlockNumber := header.Number
			currentHash := header.Hash().Hex()
			// fmt.Println(BlockNumber, currentHash)

			if previousBlock.Cmp(BlockNumber) == 0 && previousHash != currentHash {
				fmt.Println("forked block number : ", BlockNumber)
			}
			previousHash = currentHash
			previousBlock = BlockNumber
		}
	}
}
