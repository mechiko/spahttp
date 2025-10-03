package selfdb

import (
	"fmt"
	"spahttp/repo/selfdb/models"

	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

func (c *DbSelf) Authenticate(email, password string) (user *models.User, err error) {
	user = &models.User{}
	col := c.dbSession.Collection("users")
	res := col.Find("email", email)
	err = res.One(&user)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, fmt.Errorf("error looking up user: %w", err)
		}
		return nil, fmt.Errorf("error looking up user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return user, nil
}
