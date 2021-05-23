package main

import (
	"errors"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/garfeng/n2n_user_manager/server"
)

type UserInfo struct {
	Id         string `toml:"id"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Restricted bool   `toml:"restricted"`
}

type HostConfig struct {
	SuperNodeServer  string `toml:"super_node_server"`
	NetworkGroupName string `toml:"network_group_name"`
	EncodeType       string `toml:"encode_type"`

	TimePadding time.Duration `toml:"time_padding"` // hour, eg. -2, key will be updated at 2:00am
	BaseKey     string        `toml:"base_key"`

	ServerPort string `toml:"server_port"`
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

func (u *UserInfo) UID() string {
	return u.Id
}

func (u *UserInfo) Username() string {
	return u.User
}

func (u *UserInfo) Email() string {
	return ""
}

func (u *UserInfo) IsAdmin() bool {
	return false
}

func (u *UserInfo) IsRestricted() bool {
	return u.Restricted
}

type UserList struct {
	Users   []*UserInfo `toml:"users"`
	userMap map[string]*UserInfo
}

func (u *UserList) Parse() {
	u.userMap = map[string]*UserInfo{}
	for _, v := range u.Users {
		u.userMap[v.User] = v
	}
}

type FileAuthorizer struct {
	users *UserList
	// TODO: Sync safety
	mutex sync.RWMutex
}

// reload user config from file
func (l *FileAuthorizer) Refresh() error {
	_, err := toml.DecodeFile("./users.toml", l.users)
	if err != nil {
		return err
	}
	l.users.Parse()
	return nil
}

func (f *FileAuthorizer) GetUserInfo(username, password string) (server.UserInfo, error) {
	err := f.Refresh()
	if err != nil {
		return nil, err
	}

	userInfo, ok := f.users.userMap[username]
	if !ok {
		return nil, errors.New("invalid username")
	}

	if userInfo.Password != password {
		return nil, errors.New("invalid username and password")
	}

	return userInfo, nil
}

func main() {
	manager := server.NewN2NManagerServer(
		&FileAuthorizer{
			users: new(UserList),
		},
		&server.ChangeKeyEveryDayGenerator{
			BaseKey:          _hostConfig.BaseKey,
			TimePadding:      _hostConfig.TimePadding,
			SuperNodeServer:  _hostConfig.SuperNodeServer,
			NetworkGroupName: _hostConfig.NetworkGroupName,
			EncodeType:       _hostConfig.EncodeType,
			MacAddrInt:       uint64(time.Now().Unix() << 8),
		},
	)

	server.SetupServer(":8080", manager)
}
