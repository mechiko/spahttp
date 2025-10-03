package repo

import (
	"errors"
	"fmt"
	"spahttp/repo/selfdb"

	"github.com/mechiko/dbscan"
)

// if err is nil then must after Lock launch UnLock
// всегда или открывает базу и проверяет объект или возвращает ошибку
func (r *Repository) LockOther() (*selfdb.DbSelf, error) {
	info := r.dbs.Info(dbscan.Other)
	if info == nil || !info.Exists {
		return nil, fmt.Errorf("%s lock info %v is nil or not exists", modError, dbscan.Other)
	}
	mu, ok := r.dbMutex[dbscan.Other]
	if ok {
		mu.mutex.Lock()
		// ensure we don't leak the lock on panic
		defer func() {
			if r := recover(); r != nil {
				mu.mutex.Unlock()
				panic(r)
			}
		}()
	} else {
		return nil, fmt.Errorf("%s lock not present mutex %v", modError, dbscan.Other)
	}
	db, err := selfdb.New(info)
	if err != nil {
		mu.mutex.Unlock()
		return nil, fmt.Errorf("%s lock open %v error %w", modError, dbscan.Other, err)
	}
	return db, nil
}

func (r *Repository) UnlockOther(db *selfdb.DbSelf) (retErr error) {
	var errClose error
	if db == nil {
		errClose = fmt.Errorf("%s unlock db %v is nil", modError, dbscan.Other)
	} else {
		errClose = db.Close()
	}
	mu, ok := r.dbMutex[dbscan.Other]
	if ok {
		defer func() {
			if rec := recover(); rec != nil {
				retErr = errors.Join(errClose, fmt.Errorf("%s unlock panic: %v", modError, rec))
			}
		}()
		mu.mutex.Unlock()
	} else {
		errUnlock := fmt.Errorf("%s unlock not present mutex %v", modError, dbscan.Other)
		return errors.Join(errClose, errUnlock)
	}
	return errors.Join(retErr, errClose)
}
