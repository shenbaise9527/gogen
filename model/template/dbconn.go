package template

var DBConn = `package {{.pkg}}

import (
	"database/sql"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxIdleConns = 64
	maxOpenConns = 64
	maxLifetime  = time.Minute
)

var (
	exclusiveCalls = syncx.NewSharedCalls()
	stats          = cache.NewCacheStat("dbc")
	dbOnce         sync.Once
)

type defaultResult struct {
	rowsAffected int64
	lastInsertID int64
	err          error
}

func newSqlResult(rowsAffected, lastInsertID int64, err error) sql.Result {
	return &defaultResult{
		rowsAffected: rowsAffected,
		lastInsertID: lastInsertID,
		err:          err,
	}
}

func (r *defaultResult) LastInsertId() (int64, error) {
	return r.lastInsertID, r.err
}

func (r *defaultResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

// DBConn gorm db.
type DBConn struct {
	*gorm.DB
}

// NewDBConn new gorm object.
func NewDBConn(datasource string) DBConn {
	db := DBConn{}
	dbOnce.Do(func() {
		conn, err := sql.Open("mysql", datasource)
		if err != nil {
			panic(err)
		}

		conn.SetMaxIdleConns(maxIdleConns)
		conn.SetMaxOpenConns(maxOpenConns)
		conn.SetConnMaxLifetime(maxLifetime)
		err = conn.Ping()
		if err != nil {
			panic(err)
		}

		db.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: conn}), nil)
		if err != nil {
			panic(err)
		}
	})

	return db
}

// CachedDBConn with cache.
type CachedDBConn struct {
	conn  DBConn
	cache cache.Cache
}

// NewCachedDBConn with cache.
func NewCachedDBConn(datasource string, c cache.CacheConf, opts ...cache.Option) CachedDBConn {
	cc := CachedDBConn{
		conn:  NewDBConn(datasource),
		cache: cache.NewCache(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}

	return cc
}
`
