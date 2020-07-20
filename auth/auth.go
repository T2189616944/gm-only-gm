package auth

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

var (
	gAuther *Auther

	deployAccountAddr = common.HexToAddress("0X8E514EAEB5BADB3B5C929C41BF78F06D7271C8AF")

	ERROR_UNAUTHORIZED_ACCOUNT = errors.New("Unauthorized deployment contract")
	ERROR_NODE_AUTH_FAILED     = errors.New("node auth failed")
)

func TxAuth(contractAddr *common.Address, accountAddr common.Address) error {
	return gAuther.TxAuth(contractAddr, accountAddr)

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
		return errors.New("auth.server address is empty ")
	}
	if auther.nodeId == "" {
		return errors.New("NodeId is empty")
	}
	if auther.secret == "" {
		return errors.New("auth.code is empty")
	}

	if err := SendNodeAuth(auther.nodeId, auther.secret); err != nil {
		return err
	}

	if err := SendNodeUpdateStatus(auther.secret, NODE_STATUS_ONLINE); err != nil {
		return err
	}

	return nil
}

func (auther *Auther) TxAuth(contractAddr *common.Address, accountAddr common.Address) error {

	//  合约安装
	// 判断是否可安装合约用户
	// 这个可安装合约用户写在代码里？还是通过在线方式获得？
	// 目前简单写死吧
	if contractAddr == nil {
		if accountAddr == deployAccountAddr {
			return nil
		} else {
			return ERROR_UNAUTHORIZED_ACCOUNT
		}
	}

	return SendTxAuth(contractAddr, accountAddr)

}

func (auther *Auther) Close() {
	SendNodeUpdateStatus(auther.secret, NODE_STATUS_OFFLINE)
}
