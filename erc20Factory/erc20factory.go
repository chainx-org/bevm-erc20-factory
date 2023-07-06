package erc20Factory

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type ERC20Factory struct {
	ContractABI abi.ABI
	FactoryAddr common.Address
	Client      *ethclient.Client
}

func GetPrivateKeyFromEnv() (*ecdsa.PrivateKey, common.Address, error) {
	if _, err := os.Stat(".env"); err == nil {
		// .env exists, load it
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, common.Address{}, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to convert private key: %v", err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	addressHex := crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, addressHex, nil
}

func (f *ERC20Factory) EstimateGas(input []byte, from, to common.Address, gasPrice *big.Int) (uint64, error) {
	callMsg := ethereum.CallMsg{
		From:     from,
		To:       &to,
		GasPrice: gasPrice,
		Data:     input,
	}

	estimatedGas, err := f.Client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return 0, fmt.Errorf("Failed to estimate gas: %v", err)
	}

	return estimatedGas, nil
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
	privateKey, addressHex, err := GetPrivateKeyFromEnv()

	if err != nil {
		return nil, err
	}

	nonce, err := f.Client.PendingNonceAt(context.Background(), addressHex)

	if err != nil {
		log.Fatalf("Failed to get the nonce: %v", err)
	}

	gasPrice, err := f.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
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
	estimatedGas, err := f.EstimateGas(input, addressHex, f.FactoryAddr, gasPrice)
	if err != nil {
		log.Printf("Failed to estimate gas: %v, fallback to default gas limit", err)
		estimatedGas = uint64(2200000) // Fallback to default gas limit if estimate failed
	}

	auth.GasLimit = estimatedGas

	tx := types.NewTransaction(nonce, f.FactoryAddr, big.NewInt(0), auth.GasLimit, auth.GasPrice, input)

	// Get chain ID from the Ethereum client
	chainID, err := f.Client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	fmt.Printf("User Address: %s, Nonce: %d, Chain ID: %s \n", addressHex.Hex(), nonce, chainID.String())

	signer := types.NewEIP155Signer(chainID)
	sighash := signer.Hash(tx)
	signature, err := crypto.Sign(sighash[:], privateKey)
	if err != nil {
		return nil, err
	}

	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return nil, err
	}

	err = f.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Signed Tx Hash: %s\n", signedTx.Hash())

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Set your desired timeout duration
	defer cancel()

	var txReceipt *types.Receipt
	for txReceipt == nil {
		txReceipt, err = f.Client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				// The wait timed out
				return nil, fmt.Errorf("transaction receipt was not found within the specified time")
			} else if err.Error() != "not found" {
				return nil, err
			}
		}
		time.Sleep(time.Second * 5)
	}

	// Assuming the new contract address is logged in the first log
	logData := txReceipt.Logs[2].Data
	contractAddressHex := common.Bytes2Hex(logData)               // convert the hex data to bytes
	newContractAddress := common.HexToAddress(contractAddressHex) // convert the bytes to an Address

	return &newContractAddress, nil
}
