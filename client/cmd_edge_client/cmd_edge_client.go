package main

import (
	"fmt"
	"time"

	"github.com/garfeng/n2n_user_manager/client"
)

// TODO

var ()

func main() {
	controller := new(client.Controller)
	controller.SetConfigPath("./client/cmd_edge_client/config.toml")
	err := controller.ReadConfig()
	if err != nil {
		fmt.Println("fail to read config", err)
	}

	err = controller.LoginAndSetupN2NEdge("jiaru.yuan", "123456")
	if err != nil {
		fmt.Println(err)
		return
	}

	<-time.After(time.Second * 10)
}
