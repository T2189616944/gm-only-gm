package mnc_raft

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/rlp"
)

var (
	actionGetSignKey = []byte(`{"jsonrpc":"2.0","method":"solo_signKey","id":1}`)
)

type Client struct {
	addr string
}

func NewClient(addr string) (*Client, error) {
	c := &Client{
		addr: addr,
	}
	return c, nil
}

func (client *Client) GetSignKey() (*ecdsa.PublicKey, error) {
	body := bytes.NewReader(actionGetSignKey)
	resp, err := http.Post(client.addr, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rpcR, err := parseResult(resp.Body)
	if err != nil {
		fmt.Println("parse resp failed")
		fmt.Println(err)
		return nil, err
	}

	if rpcR.Error != nil {
		fmt.Println("sealerror")
		fmt.Println(rpcR.Error)
		return nil, fmt.Errorf(rpcR.Error.Message)
	}
	return decodePubkey(rpcR.Result)
}

func (client *Client) SendBlockToConsensus(block *types.Block) (result *types.Block, err error) {
	tmp := bytes.NewBuffer(nil)
	tmp.WriteString(`{"jsonrpc":"2.0","method":"solo_blockToConsensus","id":1,"params":["`)
	str, err := encodeBlock(block)
	if err != nil {
		return nil, err
	}
	tmp.WriteString(str)
	tmp.WriteString(`"]}`)

	resp, err := http.Post(client.addr, "application/json", tmp)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rpcR, err := parseResult(resp.Body)
	if err != nil {
		fmt.Println("parse resp failed: " + err.Error())
		return nil, err
	}

	if rpcR.Error != nil {
		fmt.Println("seal failed: ", rpcR.Error)
		return nil, fmt.Errorf(rpcR.Error.Message)
	}

	return decodeBlock(rpcR.Result)

}

// {"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"debug:get block 265"}}
// {"jsonrpc":"2.0","id":1,"result":"03a9cbcf7f67ebbecda1d0c27cb4b0527e82a3422511f21ed67f3b1b3c70db6fef"}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResult struct {
	Jsonrpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Error   *rpcError `json:"error,omitempty"`
	Result  string    `json:"result,omitempty"`
}

func parseResult(r io.Reader) (*rpcResult, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	result := new(rpcResult)
	err = json.Unmarshal(buf, result)
	return result, err
}
