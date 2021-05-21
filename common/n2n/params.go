package n2n

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
