package server

import (
	"errors"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// TODO: remove the dhcp feature because of it was already support in n2n

type SimpleDhcpServer interface {
	GetAnValidIp(macAddr string) (string, error)
	ReleaseAnIp(ip string, macAddr string)
	RenewIp(ip string, macAddr string) error
	SubnetMask() string
}

type BaseSimpleDhcpServer struct {
	Start  net.IP
	Number int

	_start   uint32
	pool     sync.Map
	poolSize int64
}

type BaseIpUsedInfo struct {
	LastConnectTime time.Time

	// TODO: Delete ip from pool when expires
	Expire  time.Time
	MacAddr string
}

func NewBaseIpUsedInfo(macAddr string) *BaseIpUsedInfo {
	b := &BaseIpUsedInfo{
		MacAddr: macAddr,
	}
	b.Renew()
	return b
}

func (b *BaseIpUsedInfo) Renew() *BaseIpUsedInfo {
	b.LastConnectTime = time.Now()
	// 续订时间为1天，至少需要每天续一次
	b.Expire = b.LastConnectTime.Add(time.Hour * 24)
	return b
}

func (b *BaseIpUsedInfo) Duration() time.Duration {
	return time.Now().Sub(b.LastConnectTime)
}

func (b *BaseSimpleDhcpServer) GetAnValidIp(macAddr string) (string, error) {
	if b.poolSize >= int64(b.Number) {
		return "", errors.New("ip pool is full")
	}
	ipUsedInfo := NewBaseIpUsedInfo(macAddr)

	id := rand.Int() % b.Number
	ip := b._start + uint32(id)
	ipString := b.IntToIp(ip)
	if _, ok := b.pool.LoadOrStore(ipString, ipUsedInfo); ok {
		return b.GetAnValidIp(macAddr)
	}

	atomic.AddInt64(&(b.poolSize), 1)
	return ipString, nil
}

func (b *BaseSimpleDhcpServer) GetIntByte(i uint32, byt int) byte {
	return byte((i >> (8 * byt)) & 0xFF)
}

func (b *BaseSimpleDhcpServer) IntToIp(ip uint32) string {
	ipByte := net.IP{
		b.GetIntByte(ip, 3),
		b.GetIntByte(ip, 2),
		b.GetIntByte(ip, 1),
		b.GetIntByte(ip, 0),
	}

	return ipByte.String()
}

// 续订IP，如果返回错误，则需要重新订阅新的IP
func (b *BaseSimpleDhcpServer) RenewIp(ip string, macAddr string) error {
	one, ok := b.pool.Load(ip)
	if !ok {
		return errors.New("ip not in pool")
	}

	ipInfo, ok := one.(*BaseIpUsedInfo)
	if !ok {
		atomic.AddInt64(&(b.poolSize), -1)
		b.pool.Delete(ip)
		return errors.New("ip not valid")
	}
	if ipInfo.MacAddr != macAddr {
		return errors.New("forbidden to this ip")
	}

	ipInfo.Renew()

	return nil
}

func (b *BaseSimpleDhcpServer) ReleaseAnIp(ip string, macAddr string) {
	one, ok := b.pool.Load(ip)
	if !ok {
		return
	}
	ipInfo, ok := one.(*BaseIpUsedInfo)
	if !ok {

		b.pool.Delete(ip)
		return
	}

	if ipInfo.MacAddr != macAddr {
		return
	}
	b.pool.Delete(ip)
}

func (b *BaseSimpleDhcpServer) SubnetMask() string {
	s := b._start
	e := b._start + uint32(b.Number)

	vp := ^(s ^ e)

	vp2 := uint32(0)

	ok := true

	for i := 31; i >= 0; i-- {
		vpi := (vp >> i) & 1
		if vpi == 1 {

		} else {
			ok = false
		}

		if ok {
			vp2 += 1 << i
		} else {
			//vp2 = vp2 << 1
		}
	}

	return b.IntToIp(vp2)
}

func (b *BaseSimpleDhcpServer) IpToInt(ip string) uint32 {
	ipNet := net.ParseIP(ip)

	res := uint32(0)
	l := len(ipNet)
	for i := 0; i < 4; i++ {
		res += (uint32(ipNet[l-4+i]) << ((3 - i) * 8))
	}

	return res
}

func NewDhcpServer(startIp string, number int) *BaseSimpleDhcpServer {
	d := &BaseSimpleDhcpServer{
		Start:    net.ParseIP(startIp),
		Number:   number,
		_start:   0,
		pool:     sync.Map{},
		poolSize: 0,
	}
	d._start = d.IpToInt(d.Start.String())
	return d
}
