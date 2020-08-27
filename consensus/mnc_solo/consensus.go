package mnc_solo

import (
	// "bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	// "golang.org/x/crypto/sha3"
	"github.com/tjfoc/gmsm/sm3"
)

var (
	lock = new(sync.Mutex)
)

type MNCSolo struct {
	signPubkey         *ecdsa.PublicKey
	compressSignPubkey []byte
	config             *params.SoloConfig
	db                 ethdb.Database
	chain              *core.BlockChain

	sealTime                  time.Duration
	updateConsensusStatusLock sync.RWMutex

	client *Client
}

func New(config *params.SoloConfig, db ethdb.Database) (*MNCSolo, error) {

	solo := &MNCSolo{
		// signPriv: toSm2Key(config.Priv),
		config: config,
		db:     db,
	}
	solo.sealTime = time.Second * time.Duration(config.SealTime)

	if !config.IsSoloNode {
		c, err := NewClient(config.SoloNodeAddr)
		if err != nil {
			return nil, err
		}
		solo.client = c

		solo.signPubkey, err = solo.client.GetSignKey()
		if err != nil {
			return solo, err
		}
	} else {
		solo.signPubkey = &solo.config.Priv.PublicKey
	}

	solo.compressSignPubkey = crypto.CompressPubkey(solo.signPubkey)

	return solo, nil
}

func (solo *MNCSolo) SetBlockchain(chain *core.BlockChain) {
	solo.chain = chain

}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
// coinbase , 区块作者
// signatures
func (solo *MNCSolo) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// 使用本机签名的头部就是有效的
func (solo *MNCSolo) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {

	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if len(header.Extra) < crypto.SignatureLength {
		return errors.New("Extra to small")
	}

	// check sign
	hash := solo.SealHash(header)
	sig := header.Extra[len(header.Extra)-crypto.SignatureLength:]

	if len(sig) != crypto.SignatureLength {
		fmt.Println(sig)
		return errors.New("sign data len error ")
	}
	ok := solo.verifySignature(hash, sig)
	if !ok {

		return ErrVerifyHeaderFailed
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
//
func (solo *MNCSolo) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, _ := range headers {
			err := solo.verifyHeaderWorker(chain, headers, seals, i)
			// err := solo.VerifyHeader(chain, header, true)
			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (solo *MNCSolo) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	header := headers[index]
	if chain.GetHeader(header.Hash(), header.Number.Uint64()) != nil {
		return nil // known block
	}

	if len(header.Extra) < crypto.SignatureLength {
		return errors.New("Extra to small")
	}
	// check sign
	hash := solo.SealHash(header)
	sig := header.Extra[len(header.Extra)-crypto.SignatureLength:]

	if len(sig) != crypto.SignatureLength {
		fmt.Println(sig)
		return errors.New("sign data len error ")
	}
	ok := solo.verifySignature(hash, sig)
	if !ok {
		fmt.Println("???")
		fmt.Println(sig)
		fmt.Println(hash)
		fmt.Println(header.Number.Uint64())
		return ErrVerifyHeaderFailed
	}
	return nil

}

func (solo *MNCSolo) verifySignature(hash common.Hash, sig []byte) bool {
	if solo.compressSignPubkey == nil {
		var err error
		solo.signPubkey, err = solo.client.GetSignKey()
		if err != nil {
			fmt.Println("can not connect solo node, get sign key")
			fmt.Println(err)
			return false
		}
		solo.compressSignPubkey = crypto.CompressPubkey(solo.signPubkey)
	}

	return crypto.VerifySignature(solo.compressSignPubkey, hash.Bytes(), sig[:64])
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (solo *MNCSolo) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (solo *MNCSolo) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return solo.VerifyHeader(chain, header, true)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (solo *MNCSolo) Prepare(chain consensus.ChainReader, header *types.Header) error {

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	header.Difficulty = solo.CalcDifficulty(chain, header.Time, parent)
	if len(header.Extra) == 0 {
		header.Extra = make([]byte, crypto.SignatureLength)
	} else {
		buf := make([]byte, len(header.Extra)+crypto.SignatureLength)
		copy(buf, header.Extra)
		header.Extra = buf
		// header.Extra = append(header.Extra, make([]byte, crypto.SignatureLength)...)
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// but does not assemble the block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (solo *MNCSolo) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
}

// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
// rewards) and assembles the final block.
//
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (solo *MNCSolo) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return types.NewBlock(header, txs, uncles, receipts), nil
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (solo *MNCSolo) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (err error) {

	// sign seal hash
	header := block.Header()

	number := header.Number.Uint64()
	if number == 0 {
		return errors.New("block number is 0")
	}

	if !solo.config.IsSoloNode {
		// 只有solo 可以发送
		return ErrNotSoloNode
		// 非solo 节点提交到solo 节点
		block, err = solo.client.SendBlockToConsensus(block)
		if err != nil {
			return err
		}

		err = solo.config.Mux.Post(core.NewMinedBlockEvent{Block: block})
		if err != nil {
			fmt.Println(err)
		}
		results <- block

	} else {
		header, err = solo.seal(header, chain)
		if err != nil {
			return err
		}
		results <- block.WithSeal(header)
	}

	return nil
}

func (solo *MNCSolo) seal(header *types.Header, chain consensus.ChainReader) (*types.Header, error) {
	solo.updateConsensusStatusLock.Lock()
	defer solo.updateConsensusStatusLock.Unlock()

	time.Sleep(solo.sealTime)

	cH := chain.CurrentHeader()
	if cH == nil {
		return nil, fmt.Errorf("get current header failed")
	}

	want := cH.Number.Uint64() + 1
	num := header.Number.Uint64()
	if want != num {
		return nil, fmt.Errorf("want block number %d, bug got block int %d", want, num)
	}

	// *****
	hash := solo.SealHash(header)
	var isok bool
	for i := 0; i < 5; i++ {

		sig, err := crypto.SignWithPub(hash.Bytes(), solo.config.Priv)
		if err != nil {
			return nil, err
		}

		isok = solo.verifySignature(hash, sig)
		if isok {
			copy(header.Extra[len(header.Extra)-crypto.SignatureLength:], sig)
			break
		} else {

			continue
		}
	}
	if !isok {
		return nil, errors.New("sign failed")
	}
	return header, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (solo *MNCSolo) SealHash(header *types.Header) (hash common.Hash) {

	data, err := rlp.EncodeToBytes([]interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-crypto.SignatureLength],
	})
	if err != nil {
		panic(err)
	}

	buf := sm3.Sm3Sum(data)
	return common.BytesToHash(buf)
	// hasher.Sum(hash[:0])
	// return hash

}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (solo *MNCSolo) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {

	return big.NewInt(10)
}

// APIs returns the RPC APIs this consensus engine provides.
func (solo *MNCSolo) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "solo",
		Version:   "1.0",
		Service:   &API{chain: chain, solo: solo},
		Public:    true,
	}}
	return nil
}

// Close terminates any background threads maintained by the consensus engine.
func (solo *MNCSolo) Close() error {
	printdebugLog("close")
	return nil

}

func printdebugLog(msg string) {
	fmt.Println("xxxxx")
	fmt.Println("xxxxx")
	fmt.Println(msg)
	fmt.Println("xxxxx")
	fmt.Println("xxxxx")

}
