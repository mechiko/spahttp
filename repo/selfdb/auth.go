package selfdb

import (
	"errors"
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	savedPassword := []byte(user.Passwd)
	fmt.Println(string(hashedPassword))
	fmt.Println(user.Passwd)
	fmt.Println(savedPassword)
	err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("error looking up user: %w", err)
		} else {
			return nil, err
		}
	}

	return user, nil
}
