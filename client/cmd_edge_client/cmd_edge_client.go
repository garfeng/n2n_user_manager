package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/garfeng/n2n_user_manager/client"
)

var (
	username = flag.String("u", "", "<-u username>")
	password = flag.String("p", "", "<-p password>")
	ip       = flag.String("ip", "", "[-ip [static IP]]")
	mask     = flag.String("mask", "", "[-mask [static mask]]")
)

func main() {
	flag.Parse()
	if *username == "" || *password == "" {
		flag.PrintDefaults()
		return
	}

	ipAndMask := []string{}

	if *ip != "" {
		ipAndMask = append(ipAndMask, *ip)

		if *mask != "" {
			ipAndMask = append(ipAndMask, *mask)
		}
	}

	controller := client.NewController("config.toml")
	err := controller.ReadConfig()
	if err != nil {
		fmt.Println("fail to read config", err)
	}

	err = controller.LoginAndSetupN2NEdge(*username, *password, ipAndMask...)
	if err != nil {
		log.Println(err)
		return
	}

	c := make(chan bool, 1)
	go waitForExit(&c)
	defer close(controller.ErrChan)

	select {
	case <-c:
		log.Println("user close")
		err = controller.Disconnect()
		if err != nil {
			log.Println(err)
		}
	case err = <-controller.ErrChan:
		log.Println(err)
	}

}

func waitForExit(c *chan bool) {
	fmt.Println("input `exit()` then enter to exit")
	s := ""
	for {
		fmt.Scanf("%s", &s)
		if s == "exit()" {
			break
		}
	}
	*c <- true
}
