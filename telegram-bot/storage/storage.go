package storage

type Storage interface {
	SaveCredentials(c *Credentials) error
	RemoveCredentials(c *Credentials) error
	IsExistsCredentials(userName string) (bool, error)
	GetCredentials(userName string) (Credentials, error)
}

type Credentials struct {
	Username        string
	CarParkLogin    string
	CarParkPassword string
}
