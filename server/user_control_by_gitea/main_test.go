package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/garfeng/n2n_user_manager/server"
)

func Test_Login(t *testing.T) {
	info := &server.LoginInfo{
		Username: "garfeng",
		Password: "000000",
		MacAddr:  "helloworld",
	}
	b, _ := json.MarshalIndent(info, "", "  ")
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api/n2n_params", bytes.NewBuffer(b))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}
