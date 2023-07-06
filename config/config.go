package config

type Network struct {
	RpcUrl                 string
	FactoryContractAddress string
}

var Mainnet = Network{
	RpcUrl:                 "https://mainnet.chainx.org/rpc",
	FactoryContractAddress: "0xeB789d5f6f66104AE9876175A5B9A03bDa0545A8",
}

var Ropsten = Network{
	RpcUrl:                 "https://testnet3.chainx.org/rpc",
	FactoryContractAddress: "0xeB789d5f6f66104AE9876175A5B9A03bDa0545A8",
}
