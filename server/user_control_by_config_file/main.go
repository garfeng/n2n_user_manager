package main

import (
	"errors"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/garfeng/n2n_user_manager/server"
)

type UserInfo struct {
	Id         string `toml:"id"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Restricted bool   `toml:"restricted"`
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
			BaseKey:          "123456",
			TimePadding:      -2,
			SuperNodeServer:  "127.0.0.1:8787",
			NetworkGroupName: "myGroup",
			EncodeType:       "-A2",
			Dhcp:             server.NewDhcpServer("192.168.1.2", 50),
		},
	)

	server.SetupServer(":8080", manager)
}
