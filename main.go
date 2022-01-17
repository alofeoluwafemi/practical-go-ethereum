package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/slim12kg/practical-go-ethereum/account/balance"
	"log"
	"math/big"
)

const (
	address string = "0xc90d4b1d68178Bc8cA417610DeddE7773409e097"
	rpcUrl  string = "https://rinkeby.infura.io/v3/6f20b652724a4fb198c8f3d3041d6952"
)

var (
	client *ethclient.Client
	err    error
)

func init() {
	client, err = ethclient.Dial(rpcUrl)

	if err != nil {
		log.Fatalf("cannot connect to %q", rpcUrl)
	}
}

func main() {
	// Generate Wallet
	// Private Key ccbe9d382aae335c579c3cb0ef8a1c03e7d758cfd49aaf0cfb9b54249bbad769
	// Public Key 0xc90d4b1d68178Bc8cA417610DeddE7773409e097
	// privateKey, publicKey := wallet.New()

	// Get Wallet Details from Keystore
	//publicAddress, privateAddress := wallet.RetrieveKeystore("Abcd1234","./keystore/UTC--2022-01-12T12-31-51.507222000Z--71d74b06af345c808faf7a3e1d8b42126cb100bb")
	//
	//fmt.Println(publicAddress)
	//fmt.Println(privateAddress)

	// GetBalance
	account := "0xc90d4b1d68178Bc8cA417610DeddE7773409e097"
	bal, err := balance.GetCurrent(client,account)
	if err != nil {
		log.Fatalf("Could not fetch balance of %q", account)
	}
	fmt.Println(bal)
	fmt.Println(weiToEther(bal))

	// Transfer Balance
	//	tx := wallet.Transfer(client, "ccbe9d382aae335c579c3cb0ef8a1c03e7d758cfd49aaf0cfb9b54249bbad769", "0x3472059945ee170660a9A97892a3cf77857Eba3A", big.NewInt(1000000000000))
	//	fmt.Println(tx.Hash())
	//	fmt.Println(tx.Cost())

	//link := wallet.Erc20Token{Address: "0x01be23585060835e02b77ef475b0cc51aa1e0709"}
	//tx := link.Transfer(client,"ccbe9d382aae335c579c3cb0ef8a1c03e7d758cfd49aaf0cfb9b54249bbad769","0x3472059945ee170660a9A97892a3cf77857Eba3A",big.NewInt(2000000000000000000))
	//
	//fmt.Println(tx.Hash())
	//fmt.Println(tx.Cost())

	// Listen to transactions using wss
	//query.SubscribeHeader(client)

	//daiContractAddress := common.HexToAddress("0xc60bea04bce65593bde3d872322d17311b78a991")
	//nDai, err := dai.NewDai(daiContractAddress, client)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	//opt := &bind.CallOpts{
	//	Context: context.Background(),
	//	BlockNumber: nil,
	//	From: daiContractAddress,
	//}
	//	bal, err := nDai.BalanceOf(&bind.CallOpts{},common.HexToAddress(address))
	//	if err != nil {
	//		log.Fatal("Could not check balance: ",err)
	//	}
	//	fmt.Println(bal)
}

func etherToWei(etherVal *big.Float) *big.Int {
	coerce := new(big.Int)
	coerce.SetString(etherVal.String(), 10)

	return new(big.Int).Mul(coerce, big.NewInt(100000000000000000))
}

// weiToEther converts weiVal to ether value equivalent
func weiToEther(weiVal *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(weiVal), big.NewFloat(params.Ether))
}
