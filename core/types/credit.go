package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)
var CreditTopic = common.HexToHash("1d38399075534ff764a8ae16a9ba0294e39ec2125af3846882eace07fe737d44")

type Credit struct {

	Logs              []*Log `json:"logs"              gencodec:"required"`
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}



func NewCredit(receipt *Receipt)*Credit{
	c := &Credit{TxHash: receipt.TxHash,
		BlockHash: receipt.BlockHash,
		BlockNumber:receipt.BlockNumber,
		TransactionIndex:receipt.TransactionIndex,
	}
	for _, l := range receipt.Logs{
		log.Info("MMMMMMMMMMMMMMMMMMMNewCredit", "Topics",l.Topics[0], "CreditTopic", CreditTopic, "l.Topics[0] == CreditTopic", l.Topics[0] == CreditTopic)

		if l.Topics[0] == CreditTopic{
			c.Logs = append(c.Logs, l)
		}
	}
	return c
}


type Credits []*Credit

func (c Credits) Len() int { return len(c) }

// GetRlp returns the RLP encoding of one receipt from the list.
func (c Credits) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(c[i])
	if err != nil {
		panic(err)
	}
	return bytes
}
