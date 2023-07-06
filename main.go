package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"./erc20factory"
)

func main() {
	// Read ABI from file
	abi, err := ioutil.ReadFile("../abis/BitcoinAssetsErc20Factory.json")
	if err != nil {
		log.Fatalf("Failed to read contract ABI: %v", err)
	}

	// Instantiate the factory
	factory, err := erc20factory.NewERC20Factory("/path/to/ipc", string(abi), "0xYourFactoryContractAddress")
	if err != nil {
		log.Fatalf("Failed to create ERC20 factory: %v", err)
	}

	// Replace with your own values
	name := "MyToken"
	symbol := "MTK"
	protocol := "MyProtocol"
	decimals := uint8(18)
	owner := common.HexToAddress("0xYourOwnerAddress")
	admin := common.HexToAddress("0xYourAdminAddress")

	// Call CreateERC20
	newContractAddress, err := factory.CreateERC20(name, symbol, protocol, decimals, owner, admin)
	if err != nil {
		log.Fatalf("Failed to create ERC20 token: %v", err)
	}

	fmt.Printf("New contract address: %s\n", newContractAddress.Hex())
}

