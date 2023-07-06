package erc20factory

import (
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ERC20Factory struct {
	ContractABI abi.ABI
	FactoryAddr common.Address
	Client      *ethclient.Client
}

func NewERC20Factory(ipcPath string, contractAbi string, factoryAddress string) (*ERC20Factory, error) {
	client, err := ethclient.Dial(ipcPath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractAbi))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	return &ERC20Factory{
		ContractABI: parsedABI,
		FactoryAddr: common.HexToAddress(factoryAddress),
		Client:      client,
	}, nil
}

func (f *ERC20Factory) CreateERC20(name, symbol, protocol string, decimals uint8, owner, admin common.Address) (*common.Address, error) {
	// Assuming you have the private key of the owner who can call the `create` function
	privateKey, err := getPrivateKey() // You need to implement this function yourself
	if err != nil {
		return nil, err
	}

	nonce, err := f.Client.PendingNonceAt(context.Background(), owner)
	if err != nil {
		log.Fatalf("Failed to get the nonce: %v", err)
	}

	gasPrice, err := f.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	input, err := f.ContractABI.Pack(
		"create",
		name,
		symbol,
		decimals,
		owner,
		protocol,
		admin,
	)
	if err != nil {
		return nil, err
	}

	rawTx := types.NewTransaction(nonce, f.FactoryAddr, big.NewInt(0), auth.GasLimit, auth.GasPrice, input)
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, rawTx)
	if err != nil {
		return nil, err
	}

	err = f.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	txReceipt, err := f.Client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		return nil, err
	}

	// Assuming the new contract address is logged in the first log
	newContractAddress := txReceipt.Logs[0].Address
	return &newContractAddress, nil
}

