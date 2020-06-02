package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

var TaskTopic = common.HexToHash("4dbd6453d307172d4c9ff2f95539271a96d6329aa99c736b3e4e6775bcdba202")

type Task struct {
	//PostState         []byte `json:"root"`
	//Status            uint64 `json:"status"`
	//CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
	//Bloom             Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs []*Log `json:"logs"              gencodec:"required"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	//GasUsed         uint64         `json:"gasUsed" gencodec:"required"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

func NewTask(receipt *Receipt) *Task {
	t := &Task{TxHash: receipt.TxHash,
		BlockHash:        receipt.BlockHash,
		BlockNumber:      receipt.BlockNumber,
		TransactionIndex: receipt.TransactionIndex,
	}
	for _, l := range receipt.Logs {
		log.Info("MMMMMMMMMMMMMMMMMMMNewTask", "Topics", l.Topics[0], "TaskTopic", TaskTopic, "l.Topics[0] == TaskTopic", l.Topics[0] == TaskTopic)

		if l.Topics[0] == TaskTopic {
			t.Logs = append(t.Logs, l)
		}
	}
	return t
}

type Tasks []*Task

func (t Tasks) Len() int { return len(t) }

// GetRlp returns the RLP encoding of one receipt from the list.
func (t Tasks) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(t[i])
	if err != nil {
		panic(err)
	}
	return bytes
}
