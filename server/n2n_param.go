package server

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

type N2NManagerServer interface {
	TryLoginAndGetParam(username, password, macAddr string) (*N2NParams, error)
}

type N2NParams struct {
	N2NBaseParams
	NetworkGroupName string // 网络组名
	SecretKey        string // 加密key
	EncodeType       string // 加密方式
	SubnetMask       string // 掩码
	IP               string // 本机IP
}

type N2NBaseParams struct {
	SuperNodeServer string // 中心结点IP端口
}

type ParamGenerator interface {
	GenerateParam(u UserInfo, macAddr string) (*N2NParams, error)
}

func NewN2NManagerServer(authorizer Authorizer, generator ParamGenerator) N2NManagerServer {
	return &N2NManagerServer_ByChangeParams{
		Authorizer:     authorizer,
		ParamGenerator: generator,
	}
}

type N2NManagerServer_ByChangeParams struct {
	Authorizer     Authorizer
	ParamGenerator ParamGenerator
}

func (b *N2NManagerServer_ByChangeParams) TryLoginAndGetParam(username, password, macAddr string) (*N2NParams, error) {
	userInfo, err := b.Authorizer.GetUserInfo(username, password)
	if err != nil {
		return nil, err
	}

	if userInfo.IsRestricted() {
		return nil, errors.New("user is restricted")
	}

	return b.ParamGenerator.GenerateParam(userInfo, macAddr)
}

// ------------------ an example generator --------------

type ChangeKeyEveryDayGenerator struct {
	BaseKey          string
	TimePadding      time.Duration
	SuperNodeServer  string
	NetworkGroupName string
	EncodeType       string

	Dhcp SimpleDhcpServer
}

func (c *ChangeKeyEveryDayGenerator) GenerateParam(u UserInfo, macAddr string) (*N2NParams, error) {
	now := time.Now().Add(time.Hour * c.TimePadding).Format("2006-01-02")
	keyBytes := md5.Sum([]byte(c.BaseKey + now))
	key := hex.EncodeToString(keyBytes[:])
	ip, err := c.Dhcp.GetAnValidIp(macAddr)
	if err != nil {
		return nil, err
	}
	return &N2NParams{
		N2NBaseParams: N2NBaseParams{
			SuperNodeServer: c.SuperNodeServer,
		},

		NetworkGroupName: c.NetworkGroupName,
		SecretKey:        key,
		EncodeType:       c.EncodeType,
		SubnetMask:       c.Dhcp.SubnetMask(),
		IP:               ip,
	}, nil
}
