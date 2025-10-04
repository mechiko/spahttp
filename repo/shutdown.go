package repo

import "fmt"

func (r *Repository) Shutdown() error {
	if r.dbOther != nil {
		return r.dbOther.Close()
	}
	return fmt.Errorf("dbOther is nil")
}
