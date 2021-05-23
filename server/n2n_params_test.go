package server

import (
	"fmt"
	"strings"
	"testing"
)

func TestIntToMacAddr(t *testing.T) {
	a := uint64(1234568)
	b := fmt.Sprintf("%016x", a)
	b2 := b[len(b)-12:]
	fmt.Println(b, b2)
	sl := []string{}
	for i := 0; i < 6; i++ {
		sl = append(sl, b2[i*2:i*2+2])
	}
	mac := strings.Join(sl, ":")
	fmt.Println(mac)
	//fmt.Println(b, b[:2], b[2:4])
}
