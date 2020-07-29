package auth

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_call(t *testing.T) {
	// s := NewServer(":9689")
	// defer s.Stop()

	_ = NewAuther([]string{"http://192.168.2.208:9094"}, "nodeId", "secret")

	// uid 1

	// 20200716193855
	// ecbd8ac31b5a0e7a0ad0c6db6096ab20f27eeb0053bb8ed45bfb3aaae1acf77c
	// 55b78881ac922a05f26a2725c247c374223ae7b89daf224cfa9c4d8d713e6ebb
	aAddr := common.HexToAddress("0X50D0E7CD53A476AF05245E63ED174985BC11FD57")

	// 16
	cAddr := common.HexToAddress("0X93C72E8BBC3679DED22570D9139E6BDE7335F669")

	// 4
	// 0XCBA57F50736BA27AD2E6991E3F12AB624BFBB9FF

	// tag 19

	// c_tag 40

	fmt.Println(TxAuth(&cAddr, aAddr))
	// return

	nodeAuth := map[string]string{
		"nodeId": "6d0e800d65bd5b108b56b635a86b4fafb053f42f37fd19c2b86579511c051cfe198638a18876bd207dab738fd1ffa6102a41c1ed09235a19a6f8c0b4fe81559e",
		"key":    "320389200717164634",
	}

	err := call(nodeAuth, node_auth_action, cType_json)
	if err != nil {
		t.Error(err)
		// return
	}

	nodeStatus := map[string]interface{}{
		"key":    "320389200717164634",
		"status": 3,
	}

	err = call(nodeStatus, node_status_action, cType_json)
	if err != nil {
		t.Error(err)
		return
	}
}
