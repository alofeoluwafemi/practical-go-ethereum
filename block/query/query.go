package query

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func BlockHeader(client *ethclient.Client) *types.Header {
	header, err := client.HeaderByNumber(context.Background(), nil)

	if err != nil {
		log.Fatalln("Could not fetch block header information: ", err)
	}

	return header
}

func FullBlock(client *ethclient.Client) *types.Block{
	blockNumber := BlockHeader(client).Number

	block, err := client.BlockByNumber(context.Background(), blockNumber)

	if err != nil {
		log.Fatalln("Could not fetch block details: ", err)
	}

	return block
}

func SubscribeHeader(client *ethclient.Client)  {
	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())
		}
	}
}
