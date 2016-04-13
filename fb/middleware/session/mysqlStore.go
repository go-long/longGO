package session

import (
	"github.com/gorilla/sessions"
	"jex/cn/longGo/fb/middleware/session/store/mysql"
	"database/sql"
)

type MysqlStore interface {
	Store
}




func NewMysqlStore(endpoint string, tableName string, keyPairs ...[]byte) (MysqlStore, error) {
	store, err := mysql.NewMySQLStore(endpoint,tableName,"/",86400 * 30,keyPairs...)
	if err != nil {
		return nil, err
	}
	return &mysqlStore{store}, nil
}

func NewMysQLStoreFromConnection(db *sql.DB, tableName string, keyPairs ...[]byte) (MysqlStore, error) {
	store, err := mysql.NewMySQLStoreFromConnection(db,tableName,"/",86400 * 30,keyPairs...)
	if err != nil {
		return nil, err
	}
	return &mysqlStore{store}, nil
}

type mysqlStore struct {
	*mysql.MySQLStore
}

func (c *mysqlStore) Options(options Options) {
	c.MySQLStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}