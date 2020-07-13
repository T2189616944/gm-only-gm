package auth

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	gAuther *Auther
)

func TxAuth(comtractAddr *common.Address, accountAddr common.Address) error {
	return gAuther.TxAuth(comtractAddr, accountAddr)

}

func NodeAuth() error {
	return gAuther.NodeAuth()

}

func Close() {
	if gAuther != nil {
		gAuther.Close()
		gAuther = nil
	}
}

type Auther struct {
	addrs  []string
	secret string
	nodeId string
}

func NewAuther(bc []string, nodeId, secret string) *Auther {

	gAuther = &Auther{
		addrs:  bc,
		secret: secret,
		nodeId: nodeId,
	}

	return gAuther
}

func (auther *Auther) NodeAuth() error {
	if len(auther.addrs) == 0 {
		return errors.New("License server address is empty ")
	}
	if err := SendNodeAuth(auther.nodeId, auther.secret); err != nil {
		return err
	}

	if err := SendNodeUpdateStatus(auther.secret, NODE_STATUS_ONLINE); err != nil {
		return err
	}

	return nil
}

func (auther *Auther) TxAuth(comtractAddr *common.Address, accountAddr common.Address) error {

	return SendTxAuth(comtractAddr, accountAddr)

}

func (auther *Auther) Close() {
	SendNodeUpdateStatus(auther.secret, NODE_STATUS_OFFLINE)
}
