package user

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`

	// used for Dhcp server
	MacAddr string `json:"mac_addr"`
}
