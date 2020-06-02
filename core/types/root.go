// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	// "bytes"
	"errors"
	// "fmt"
	"io"
	"math/big"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

// 根任务/根信用系统
type Root struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	// PostState         []byte `json:"root"`
	// Status            uint64 `json:"status"`
	// CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
	// Bloom Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs []*Log `json:"logs"              gencodec:"required"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	// GasUsed         uint64         `json:"gasUsed" gencodec:"required"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash `json:"blockHash,omitempty"`
	BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
	TransactionIndex uint        `json:"transactionIndex"`
}

type rootMarshaling struct {
	// PostState         hexutil.Bytes
	// Status            hexutil.Uint64
	// CumulativeGasUsed hexutil.Uint64
	// GasUsed           hexutil.Uint64
	BlockNumber      *hexutil.Big
	TransactionIndex hexutil.Uint
}

// receiptRLP is the consensus encoding of a receipt.
type rootRLP struct {
	// PostStateOrStatus []byte
	// CumulativeGasUsed uint64
	// Bloom             Bloom
	Logs []*Log
}

// storedReceiptRLP is the storage encoding of a receipt.
type storedRootRLP struct {
	// PostStateOrStatus []byte
	// CumulativeGasUsed uint64
	Logs []*LogForStorage
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewRoot(receipt *Receipt, topic common.Hash) *Root {
	r := &Root{TxHash: receipt.TxHash,
		BlockHash:        receipt.BlockHash,
		BlockNumber:      receipt.BlockNumber,
		TransactionIndex: receipt.TransactionIndex,
	}
	for _, l := range receipt.Logs {
		// log.Info("MMMMMMMMMMMMMMMMMMMNewCredit", "Topics", l.Topics[0], "CreditTopic", topic, "l.Topics[0] == CreditTopic", l.Topics[0] == topic)

		if l.Topics[0] == topic {
			r.Logs = append(r.Logs, l)
		}
	}
	// return c
	return r
}

// EncodeRLP implements rlp.Encoder, and flattens the consensus fields of a receipt
// into an RLP stream. If no post state is present, byzantium fork is assumed.
func (r *Root) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &rootRLP{r.Logs})
}

// DecodeRLP implements rlp.Decoder, and loads the consensus fields of a receipt
// from an RLP stream.
func (r *Root) DecodeRLP(s *rlp.Stream) error {
	var dec rootRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.Logs = dec.Logs
	return nil
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (r *Root) Size() common.StorageSize {
	size := common.StorageSize(unsafe.Sizeof(*r))

	size += common.StorageSize(len(r.Logs)) * common.StorageSize(unsafe.Sizeof(Log{}))
	for _, log := range r.Logs {
		size += common.StorageSize(len(log.Topics)*common.HashLength + len(log.Data))
	}
	return size
}

// RootForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type RootForStorage Root

// EncodeRLP implements rlp.Encoder, and flattens all content fields of a receipt
// into an RLP stream.
func (r *RootForStorage) EncodeRLP(w io.Writer) error {
	enc := &storedRootRLP{
		// PostStateOrStatus: (*Receipt)(r).statusEncoding(),
		// CumulativeGasUsed: r.CumulativeGasUsed,
		Logs: make([]*LogForStorage, len(r.Logs)),
	}
	for i, log := range r.Logs {
		enc.Logs[i] = (*LogForStorage)(log)
	}
	return rlp.Encode(w, enc)
}

// DecodeRLP implements rlp.Decoder, and loads both consensus and implementation
// fields of a receipt from an RLP stream.
func (r *RootForStorage) DecodeRLP(s *rlp.Stream) error {
	// Retrieve the entire receipt blob as we need to try multiple decoders
	blob, err := s.Raw()
	if err != nil {
		return err
	}
	// Try decoding from the newest format for future proofness, then the older one
	// for old nodes that just upgraded. V4 was an intermediate unreleased format so
	// we do need to decode it, but it's not common (try last).
	return decodeStoredRootRLP(r, blob)
}

func decodeStoredRootRLP(r *RootForStorage, blob []byte) error {
	var stored storedRootRLP
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	// if err := (*Receipt)(r).setStatus(stored.PostStateOrStatus); err != nil {
	// 	return err
	// }
	// r.CumulativeGasUsed = stored.CumulativeGasUsed
	r.Logs = make([]*Log, len(stored.Logs))
	for i, log := range stored.Logs {
		r.Logs[i] = (*Log)(log)
	}
	// r.Bloom = CreateBloom(Receipts{(*Receipt)(r)})

	return nil
}

// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Roots []*Root

// Len returns the number of receipts in this list.
func (r Roots) Len() int { return len(r) }

// GetRlp returns the RLP encoding of one receipt from the list.
func (r Roots) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(r[i])
	if err != nil {
		panic(err)
	}
	return bytes
}

// DeriveFields fills the receipts with their computed fields based on consensus
// data and contextual infos like containing block and transactions.
func (r Roots) DeriveFields(config *params.ChainConfig, hash common.Hash, number uint64, txs Transactions) error {
	signer := MakeSigner(config, new(big.Int).SetUint64(number))

	logIndex := uint(0)
	if len(txs) != len(r) {
		return errors.New("transaction and receipt count mismatch")
	}
	for i := 0; i < len(r); i++ {
		// The transaction hash can be retrieved from the transaction itself
		r[i].TxHash = txs[i].Hash()

		// block location fields
		r[i].BlockHash = hash
		r[i].BlockNumber = new(big.Int).SetUint64(number)
		r[i].TransactionIndex = uint(i)

		// The contract address can be derived from the transaction itself
		if txs[i].To() == nil {
			// Deriving the signer is expensive, only do if it's actually needed
			from, _ := Sender(signer, txs[i])
			r[i].ContractAddress = crypto.CreateAddress(from, txs[i].Nonce())
		}
		// The used gas can be calculated based on previous r
		// if i == 0 {
		// 	r[i].GasUsed = r[i].CumulativeGasUsed
		// } else {
		// 	r[i].GasUsed = r[i].CumulativeGasUsed - r[i-1].CumulativeGasUsed
		// }
		// The derived log fields can simply be set from the block and transaction
		for j := 0; j < len(r[i].Logs); j++ {
			r[i].Logs[j].BlockNumber = number
			r[i].Logs[j].BlockHash = hash
			r[i].Logs[j].TxHash = r[i].TxHash
			r[i].Logs[j].TxIndex = uint(i)
			r[i].Logs[j].Index = logIndex
			logIndex++
		}
	}
	return nil
}
