package auth

import (
	"testing"
)

func Test_call(t *testing.T) {
	s := NewServer(":9689")
	defer s.Stop()

	_ = NewAuther([]string{"http://127.0.0.1:9689"}, "nodeId", "secret")

	nodeAuth := map[string]string{
		"nodeId": "is nodeid",
		"key":    "secret",
	}

	err := call(nodeAuth, node_auth_action)
	if err != nil {
		t.Error(err)
		return
	}

	nodeStatus := map[string]interface{}{
		"key":    "secret",
		"status": 1,
	}

	err = call(nodeStatus, node_status_action)
	if err != nil {
		t.Error(err)
		return
	}
}
