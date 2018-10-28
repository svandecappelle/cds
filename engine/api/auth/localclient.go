package auth

import (
	"github.com/go-gorp/gorp"

	"github.com/ovh/cds/engine/api/user"
	"github.com/ovh/cds/sdk/log"
)

//LocalClient is a auth driver which store all in database
type LocalClient struct {
	dbFunc func() *gorp.DbMap
}

// Init nothing
func (c *LocalClient) Init(options interface{}) error {
	return nil
}

//Authentify check username and password
func (c *LocalClient) Authentify(username, password string) (bool, error) {
	// Load user
	u, err := user.LoadUserAndAuth(c.dbFunc(), username)
	if err != nil {
		log.Warning("Auth> Authorization failed")
		return false, err
	}

	b := user.IsCheckValid(password, u.Auth.HashedPassword)
	return b, err
}
