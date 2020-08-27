package mnc_raft

import (
	"fmt"

	"github.com/ethereum/go-ethereum/consensus"

	"github.com/ethereum/go-ethereum/crypto"
)

type API struct {
	chain consensus.ChainReader
	solo  *MNCSolo
}

func (api *API) SignKey() (string, error) {
	// if api.solo.config.IsSoloNode {
	// 	return encodePubkey(&api.solo.config.Priv.PublicKey)
	// }
	if api.solo.compressSignPubkey == nil {
		var err error
		api.solo.signPubkey, err = api.solo.client.GetSignKey()
		if err != nil {
			fmt.Println("can not connect solo node, get sign key")
			fmt.Println("err")
			return "", err
		}
		api.solo.compressSignPubkey = crypto.CompressPubkey(api.solo.signPubkey)
	}
	return encodePubkey(api.solo.signPubkey), nil
}

/*
0
f9029ef90299a00000000000000000000000000000000000000000000000000000000000000000a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001808347b760808080820400b8620000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000880000000000000000c0c0


265
f902b9f902b4a08351a4d28ccc22cac0dd6b5331d3764827085b5510b559bcfd3e0b57b0956a43a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479444cf77f78325583d9ab9378f1e7d316f8bb0b3b5a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a820109835ce138808080845f3641f8b879d683010912846765746886676f312e3135856c696e7578c8127ef7bf957cd26f831d0c447d707f6a6c98b23bc5511e000e3a56f704e579428081cc9ed797f4feb379722b5e7ffbf57ffb00a4d500c96fa5673aeff0a6050003a9cbcf7f67ebbecda1d0c27cb4b0527e82a3422511f21ed67f3b1b3c70db6fefa00000000000000000000000000000000000000000000000000000000000000000880000000000000000c0c0


curl -X POST -H "Content-type: application/json" --data '{"jsonrpc":"2.0","method":"solo_blockToConsensus","id":1,"params":["f902b9f902b4a08351a4d28ccc22cac0dd6b5331d3764827085b5510b559bcfd3e0b57b0956a43a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479444cf77f78325583d9ab9378f1e7d316f8bb0b3b5a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a820109835ce138808080845f3641f8b879d683010912846765746886676f312e3135856c696e7578c8127ef7bf957cd26f831d0c447d707f6a6c98b23bc5511e000e3a56f704e579428081cc9ed797f4feb379722b5e7ffbf57ffb00a4d500c96fa5673aeff0a6050003a9cbcf7f67ebbecda1d0c27cb4b0527e82a3422511f21ed67f3b1b3c70db6fefa00000000000000000000000000000000000000000000000000000000000000000880000000000000000c0c0"]}' 127.0.0.1:8545

*/

func (api *API) BlockToConsensus(hexRlpBlock string) (out string, err error) {

	if !api.solo.config.IsSoloNode {
		return "", fmt.Errorf("is not consensue note. need sent to: %s", api.solo.config.SoloNodeAddr)
	}

	block, err := decodeBlock(hexRlpBlock)
	if err != nil {
		return "", err
	}

	header, err := api.solo.seal(block.Header(), api.chain)
	if err != nil {
		return "", err
	}
	block = block.WithSeal(header)

	return encodeBlock(block)
}
