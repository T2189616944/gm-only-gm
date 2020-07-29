// +build !noauth

package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

const (
	cType_json      = "application/json"
	cType_urlencode = "application/x-www-form-urlencoded"

	NODE_STATUS_ONLINE  = "3"
	NODE_STATUS_OFFLINE = "4"

	node_auth_action   = "newchain/rod-api/node/getIsStart"
	node_status_action = "newchain/rod-api/node/updateStatus"
	node_keepalive     = "newchain/nodeWebsocket"
	tx_auht_action     = "newchain/rod-api/contract/checkAccountIsEmpowered"
)

var (
	gAuther *Auther

	deployAccountAddr = common.HexToAddress("0X8E514EAEB5BADB3B5C929C41BF78F06D7271C8AF")

	ERROR_UNAUTHORIZED_ACCOUNT = errors.New("Unauthorized deployment contract")
	ERROR_NODE_AUTH_FAILED     = errors.New("node auth failed")
	ERROR_TX_AUTH_FAILED       = errors.New("tx auth failed")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		CheckOrigin: func(c *http.Request) bool {
			return true
		},
	}
)

func TxAuth(contractAddr *common.Address, accountAddr common.Address) error {
	return gAuther.TxAuth(contractAddr, accountAddr)

}

func NodeAuth() error {
	return gAuther.NodeAuth()
}

func Keepalive() error {
	return gAuther.Keepalive()

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
	for i, b := range bc {
		if !strings.HasSuffix(b, "/") {
			bc[i] = b + "/"
		}
	}

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

	params := map[string]string{
		"nodeId": auther.nodeId,
		"key":    auther.secret,
	}
	if err := call(params, node_auth_action, cType_json); err != nil {
		return err
	}

	params = map[string]string{
		"status": NODE_STATUS_ONLINE,
		"key":    auther.secret,
	}
	if err := call(params, node_status_action, cType_json); err != nil {
		fmt.Println(err)
	}

	return nil
}

func (auther *Auther) Keepalive() error {
	conn, err := SelectTxAuthServer(auther.addrs, auther.secret)
	if err != nil {
		return err
	}

	_, _, err = conn.ReadMessage()
	return err
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

	params := make(url.Values)
	params.Add("contractAddress", contractAddr.String())
	params.Add("accountAddress", accountAddr.String())

	if err := call(params, tx_auht_action, cType_urlencode); err != nil {
		return ERROR_TX_AUTH_FAILED
	}
	return nil

}

func (auther *Auther) Close() {
	params := map[string]string{
		"status": NODE_STATUS_OFFLINE,
		"key":    auther.secret,
	}
	if err := call(params, node_status_action, cType_json); err != nil {
		fmt.Println(err)
	}
}

type Response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Timestamp int64  `json:"timestamp"`
	Result    string `json:"result"`
}

func call(params interface{}, action string, ctype string) error {
	var body io.Reader
	if ctype == cType_json {
		buf, err := json.Marshal(params)
		if err != nil {
			return err
		}

		body = bytes.NewReader(buf)

	} else {
		body = strings.NewReader(params.(url.Values).Encode())
	}
	for _, bc := range gAuther.addrs {

		resp, err := http.Post(bc+action, ctype, body)
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
		if tmp.Code == 200 {
			return nil
		}

		fmt.Println("response code is not 0")
		fmt.Println(string(buf))

		return ERROR_NODE_AUTH_FAILED
	}

	return ERROR_NODE_AUTH_FAILED
}

// 简单来
func SelectTxAuthServer(auths []string, secret string) (conn *websocket.Conn, err error) {
	// ......
	// ......
	dialer := websocket.Dialer{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
	}
	for _, url := range auths {
		if strings.HasPrefix(url, "https") {
			url = strings.Replace(url, "https", "wss", -1)
		} else {
			url = strings.Replace(url, "http", "ws", -1)
		}

		url = fmt.Sprintf("%s%s/%s", url, node_keepalive, secret)
		conn, _, err = dialer.Dial(url, nil)

		if err == nil {
			return conn, nil
		} else {
			fmt.Println(err)
		}
	}
	return nil, fmt.Errorf("all backend is failed:%s", err.Error())
}
