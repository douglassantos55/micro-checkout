package pkg

type Auth interface {
	GetAuthenticated(user string) string
}

type auth struct{}

func NewAuth() Auth {
	return &auth{}
}

func (a *auth) GetAuthenticated(user string) string {
	return user
}
