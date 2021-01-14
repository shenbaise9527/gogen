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
	maxIdleConns                       = 64
	maxOpenConns                       = 64
	maxLifetime                        = time.Minute
	CacheSafeGapBetweenIndexAndPrimary = time.Second * 5
)

var (
	exclusiveCalls = syncx.NewSharedCalls()
	stats          = cache.NewCacheStat("dbc")
	dbOnce         sync.Once
)

type (
	// DBConn gorm db.
	DBConn struct {
		*gorm.DB
	}

	// CachedDBConn with cache.
	CachedDBConn struct {
		conn  DBConn
		cache cache.Cache
	}

	QueryFn func(conn DBConn, v interface{}) error
	ExecFn  func(conn DBConn) (int64, error)
)

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

// NewCachedDBConn with cache.
func NewCachedDBConn(datasource string, c cache.CacheConf, opts ...cache.Option) CachedDBConn {
	cc := CachedDBConn{
		conn:  NewDBConn(datasource),
		cache: cache.NewCache(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}

	return cc
}

// QueryRow single row with cache.
func (cc CachedDBConn) QueryRow(v interface{}, key string, query QueryFn) error {
	return cc.cache.Take(v, key, func(v interface{}) error {
		return query(cc.conn, v)
	})
}

// Exec exec with cache.
func (cc CachedDBConn) Exec(exec ExecFn, keys ...string) (int64, error) {
	rowsaffected, err := exec(cc.conn)
	if err != nil {
		return rowsaffected, err
	}

	err = cc.cache.DelCache(keys...)

	return rowsaffected, err
}

// Transact start a transaction.
func (cc CachedDBConn) Transact(fn func(DBConn) error) error {
	return cc.conn.Transaction(func(tx *gorm.DB) error {
		db := DBConn{tx}
		return fn(db)
	})
}
`
