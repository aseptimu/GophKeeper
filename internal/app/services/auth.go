package services

type Auth interface {
	Register()
	Login()
}

type UserStore interface {
	GetUserById()
	GetUserByUsername()
	SaveUser()
}
