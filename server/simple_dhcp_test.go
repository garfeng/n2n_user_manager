package server

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestBaseSimpleDhcpServer_GetAnValidIp(t *testing.T) {
	b := &BaseSimpleDhcpServer{
		Start:    net.IP{192, 168, 1, 23},
		Number:   100,
		pool:     sync.Map{},
		poolSize: 0,
	}

	fmt.Println(b.Start.String())

	b._start = b.IpToInt(b.Start.String())

	fmt.Println(b._start)

	fmt.Println(b.IntToIp(b._start))

	for i := 0; i < 10; i++ {
		ip, err := b.GetAnValidIp("hello")
		fmt.Println(ip, err)
	}

	fmt.Println(b.SubnetMask())
}
