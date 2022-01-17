package wallet

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"log"
	"math/big"
)

type Erc20Token struct {
  Address string
}

func New() (privateKeyString string, publicKeyAddress string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("error casting public key to ECDSA: ", err)
	}

	publicKeyAddress = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return  privateKeyString, publicKeyAddress
}

func NewKeystore(password string) string {
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	return account.Address.Hex()
}

func RetrieveKeystore(password string, filepath string) (publicKeyString, privateKeyString string) {
	jsonBytes, err := ioutil.ReadFile(filepath)
	account, err := keystore.DecryptKey(jsonBytes, password)

	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := keystore.DecryptKey(jsonBytes, password)
	privateKeyBytes := crypto.FromECDSA(privateKey.PrivateKey)
	privateKeyString = hexutil.Encode(privateKeyBytes)[2:]

	if err != nil {
		log.Fatal("Could not decrypt keystore: ", filepath)
	}

	return account.Address.Hex(), privateKeyString
}

func IsContractAddress(client *ethclient.Client,address string) bool {
	account := common.HexToAddress(address)

	bytecode, err := client.CodeAt(context.Background(), account, nil)

	if err != nil {
		log.Fatal("Could not determine if account has code: ",err)
	}

	isContract := len(bytecode) > 1

	return isContract
}

func Transfer(client *ethclient.Client, hexKey string, to string, amount *big.Int) *types.Transaction {
	privateKey , err := crypto.HexToECDSA(hexKey)

	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA: ", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(to)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Could not determine Network ID: ",err)
	}

	gasLimit := uint64(21000)
	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal("Could not determine Gas tip cap: ",err)
	}
	gasMaxFee, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Could not determine Gas tip cap: ",err)
	}

	transaction := &types.DynamicFeeTx{
		Nonce: nonce,
		ChainID: chainID,
		Gas: gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasMaxFee,
		Value: amount,
		To: &toAddress,
	}

	signedTx := types.MustSignNewTx(privateKey, types.LatestSignerForChainID(chainID), transaction)
	if err != nil {
		log.Fatal("Could not Sign Transaction: ",err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("Cannot send Transaction: ",err)
	}

	return signedTx
}

func (token Erc20Token) Transfer(client *ethclient.Client, hexKey string, to string, amount *big.Int) *types.Transaction {
	privateKey , err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA: ", err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(to)
	tokenAddress := common.HexToAddress(token.Address)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal("Could not determine Network ID: ",err)
	}

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal("Could not determine Gas tip cap: ",err)
	}
	gasMaxFee, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("Could not determine Gas tip cap: ",err)
	}

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal("Gas limited: ", err)
	}

	transaction := &types.DynamicFeeTx{
		Nonce: nonce,
		ChainID: chainID,
		Gas: gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasMaxFee,
		To: &tokenAddress,
		Data: data,
	}

	signedTx := types.MustSignNewTx(privateKey, types.LatestSignerForChainID(chainID), transaction)
	if err != nil {
		log.Fatal("Could not Sign Transaction: ",err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("Cannot send Transaction: ",err)
	}

	return signedTx
}