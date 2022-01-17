package balance

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetCurrent(client *ethclient.Client, address string) (*big.Int, error) {
	account := common.HexToAddress(address)

	//Fetch balance from the latest block
	balance, err := client.BalanceAt(context.Background(), account, nil)

	return balance, err
}

func GetPending(client *ethclient.Client, address string) (*big.Int, error) {
	account := common.HexToAddress(address)

	//Fetch balance from the latest block
	balance, err := client.PendingBalanceAt(context.Background(), account)

	return balance, err
}
