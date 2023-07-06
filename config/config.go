package config

type Network struct {
	RpcUrl                 string
	FactoryContractAddress string
}

var BEVMMainnet = Network{
	RpcUrl:                 "https://mainnet.chainx.org/rpc",
	FactoryContractAddress: "0x124e3E8D56db6ADA37aF2b7662F275D49BA850e6",
}

var BEVMTestnet = Network{
	RpcUrl:                 "https://testnet3.chainx.org/rpc",
	FactoryContractAddress: "0xeB789d5f6f66104AE9876175A5B9A03bDa0545A8",
}
