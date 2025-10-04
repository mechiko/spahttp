package selfdb

import (
	_ "embed"
	"fmt"

	"github.com/mechiko/dbscan"
	"github.com/upper/db/v4"
)

const modError = "selfdb"

type DbSelf struct {
	dbSession db.Session // открытый хэндл тут
	dbInfo    *dbscan.DbInfo
	infoType  dbscan.DbInfoType
	version   int64
}

func New(info *dbscan.DbInfo) (*DbSelf, error) {
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	db := &DbSelf{
		dbInfo:   info,
		infoType: dbscan.Other,
	}
	if err := db.Connect(); err != nil {
		return nil, fmt.Errorf("%s error check %w", modError, err)
	}
	return db, nil
}

// при первом запуске программы проверяет наличие создает пустую и миграция
func NewOnceOnStart(info *dbscan.DbInfo) (*DbSelf, error) {
	if info == nil {
		return nil, fmt.Errorf("%s dbinfo is nil", modError)
	}
	db := &DbSelf{
		dbInfo:   info,
		infoType: dbscan.Other,
	}
	// передаем флаг о необходимости создания, это при запуске приложения из repo
	// проверяем, если нет создаем, если надо мигрируем
	// открываем сесиию в этом методе если нет ошибки
	if err := db.CheckCreateAndMigrate(); err != nil {
		return nil, fmt.Errorf("%s error check %w", modError, err)
	}
	return db, nil
}

func (c *DbSelf) Close() (err error) {
	if c.dbSession == nil {
		return nil
	}
	return c.dbSession.Close()
}

func (c *DbSelf) Sess() db.Session {
	return c.dbSession
}

func (c *DbSelf) Version() int64 {
	return c.version
}

func (c *DbSelf) Info() dbscan.DbInfo {
	if c.dbInfo == nil {
		return dbscan.DbInfo{}
	}
	return *c.dbInfo
}

func (c *DbSelf) InfoType() dbscan.DbInfoType {
	return c.infoType
}
