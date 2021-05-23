package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/garfeng/n2n_user_manager/common/httputils"
	"github.com/garfeng/n2n_user_manager/server"
)

type HostConfig struct {
	GiteaHost        string `toml:"gitea_host"`
	SuperNodeServer  string `toml:"super_node_server"`
	NetworkGroupName string `toml:"network_group_name"`
	EncodeType       string `toml:"encode_type"`

	TimePadding time.Duration `toml:"time_padding"` // hour, eg. -2, key will be updated at 2:00am
	BaseKey     string        `toml:"base_key"`

	DhcpStartIp  string `toml:"dhcp_start_ip"`
	DhcpIpNumber int    `toml:"dhcp_ip_number"`
	ServerPort   string `toml:"server_port"`
}

var (
	_hostConfig *HostConfig // need init
)

func init() {
	_hostConfig = new(HostConfig)
	_, err := toml.DecodeFile("./config.toml", _hostConfig)
	if err != nil {
		panic(err)
	}
}

type GiteaUserInfo struct {
	Id              int    `json:"id"`
	GiteaUsername   string `json:"username"`
	GiteaEmail      string `json:"email"`
	GiteaIsAdmin    bool   `json:"is_admin"`
	GiteaRestricted bool   `json:"restricted"`
	Message         string `json:"message"`
}

func (g *GiteaUserInfo) UID() string {
	return fmt.Sprint(g.Id)
}

func (g *GiteaUserInfo) Username() string {
	return g.GiteaUsername
}
func (g *GiteaUserInfo) Email() string {
	return g.GiteaEmail
}

func (g *GiteaUserInfo) IsAdmin() bool {
	return g.GiteaIsAdmin
}
func (g *GiteaUserInfo) IsRestricted() bool {
	return g.GiteaRestricted
}

type GiteaAuthorizer struct {
	GiteaHost string
}

func (g *GiteaAuthorizer) GetUserInfo(username, password string) (server.UserInfo, error) {
	giteaUserInfo := new(GiteaUserInfo)
	h := httputils.NewRequest(http.MethodGet, g.GiteaHost+"/api/v1/user").
		SetBasicAuth(username, password).JSON(giteaUserInfo)

	if h.Err != nil {
		return nil, h.Err
	}

	if giteaUserInfo.Message != "" {
		return nil, errors.New(giteaUserInfo.Message)
	}

	return giteaUserInfo, nil
}

func main() {
	manager := server.NewN2NManagerServer(
		&GiteaAuthorizer{
			GiteaHost: _hostConfig.GiteaHost,
		}, &server.ChangeKeyEveryDayGenerator{
			BaseKey:          _hostConfig.BaseKey,
			TimePadding:      _hostConfig.TimePadding,
			SuperNodeServer:  _hostConfig.SuperNodeServer,
			NetworkGroupName: _hostConfig.NetworkGroupName,
			EncodeType:       _hostConfig.EncodeType,
			MacAddrInt:       uint64(time.Now().Unix()),
		})
	server.SetupServer(_hostConfig.ServerPort, manager)
}
