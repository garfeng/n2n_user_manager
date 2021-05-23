package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/garfeng/n2n_user_manager/common/httputils"
	"github.com/garfeng/n2n_user_manager/common/n2n"
	"github.com/garfeng/n2n_user_manager/common/user"
)

type Config struct {
	ServerHost string `toml:"server_host"`
	MacAddr    string `toml:"mac_addr"`
	EdgePath   string `toml:"edge_path"`
}

func NewController(cfgPath string) *Controller {
	return &Controller{
		Config:     nil,
		ConfigPath: cfgPath,
		ErrChan:    make(chan error, 1),
		cmd:        nil,
	}
}

type Controller struct {
	Config     *Config
	ConfigPath string

	ErrChan chan error

	cmd *exec.Cmd

	currentUserInfo *CurrentUserInfo
}

type CurrentUserInfo struct {
	Username  string
	Password  string
	IpAndMask []string
}

func (c *Controller) SetConfigPath(configPath string) {
	c.ConfigPath = configPath
}

func (c *Controller) ReadConfig() error {
	c.Config = new(Config)
	_, err := toml.DecodeFile(c.ConfigPath, c.Config)
	return err
}

func (c *Controller) WriteConfig() error {
	w := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(w)
	err := encoder.Encode(c.Config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.ConfigPath, w.Bytes(), 0755)
	return err
}

func GetMacAddrs() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, v := range interfaces {
		res = append(res, v.HardwareAddr.String())
	}

	return res, nil
}

func (c *Controller) LoginAndGetN2NParam(username, password string) (*n2n.N2NParams, error) {
	loginInfo := &user.LoginInfo{
		Username: username,
		Password: password,
	}

	b, _ := json.MarshalIndent(loginInfo, "", "  ")

	fmt.Println(string(b))

	params := new(n2n.N2NParams)
	h := httputils.NewRequest(http.MethodPost,
		c.Config.ServerHost+"/api/n2n_params",
		bytes.NewBuffer(b)).JSON(params)
	if h.Err != nil {
		return nil, h.Err
	}

	return params, nil
}

func (c *Controller) InitUserInfo(username, password string, ipAndMask ...string) {
	c.currentUserInfo = &CurrentUserInfo{
		Username:  username,
		Password:  password,
		IpAndMask: ipAndMask,
	}
}

func (c *Controller) LoginAndSetupN2NEdge(username, password string, ipAndMask ...string) error {
	params, err := c.LoginAndGetN2NParam(username, password)
	if err != nil {
		return err
	}

	args := []string{
		"-l", params.SuperNodeServer,
		"-c", params.NetworkGroupName,
		"-k", params.SecretKey,
		"-m", params.MacAddr,
		params.EncodeType,
	}

	if len(ipAndMask) >= 1 {
		args = append(args, "-a", ipAndMask[0])
	}
	if len(ipAndMask) >= 2 {
		args = append(args, "-s", ipAndMask[1])
	}

	c.cmd = exec.Command(c.Config.EdgePath, args...)

	c.cmd.Stderr = os.Stderr
	c.cmd.Stdout = os.Stdout

	return c.cmd.Start()
}

/*
func (c *Controller) setup() {
	err := c.cmd.Run()
	if err != nil {
		c.ErrChan <- err
	}
}
*/

func (c *Controller) Cmd() *exec.Cmd {
	return c.cmd
}

func (c *Controller) Reconnect() error {
	fmt.Println("edge reconnect")
	if c.cmd != nil {
		c.Disconnect()
		c.cmd = nil
	}

	return c.LoginAndSetupN2NEdge(c.currentUserInfo.Username, c.currentUserInfo.Password, c.currentUserInfo.IpAndMask...)
}
