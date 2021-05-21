package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`

	// used for Dhcp server
	MacAddr string `json:"mac_addr"`
}

func getHandler(manager N2NManagerServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("invalid method"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request " + err.Error()))
			return
		}

		loginInfo := new(LoginInfo)
		err = json.Unmarshal(body, loginInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error " + err.Error()))
			return
		}

		params, err := manager.TryLoginAndGetParam(loginInfo.Username, loginInfo.Password, loginInfo.MacAddr)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden " + err.Error()))
			return
		}

		res, _ := json.MarshalIndent(params, "", "  ")
		w.Write(res)
	}
}

func SetupServer(port string, manager N2NManagerServer) error {
	http.HandleFunc("/api/n2n_params", getHandler(manager))
	return http.ListenAndServe(port, nil)
}
