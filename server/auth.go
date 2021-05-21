package server

type Authorizer interface {
	GetUserInfo(username string, password string) (UserInfo, error)
}

type UserInfo interface {
	UID() string
	Username() string
	Email() string
	IsAdmin() bool
	IsRestricted() bool
}
