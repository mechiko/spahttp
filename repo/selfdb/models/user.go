package models

type User struct {
	Id      int    `db:"id,omitempty"`
	Login   string `db:"login"`
	Passwd  string `db:"passwd"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Active  int    `db:"active"`
	IsAdmin int    `db:"is_admin"`
	Rem     string `db:"rem"`
}

// LookupField represents an enumeration used to specify fields for lookup operations in user-related database queries.
type LookupField int

const (
	ID LookupField = iota
	Email
	Username
)
