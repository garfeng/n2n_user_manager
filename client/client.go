package client

import (
	"bytes"
	"encoding/json"
	"errors"
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

type Controller struct {
	Config     *Config
	ConfigPath string
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

func (c *Controller) GetMacAddr() (string, error) {
	addrs, err := GetMacAddrs()
	if err != nil {
		return "", err
	}

	if len(addrs) == 0 {
		return "", errors.New("no mac addr available")
	}

	for _, v := range addrs {
		if v == c.Config.MacAddr && v != "" {
			return v, nil
		}
	}

	for _, v := range addrs {
		if v != "" {
			c.Config.MacAddr = v
			err = c.WriteConfig()
			if err != nil {
				return "", err
			}
			return c.Config.MacAddr, nil
		}
	}

	return "", errors.New("no mac addr available")
}

func (c *Controller) LoginAndGetN2NParam(username, password string) (*n2n.N2NParams, error) {
	macAddr, err := c.GetMacAddr()
	if err != nil {
		return nil, err
	}

	loginInfo := &user.LoginInfo{
		Username: username,
		Password: password,
		MacAddr:  macAddr,
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

func (c *Controller) LoginAndSetupN2NEdge(username, password string) error {
	params, err := c.LoginAndGetN2NParam(username, password)
	if err != nil {
		return err
	}

	cmd := exec.Command(c.Config.EdgePath,
		"-l", params.SuperNodeServer,
		"-c", params.NetworkGroupName,
		"-k", params.SecretKey,
		params.EncodeType,
		"-s", params.SubnetMask,
		"-a", params.IP)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Start()
}
