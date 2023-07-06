# bevm-erc20-factory

**bevm-erc20-factory** is a Golang-based library that allows the deployment of ERC20 tokens by interacting with an Ethereum smart contract factory.

## Installation
First, you need to install this library in your project. If you're using Go Modules, you can run:

```
go get github.com/chainx-org/bevm-erc20-factory
```

## Configuration
Create a .env file in the root directory of your project and add the following line to set your private key:
```
PRIVATE_KEY=YourPrivateKey
```
You also need to configure the network settings. You can edit the config.go file to suit your needs. For example:

```
package config

type Network struct {
	RpcUrl                 string
	FactoryContractAddress string
}

var BEVMTestnet = Network{
	RpcUrl:                 "https://yourtestnet.infura.io/v3/YOUR_INFURA_PROJECT_ID",
	FactoryContractAddress: "0xYourTestnetFactoryContractAddress",
}

```

## Usage

To use this library, first import it:

import (
	"bevm-erc20-factory/factory"
	"github.com/ethereum/go-ethereum/common"
)

Then, you can create a new instance of ERC20Factory:

```
abi := "<your abi here>"
factoryAddr := common.HexToAddress("<factory contract address here>")
factory, err := factory.NewERC20Factory("<your rpc url here>", abi, factoryAddr)
if err != nil {
	// handle error
}

```
Now, you can use this instance to create new ERC20 tokens:

```
name := "ABCD"
symbol := "ABCD"
protocol := "brc-20"
decimals := uint8(18)
owner := common.HexToAddress("<owner address here>")
admin := common.HexToAddress("<admin address here>")
newContractAddress, err := factory.CreateERC20(name, symbol, protocol, decimals, owner, admin)
if err != nil {
	// handle error
}

```
In the above code, `name`, `symbol`, `protocol` are the name, symbol, and protocol of your token. `decimals` is the number of decimal places for your token. `owner` and `admin` are the owner and admin addresses for your token contract.

You can refer to the main.go file in the project repository for a more detailed usage example.

## Note
Make sure that you've correctly set up your .env file, and that your contract ABI and address are correct. If you encounter any issues, don't hesitate to ask.
