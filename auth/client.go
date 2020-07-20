package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	NODE_STATUS_ONLINE  = 3
	NODE_STATUS_OFFLINE = 4

	node_auth_action   = "/rod-api/node/getIsStart"
	node_status_action = "/rod-api/node/updateStatus"
)

//
func SendNodeUpdateStatus(secret string, status int) error {

	params := map[string]interface{}{
		"status": status,
		"key":    secret,
	}

	return call(params, node_status_action)
}

func SendNodeAuth(nodeId, secret string) error {
	params := map[string]string{
		"nodeId": nodeId,
		"key":    secret,
	}

	return call(params, node_auth_action)

}

func SendTxAuth(contractAddr *common.Address, accountAddr common.Address) error {
	log.Info("check tx auth")
	return nil
}

type Response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
	Result    string `json:"result"`
}

func call(params interface{}, action string) error {
	buf, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// fmt.Println(action)
	// fmt.Println("send ", string(buf))

	body := bytes.NewReader(buf)

	for _, bc := range gAuther.addrs {
		resp, err := http.Post(bc+action, "application/json", body)
		if err != nil {
			fmt.Println("connect auth serve failed:", err)
			continue
		}
		defer resp.Body.Close()

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("resad resp body failed", err)
			return ERROR_NODE_AUTH_FAILED
		}
		tmp := new(Response)
		if err = json.Unmarshal(buf, tmp); err != nil {
			fmt.Println("unmarshal failed", err.Error())
			fmt.Println(string(buf))
			return ERROR_NODE_AUTH_FAILED
		}
		if tmp.Code == 0 {
			return nil
		}

		// fmt.Println("response code is not 0")
		// fmt.Println(string(buf))

		return ERROR_NODE_AUTH_FAILED
	}

	return ERROR_NODE_AUTH_FAILED
}
