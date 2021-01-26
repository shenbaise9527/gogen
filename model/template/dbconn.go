package template

var DBConn = `package {{.pkg}}

import (
	"database/sql"
	"errors"
	"io"
	"time"

	"github.com/tal-tech/go-zero/core/breaker"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/syncx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxIdleConns                       = 64
	maxOpenConns                       = 64
	maxLifetime                        = time.Minute
	cacheSafeGapBetweenIndexAndPrimary = time.Second * 5
)

var (
	exclusiveCalls = syncx.NewSharedCalls()
	stats          = cache.NewCacheStat("dbc")
	connManager    = syncx.NewResourceManager()
)

type (
	// DBConn gorm db.
	DBConn struct {
		*gorm.DB
		breaker.Breaker
	}

	// CachedDBConn with cache.
	CachedDBConn struct {
		conn  *DBConn
		cache cache.Cache
	}

	QueryFn func(conn *DBConn, v interface{}) error
	ExecFn  func(conn *DBConn) (int64, error)
	PrimaryFn func(data interface{}) (string, error)
)

// NewDBConn new gorm object.
func NewDBConn(datasource string) *DBConn {
	val, err := connManager.GetResource(datasource, func() (io.Closer, error) {
		db, err := newDBConnection(datasource)
		if err != nil {
			return nil, err
		}

		return db, nil
	})
	if err != nil {
		return nil
	}

	return val.(*DBConn)
}

func newDBConnection(datasource string) (*DBConn, error) {
	db := &DBConn{Breaker: breaker.NewBreaker()}
	conn, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetConnMaxLifetime(maxLifetime)
	db.DB, err = gorm.Open(mysql.New(mysql.Config{Conn: conn}), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Transact start a transaction.
func (conn *DBConn) Transact(fn func(*DBConn) error) error {
	return conn.DoWithAcceptable(
		func() error {
			return conn.Transaction(func(tx *gorm.DB) error {
				db := &DBConn{DB: tx}
				return fn(db)
			})
		}, conn.Acceptable)
}

// Close closer interface
func (conn *DBConn) Close() error {
	return nil
}

// Acceptable accept
func (conn *DBConn) Acceptable(err error) bool {
	ok := err == nil || errors.Is(err, sql.ErrNoRows) || errors.Is(err, sql.ErrTxDone)
	return ok
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
		err := cc.conn.DoWithAcceptable(func() error {
			return query(cc.conn, v)
		}, cc.conn.Acceptable)

		return err
	})
}

// QueryRowIndex single row with cache by unique index.
func (cc CachedDBConn) QueryRowIndex(key string, primaryValue, data interface{}, query QueryFn, primaryFn PrimaryFn, queryByPrimary QueryFn) error {
	var found bool
	err := cc.cache.TakeWithExpire(primaryValue, key, func(v interface{}, expire time.Duration) error {
		err := cc.conn.DoWithAcceptable(func() error {
			return query(cc.conn, data)
		}, cc.conn.Acceptable)
		if err != nil {
			return nil
		}

		found = true
		primaryKey, err := primaryFn(data)
		if err != nil {
			return err
		}

		return cc.cache.SetCacheWithExpire(primaryKey, data, expire+cacheSafeGapBetweenIndexAndPrimary)
	})

	if err != nil {
		return err
	}

	if found {
		return nil
	}

	primaryKey, err := primaryFn(primaryValue)
	if err != nil {
		return err
	}

	err = cc.QueryRow(data, primaryKey, queryByPrimary)

	return err
}

// Exec exec with cache.
func (cc CachedDBConn) Exec(exec ExecFn, keys ...string) (int64, error) {
	var rowsaffected int64
	var err error
	err = cc.conn.DoWithAcceptable(
		func() error {
			rowsaffected, err = exec(cc.conn)

			return err
		}, cc.conn.Acceptable)
	if err != nil {
		return rowsaffected, err
	}

	err = cc.cache.DelCache(keys...)

	return rowsaffected, err
}

// Transact start a transaction.
func (cc CachedDBConn) Transact(fn func(*DBConn) error) error {
	return cc.conn.Transact(fn)
}
`
