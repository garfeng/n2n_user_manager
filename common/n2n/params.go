package n2n

type N2NParams struct {
	N2NBaseParams
	NetworkGroupName string `json:"network_group_name"` // 网络组名
	MacAddr          string `json:"mac_addr"`           // Mac Addr
	SecretKey        string `json:"secret_key"`         // 加密key
	EncodeType       string `json:"encode_type"`        // 加密方式
}

type N2NBaseParams struct {
	SuperNodeServer string `json:"super_node_server"` // 中心结点IP端口
}
