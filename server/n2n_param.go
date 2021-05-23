package server

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/garfeng/n2n_user_manager/common/n2n"
)

type N2NManagerServer interface {
	TryLoginAndGetParam(username, password string) (*n2n.N2NParams, error)
}

type ParamGenerator interface {
	GenerateParam(u UserInfo) (*n2n.N2NParams, error)
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

func (b *N2NManagerServer_ByChangeParams) TryLoginAndGetParam(username, password string) (*n2n.N2NParams, error) {
	userInfo, err := b.Authorizer.GetUserInfo(username, password)
	if err != nil {
		return nil, err
	}

	if userInfo.IsRestricted() {
		return nil, errors.New("user is restricted")
	}

	return b.ParamGenerator.GenerateParam(userInfo)
}

// ------------------ an example generator --------------

type ChangeKeyEveryDayGenerator struct {
	BaseKey          string
	TimePadding      time.Duration
	SuperNodeServer  string
	NetworkGroupName string
	EncodeType       string

	MacAddrInt uint64
}

func (c *ChangeKeyEveryDayGenerator) intToMacAddr(num uint64) string {
	b := fmt.Sprintf("%016x", num)
	b2 := b[len(b)-12:]
	sl := []string{}
	for i := 0; i < 6; i++ {
		sl = append(sl, b2[i*2:i*2+2])
	}
	mac := strings.Join(sl, ":")
	return mac
}

func (c *ChangeKeyEveryDayGenerator) GenerateParam(u UserInfo) (*n2n.N2NParams, error) {
	newAddr := atomic.AddUint64(&c.MacAddrInt, 1)
	now := time.Now().Add(time.Hour * c.TimePadding).Format("2006-01-02")
	hostName, _ := os.Hostname()
	keyBytes := md5.Sum([]byte(c.BaseKey + hostName + now))
	key := hex.EncodeToString(keyBytes[:])
	return &n2n.N2NParams{
		N2NBaseParams: n2n.N2NBaseParams{
			SuperNodeServer: c.SuperNodeServer,
		},

		NetworkGroupName: c.NetworkGroupName,
		SecretKey:        key,
		EncodeType:       c.EncodeType,
		MacAddr:          c.intToMacAddr(newAddr),
	}, nil
}
