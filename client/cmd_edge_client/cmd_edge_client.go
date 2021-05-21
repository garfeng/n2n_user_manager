package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garfeng/n2n_user_manager/client"
)

// TODO

var ()

func main() {
	if len(os.Args) < 3 {
		fmt.Println("cmd_edge_client <username> <password>")
		return
	}

	controller := new(client.Controller)
	controller.SetConfigPath("config.toml")
	err := controller.ReadConfig()
	if err != nil {
		fmt.Println("fail to read config", err)
	}

	err = controller.LoginAndSetupN2NEdge(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		<-time.After(time.Second * 10)
		return
	}

	defer controller.Disconnect()
	pause()
}

func pause() {
	fmt.Println("press Ctrl+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	<-c
}
