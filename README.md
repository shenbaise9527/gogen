# gogen
根据数据库配置自动生成model代码(Golang语言),主要是使用[go-zero](https://github.com/tal-tech/go-zero)和[gorm](https://gorm.io/)

## 支持功能
* 支持mysql数据库,使用的是[gorm](https://gorm.io/).
* 支持数据库熔断,使用的[go-zero](https://github.com/tal-tech/go-zero)的`core/breaker`组件.
* 支持redis缓存,使用的[go-zero](https://github.com/tal-tech/go-zero)的`core/stores/cache`组件.
* 支持OpenTracing链路追踪,使用的[tracing](https://github.com/shenbaise9527/tracing).
* 支持单字段的主键,包括查询、更新和删除操作.
* 支持多字段的组合主键,包括查询、更新和删除操作.
* 支持单字段的唯一索引,包括查询、更新和删除操作.
* 支持多字段的组合唯一索引,包括查询、更新和删除操作.

## 命令
``` bash
$ ./gogen model datasource -h
NAME:
   gogen model datasource - generate model from datasource

USAGE:
   gogen model datasource [command options] [arguments...]

OPTIONS:
   --url value      data soucre of database, mysql: "root:password@tcp(127.0.0.1:3306)/database"
   --table value    the tables in the database,support for comma separation
   --dir value      the target dir
   --cache value    generate code with cache [optional]
   --tracing value  generate code with tracing [optional]
   --help, -h       show help (default: false)
```

## 例子
``` bash
$ ./gogen model datasource --url "gozero:123456@tcp(192.168.20.151:3406)/gozero" --table "book,goods,user" --dir ../models --tracing true --cache true

$ ll 
total 40K
-rwxrwxr-x 1 zhou.yingan zhou.yingan 3.9K Jan 25 16:29 bookmodel.go
-rwxrwxr-x 1 zhou.yingan zhou.yingan  826 Jan 25 16:29 builder.go
-rwxrwxr-x 1 zhou.yingan zhou.yingan 4.1K Jan 25 16:29 dbconn.go
-rwxrwxr-x 1 zhou.yingan zhou.yingan 9.8K Jan 25 16:29 goodsmodel.go
-rwxrwxr-x 1 zhou.yingan zhou.yingan  12K Jan 25 16:29 usermodel.go
```

数据库相关的封装在`dbconn.go`文件中.
``` go
package models

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

	QueryFn   func(conn *DBConn, v interface{}) error
	ExecFn    func(conn *DBConn) (int64, error)
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
```

`bookmodel.go`对应`book`表,`goodsmodel.go`对应`goods`表,`usermodel.go`对应`user`表.

`goodsmodel.go`的内容:
``` go
package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/shenbaise9527/tracing"
)

var (
	goodsFieldNames = fieldNames(&Goods{})

	goodsRowsNoPA = strings.Join(removeField(goodsFieldNames, "GOODSID", "GOODSCODE"), "=?,") + "=?"

	goodsRowsUkGoodsGoodsnameNoPA = strings.Join(removeField(goodsFieldNames, "GOODSNAME", "GOODSID", "GOODSCODE"), "=?,") + "=?"

	cacheGoodsPKPrefix = "cache#Goods#PK"

	cacheGoodsUkGoodsGoodsnamePrefix = "cache#Goods#UkGoodsGoodsname"
)

type (
	// GoodsModel model interface
	GoodsModel interface {
		// Insert insert the record
		Insert(ctx context.Context, data *Goods) error

		// FindOne query the record by the primary key
		FindOne(ctx context.Context, goodsid int64, goodscode string) (*Goods, error)

		// FindByUkGoodsGoodsname query the record by the unique key-UK_GOODS_GOODSNAME
		FindByUkGoodsGoodsname(ctx context.Context, goodsname string) (*Goods, error)

		// Update update the record by the primary key
		Update(ctx context.Context, data *Goods) error

		// UpdateByUkGoodsGoodsname update the record by the unique key-UK_GOODS_GOODSNAME
		UpdateByUkGoodsGoodsname(ctx context.Context, data *Goods) error

		// Delete delete the record
		Delete(ctx context.Context, data *Goods) error

		// Delete delete the record by the primary key
		DeleteByPrimary(ctx context.Context, goodsid int64, goodscode string) error

		// DeleteByUkGoodsGoodsname delete the record by the unique key-UK_GOODS_GOODSNAME
		DeleteByUkGoodsGoodsname(ctx context.Context, goodsname string) error
	}

	// defaultGoodsModel model object
	defaultGoodsModel struct {
		CachedDBConn
		table string
	}

	// Goods .
	Goods struct {
		Goodsid       int64     `json:"goodsid" gorm:"column:GOODSID;autoIncrement;primaryKey"`
		Goodscode     string    `json:"goodscode" gorm:"column:GOODSCODE;primaryKey"` // ()
		Goodsname     string    `json:"goodsname" gorm:"column:GOODSNAME"`
		Marketid      int64     `json:"marketid" gorm:"column:MARKETID"`         // ID
		Goodsgroupid  int64     `json:"goodsgroupid" gorm:"column:GOODSGROUPID"` // ID
		Goodsstatus   int32     `json:"goodsstatus" gorm:"column:GOODSSTATUS"`   // - 1: 2: 3: 4: 5: 6: 7:
		Currencyid    int64     `json:"currencyid" gorm:"column:CURRENCYID"`     // ID
		Goodunitid    int64     `json:"goodunitid" gorm:"column:GOODUNITID"`     // ID
		Agreeunit     float64   `json:"agreeunit" gorm:"column:AGREEUNIT"`
		Decimalplace  int32     `json:"decimalplace" gorm:"column:DECIMALPLACE"`
		Listingdate   time.Time `json:"listingdate" gorm:"column:LISTINGDATE"`
		Lasttradedate time.Time `json:"lasttradedate" gorm:"column:LASTTRADEDATE"` // ()
	}

	// goodsPrimary primary key struct.
	goodsPrimary struct {
		Goodsid   int64  `json:"goodsid" gorm:"column:GOODSID;autoIncrement;primaryKey"`
		Goodscode string `json:"goodscode" gorm:"column:GOODSCODE;primaryKey"` // ()
	}
)

// TableName is goods
func (Goods) TableName() string {
	return "goods"
}

// NewGoodsModel new model object
func NewGoodsModel(conn CachedDBConn) GoodsModel {
	return &defaultGoodsModel{
		CachedDBConn: conn,
		table:        "goods",
	}
}

// Insert insert the record
func (m *defaultGoodsModel) Insert(ctx context.Context, data *Goods) error {
	var err error
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "Insert")
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	err = m.conn.DoWithAcceptable(
		func() error {
			err := m.conn.Create(data).Error

			return err
		}, m.conn.Acceptable)

	return err
}

// FindOne query the record by the primary key
func (m *defaultGoodsModel) FindOne(ctx context.Context, goodsid int64, goodscode string) (*Goods, error) {
	var err error
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, goodsid, goodscode)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "FindOne")
	span.SetTag("key", primaryKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	var resp Goods
	err = m.QueryRow(&resp, primaryKey, func(conn *DBConn, v interface{}) error {
		return conn.Where("goodsid = ? and goodscode = ?", goodsid, goodscode).First(v).Error
	})

	return &resp, err
}

// FindByUkGoodsGoodsname query the record by the unique key-UK_GOODS_GOODSNAME
func (m *defaultGoodsModel) FindByUkGoodsGoodsname(ctx context.Context, goodsname string) (*Goods, error) {
	var err error
	ukgoodsgoodsnameKey := fmt.Sprintf("%s#%v", cacheGoodsUkGoodsGoodsnamePrefix, goodsname)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "FindByUkGoodsGoodsname")
	span.SetTag("key", ukgoodsgoodsnameKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	var data Goods

	var primaryKey goodsPrimary
	err = m.QueryRowIndex(
		ukgoodsgoodsnameKey, &primaryKey, &data,
		func(conn *DBConn, v interface{}) error {
			return conn.Where("goodsname = ?", goodsname).First(v).Error
		},
		m.buildPrimaryKey,
		func(conn *DBConn, v interface{}) error {
			return conn.Where("goodsid = ? and goodscode = ?", primaryKey.Goodsid, primaryKey.Goodscode).First(v).Error
		},
	)

	return &data, err
}

func (m *defaultGoodsModel) buildPrimaryKey(v interface{}) (string, error) {
	switch d := v.(type) {
	case *Goods:
		return fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, d.Goodsid, d.Goodscode), nil
	case *goodsPrimary:
		return fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, d.Goodsid, d.Goodscode), nil
	default:
		return "", fmt.Errorf("cant support type: %v", d)
	}
}

// Update update the record by the primary key
func (m *defaultGoodsModel) Update(ctx context.Context, data *Goods) error {
	var err error
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, data.Goodsid, data.Goodscode)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "Update")
	span.SetTag("key", primaryKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	_, err = m.Exec(func(conn *DBConn) (int64, error) {
		query := fmt.Sprintf("update %s set %s where goodsid = ? and goodscode = ?", m.table, goodsRowsNoPA)
		db := conn.Exec(query, data.Goodsname, data.Marketid, data.Goodsgroupid, data.Goodsstatus, data.Currencyid, data.Goodunitid, data.Agreeunit, data.Decimalplace, data.Listingdate, data.Lasttradedate, data.Goodsid, data.Goodscode)

		return db.RowsAffected, db.Error
	}, primaryKey)

	return err
}

// UpdateByUkGoodsGoodsname update the record by the unique key-UK_GOODS_GOODSNAME
func (m *defaultGoodsModel) UpdateByUkGoodsGoodsname(ctx context.Context, data *Goods) error {
	var err error
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, data.Goodsid, data.Goodscode)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "UpdateByUkGoodsGoodsname")
	span.SetTag("key", primaryKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	_, err = m.Exec(func(conn *DBConn) (int64, error) {
		query := fmt.Sprintf("update %s set %s where goodsname = ?", m.table, goodsRowsUkGoodsGoodsnameNoPA)
		db := conn.Exec(query, data.Marketid, data.Goodsgroupid, data.Goodsstatus, data.Currencyid, data.Goodunitid, data.Agreeunit, data.Decimalplace, data.Listingdate, data.Lasttradedate, data.Goodsname)

		return db.RowsAffected, db.Error
	}, primaryKey)

	return err
}

// Delete delete the record
func (m *defaultGoodsModel) Delete(ctx context.Context, data *Goods) error {
	var err error
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, data.Goodsid, data.Goodscode)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "Delete")
	span.SetTag("key", primaryKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	err = m.delete(data)

	return err
}

// Delete delete the record by the primary key
func (m *defaultGoodsModel) DeleteByPrimary(ctx context.Context, goodsid int64, goodscode string) error {
	var err error
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, goodsid, goodscode)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "DeleteByPrimary")
	span.SetTag("key", primaryKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	data, err := m.FindOne(ctx, goodsid, goodscode)
	if err != nil {
		return err
	}

	err = m.delete(data)

	return err
}

// DeleteByUkGoodsGoodsname delete the record by the unique key-UK_GOODS_GOODSNAME
func (m *defaultGoodsModel) DeleteByUkGoodsGoodsname(ctx context.Context, goodsname string) error {
	var err error
	ukgoodsgoodsnameKey := fmt.Sprintf("%s#%v", cacheGoodsUkGoodsGoodsnamePrefix, goodsname)
	span := tracing.ChildOfSpanFromContext(ctx, "goodsmodel")
	defer span.Finish()
	ext.DBStatement.Set(span, "DeleteByUkGoodsGoodsname")
	span.SetTag("key", ukgoodsgoodsnameKey)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	data, err := m.FindByUkGoodsGoodsname(ctx, goodsname)
	if err != nil {
		return err
	}

	err = m.delete(data)

	return err
}

func (m *defaultGoodsModel) delete(data *Goods) error {
	primaryKey := fmt.Sprintf("%s#%v#%v", cacheGoodsPKPrefix, data.Goodsid, data.Goodscode)

	ukgoodsgoodsnameKey := fmt.Sprintf("%s#%v", cacheGoodsUkGoodsGoodsnamePrefix, data.Goodsname)

	_, err := m.Exec(func(conn *DBConn) (int64, error) {
		db := conn.Delete(Goods{}, "goodsid = ? and goodscode = ?", data.Goodsid, data.Goodscode)

		return db.RowsAffected, db.Error
	}, primaryKey, ukgoodsgoodsnameKey)

	return err
}
```