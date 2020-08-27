package params

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

type SoloConfig struct {
	IsSoloNode bool
	Priv       *ecdsa.PrivateKey `json:"-"`
	Mux        *event.TypeMux    `json:"-"`
	// Chain        *core.BlockChain  `json:"-"`
	SoloNodeAddr string
	SealTime     int `json:"-"`
}

func (c *SoloConfig) String() string {
	return "mnc_solo"
}

type RaftConfig struct {
	IsSoloNode bool
	Priv       *ecdsa.PrivateKey `json:"-"`
	Mux        *event.TypeMux    `json:"-"`
	// Chain        *core.BlockChain  `json:"-"`
	SoloNodeAddr string
	SealTime     int `json:"-"`
}

func (c *RaftConfig) String() string {
	return "mnc_raft"
}

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	SoloChainConfig = &ChainConfig{
		ChainID:             big.NewInt(110),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		Solo:                new(SoloConfig),
	}

	RaftChainConfig = &ChainConfig{
		ChainID:             big.NewInt(110),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      false,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		Raft:                new(RaftConfig),
	}
)

func NewRaftConfig() error {
	RaftChainConfig.Raft = &RaftConfig{}
	return nil

}

func NewSoloConfig(isMainNode bool, key string, addr string, sealTime int) error {
	SoloChainConfig.Solo = &SoloConfig{
		IsSoloNode: isMainNode,
		// Priv:         priv,
		SealTime:     sealTime,
		SoloNodeAddr: addr,
	}
	if isMainNode {
		priv, err := crypto.HexToECDSA(key)
		if err != nil {
			return fmt.Errorf("parse solo.key failed:%s", err.Error())
		}
		SoloChainConfig.Solo.Priv = priv
		SoloChainConfig.Solo.SoloNodeAddr = "localhost:8545"
	}

	return nil
}
