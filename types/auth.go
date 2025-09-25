package types

type Auth interface {
	CheckUserByEmail(email string) (bool, error)
	
}
