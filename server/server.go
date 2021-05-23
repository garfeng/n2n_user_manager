package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/garfeng/n2n_user_manager/common/user"
)

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

		loginInfo := new(user.LoginInfo)
		err = json.Unmarshal(body, loginInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error " + err.Error()))
			return
		}
		log.Println("user login:", loginInfo.Username)
		params, err := manager.TryLoginAndGetParam(loginInfo.Username, loginInfo.Password)
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
