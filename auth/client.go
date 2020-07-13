package auth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	NODE_STATUS_ONLINE  = "online"
	NODE_STATUS_OFFLINE = "offline"
)

func SendNodeUpdateStatus(secret, status string) error {
	log.Info("udpate node status to " + status)
	return nil

}
func SendNodeAuth(nodeId, secret string) error {
	log.Info("check node auth")
	return nil

}

func SendTxAuth(comtractAddr *common.Address, accountAddr common.Address) error {
	log.Info("check tx auth")
	return nil
}
