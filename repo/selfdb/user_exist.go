package selfdb

import (
	"fmt"
	"spahttp/repo/selfdb/models"

	"github.com/upper/db/v4"
)

func (c *DbSelf) UserExists(field models.LookupField, value any) (bool, error) {
	var res db.Result
	var user models.User
	col := c.dbSession.Collection("users")
	switch field {
	case models.ID:
		res = col.Find("id", value)
	case models.Email:
	case models.Username:
	default:
		return false, fmt.Errorf("invalid lookup field")
	}
	err := res.One(&user)
	if err != nil {
		return false, fmt.Errorf("error lookup field %w", err)
	}
	return false, nil
}
