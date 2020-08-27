package mnc_raft

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func encodeBlock(block *types.Block) (str string, err error) {
	tmp := bytes.NewBuffer(nil)
	err = block.EncodeRLP(tmp)
	if err != nil {
		return "", err
	}
	str = hex.EncodeToString(tmp.Bytes())
	return str, err

}

func decodeBlock(hexRlpBlock string) (*types.Block, error) {
	buf, err := hex.DecodeString(hexRlpBlock)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stream := rlp.NewStream(bytes.NewReader(buf), 0)
	block := new(types.Block)
	err = block.DecodeRLP(stream)
	return block, err
}

func encodePubkey(key *ecdsa.PublicKey) string {
	buf := crypto.CompressPubkey(key)
	return hex.EncodeToString(buf)
}

func decodePubkey(keyStr string) (key *ecdsa.PublicKey, err error) {
	buf, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	return crypto.DecompressPubkey(buf)
}
