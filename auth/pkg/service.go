package pkg

type Auth interface {
	GetAuthenticated() string
}

type auth struct{}

func NewAuth() Auth {
	return &auth{}
}

func (a *auth) GetAuthenticated() string {
	return "john doe"
}
