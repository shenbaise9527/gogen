package schemas

import (
	"fmt"
	"strings"
)

// Column 列.
type Column struct {
	Name                   string // 列名.
	DataType               string // 类型.
	CharacterMaximumLength int64  // 字符串长度.
	NumericPrecision       int64  // 整数部分长度.
	NumericScale           int64  // 精度.
	ColumnKey              string // key,PRI/UNI
	Extra                  string // auto_increment
	ColumnComment          string // 注释.
}

// IsAutoIncrement 是否是自增字段.
func (col *Column) IsAutoIncrement() bool {
	return col.Extra == "auto_increment"
}

// IsTime 是否为time.Time类型.
func (col *Column) IsTime() bool {
	dt, _ := col.ConvertType()
	return dt == "time.Time"
}

// ConvertType 转换类型.
func (col *Column) ConvertType() (string, error) {
	dt := strings.ToLower(col.DataType)
	switch dt {
	case "bool", "boolean":
		return "bool", nil
	case "tinyint", "smallint", "mediumint", "int":
		return "int32", nil
	case "integer", "bigint":
		return "int64", nil
	case "float", "double":
		return "float64", nil
	case "decimal":
		if 0 == col.NumericScale && (0 == col.NumericPrecision || col.NumericPrecision >= 10) {
			return "int64", nil
		} else if 0 == col.NumericScale && col.NumericPrecision < 10 {
			return "int32", nil
		} else {
			return "float64", nil
		}
	case "date", "datetime", "timestamp":
		return "time.Time", nil
	case "time":
		return "string", nil
	case "year":
		return "int64", nil
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext", "enum", "set", "json":
		return "string", nil
	default:
		return "", fmt.Errorf("unexpected database type: %s", col.DataType)
	}
}

// GetUpperNameByDoubleQuotation 获取名字,带双引号.
func (col *Column) GetUpperNameByDoubleQuotation() string {
	return `"` + col.GetUpperName() + `"`
}

// GetUpperStartName 获取名字.
func (col *Column) GetUpperStartName() string {
	return UpperStartCamel(col.Name)
}

// GetUpperName 获取大写名字.
func (col *Column) GetUpperName() string {
	return strings.ToUpper(col.Name)
}

// GetLowerName 获取小写名字.
func (col *Column) GetLowerName() string {
	return strings.ToLower(col.Name)
}
