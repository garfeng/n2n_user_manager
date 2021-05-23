package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("edge start")
	defer fmt.Println("edge end")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	select {
	case sig := <-c:
		fmt.Println("edge signal:", sig)
	case <-time.After(time.Second * 10):
		fmt.Println("edge timeout")
	}
	<-time.After(time.Second)
}
