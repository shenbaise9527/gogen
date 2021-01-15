package schemas

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// GetTableInfos 获取表.
func GetTableInfos(tablenames []string, url string, withCached, withTracing bool) ([]*Table, error) {
	// parse dsn
	dsn, err := mysql.ParseDSN(url)
	if err != nil {
		return nil, err
	}

	// 创建数据库连接.
	databaseSource := strings.TrimSuffix(url, "/"+dsn.DBName) + "/information_schema"
	conn, err := sql.Open("mysql", databaseSource)
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	tables := make([]*Table, 0, len(tablenames))
	for _, name := range tablenames {
		t, err := GetTableInfoByName(name, dsn.DBName, conn, withCached, withTracing)
		if err != nil {
			return nil, err
		}

		tables = append(tables, t)
	}

	return tables, nil
}

// 获取表信息.
func GetTableInfoByName(name, dbname string, conn *sql.DB, withCached, withTracing bool) (*Table, error) {
	// 查询表注释.
	query := "select TABLE_COMMENT from TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	row := conn.QueryRow(query, dbname, name)
	var comment string
	err := row.Scan(&comment)
	if err != nil {
		return nil, err
	}

	t := NewTable(name, comment, withCached, withTracing)

	// 查询字段信息.
	query = "SELECT COLUMN_NAME, DATA_TYPE, IFNULL(CHARACTER_MAXIMUM_LENGTH, 0) CHARACTER_MAXIMUM_LENGTH, IFNULL(NUMERIC_PRECISION, 0) NUMERIC_PRECISION, IFNULL(NUMERIC_SCALE, 0) NUMERIC_SCALE, COLUMN_KEY, EXTRA, COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	rows, err := conn.Query(query, dbname, name)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var col Column
		err = rows.Scan(&col.Name, &col.DataType, &col.CharacterMaximumLength, &col.NumericPrecision, &col.NumericScale, &col.ColumnKey, &col.Extra, &col.ColumnComment)
		if err != nil {
			return nil, err
		}

		t.AddColumn(&col)
	}

	// 查询索引信息.
	query = "SELECT NON_UNIQUE, INDEX_NAME, COLUMN_NAME FROM STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?"
	indexrows, err := conn.Query(query, dbname, name)
	if err != nil {
		return nil, err
	}

	defer indexrows.Close()
	var nonUnique int64
	var indexName, indexColumnName string
	indexs := make(map[string]*Index)
	for indexrows.Next() {
		err = indexrows.Scan(&nonUnique, &indexName, &indexColumnName)
		if err != nil {
			return nil, err
		}

		col := t.GetColumn(indexColumnName)
		if col == nil {
			return nil, fmt.Errorf("cant find column: %s", indexColumnName)
		}

		ix, ok := indexs[indexName]
		if !ok {
			ix = NewIndex(indexName, nonUnique)
			indexs[indexName] = ix
		}

		ix.AddColumn(col)
	}

	for _, ix := range indexs {
		t.AddIndex(ix)
	}

	return t, nil
}
