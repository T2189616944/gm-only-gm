// +build noauth

package auth

import (
	"github.com/ethereum/go-ethereum/common"
)

type Auther struct {
}

func NewAuther(bc []string, nodeId, secret string) *Auther {
	return nil
}

func TxAuth(contractAddr *common.Address, accountAddr common.Address) error {

	return nil
}

func NodeAuth() error {
	return nil
}

func Keepalive() error {
	return nil
}

func Close() {
}
