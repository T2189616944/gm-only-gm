package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc(node_auth_action, ActionNodeAuth)
	mux.HandleFunc(node_status_action, ActionNodeUpdateStatus)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
	}

	go server.ListenAndServe()
	time.Sleep(time.Second * 3)
	s := &Server{
		server: server,
	}

	return s
}

func (s *Server) Stop() {
	s.server.Shutdown(context.TODO())
}

func ActionNodeAuth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params, err := parseBody(r.Body)
	if err != nil {
		w.Write(buildResp(1, err.Error()))
		return

	}
	_, has := params["nodeId"]
	if !has {
		w.Write(buildResp(1, "miss params nodeId"))
		return
	}
	_, has = params["key"]
	if !has {
		w.Write(buildResp(1, "miss params key"))
		return

	}
	w.Write(buildResp(0, "ok"))

}

func ActionNodeUpdateStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params, err := parseBody(r.Body)
	if err != nil {
		w.Write(buildResp(1, err.Error()))
		return
	}
	status, has := params["status"]
	if !has {
		w.Write(buildResp(1, "miss params nodeId"))
		return

	}
	_, has = params["key"]
	if !has {
		w.Write(buildResp(1, "miss params key"))
		return

	}
	fmt.Println(status)
	w.Write(buildResp(0, "ok"))
}

func ActionTxAuth(w http.ResponseWriter, r *http.Request) {

}

func buildResp(code int, message string) []byte {
	resp := &Response{
		Code:      code,
		Success:   true,
		Message:   message,
		Timestamp: time.Now().Unix(),
		Result:    "is reuslt",
	}
	data, _ := json.Marshal(resp)
	return data
}

func parseBody(body io.Reader) (params map[string]interface{}, err error) {
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	params = make(map[string]interface{})
	err = json.Unmarshal(buf, &params)
	return params, err
}
