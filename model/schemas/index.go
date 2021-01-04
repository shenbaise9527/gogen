package schemas

import (
	"strings"
)

type IndexType int

const (
	InvalidIndex IndexType = iota
	PrimaryKeyIndex
	UniqueKeyIndex
	NormalKeyIndex
)

// Index 索引.
type Index struct {
	Name      string    // 索引名.
	NonUnique int64     // 是否唯一.
	Columns   []*Column // 索引对应的列.
}

// NewIndex 创建索引.
func NewIndex(name string, nonUnique int64) *Index {
	return &Index{
		Name:      name,
		NonUnique: nonUnique,
		Columns:   make([]*Column, 0),
	}
}

// Type 索引类型.
func (ix *Index) Type() IndexType {
	if 0 != ix.NonUnique {
		return NormalKeyIndex
	} else if ix.Name == "PRIMARY" {
		return PrimaryKeyIndex
	} else {
		return UniqueKeyIndex
	}
}

// AddColumn 添加列.
func (ix *Index) AddColumn(col *Column) {
	ix.Columns = append(ix.Columns, col)
}

// GetColumns 获取列.
func (ix *Index) GetColumns() map[string]*Column {
	res := make(map[string]*Column, len(ix.Columns))
	for _, col := range ix.Columns {
		res[col.Name] = col
	}

	return res
}

// GetColumnsNameByDq 获取列名,带双引号,以逗号分隔.
func (ix *Index) GetColumnsNameByDq() string {
	colnames := ix.GetColumnNameSliceByDq()
	return strings.Join(colnames, ",")
}

// GetColumnNameSliceByDq 获取列名的列表,带双引号.
func (ix *Index) GetColumnNameSliceByDq() []string {
	colnames := make([]string, 0)
	for _, col := range ix.Columns {
		colnames = append(colnames, col.GetNameByDoubleQuotation())
	}

	return colnames
}

// GetSuffixName 获取后缀名字.
func (ix *Index) GetSuffixName() string {
	return UpperStartCamel(ix.Name)
}

// GetColumnsNameAndType 获取列名和类型.
func (ix *Index) GetColumnsNameAndType() (string, error) {
	results := make([]string, 0)
	for _, col := range ix.Columns {
		name := col.GetLowerName()
		Type, err := col.ConvertType()
		if err != nil {
			return "", err
		}

		results = append(results, name+" "+Type)
	}

	return strings.Join(results, ", "), nil
}

// GetColumnsNameAndMark 获取列名+问号.
func (ix *Index) GetColumnsNameAndMark() string {
	results := make([]string, 0)
	for _, col := range ix.Columns {
		name := col.GetLowerName()
		results = append(results, name+" = ?")
	}

	return strings.Join(results, " and ")
}

// GetColumnsName 获取列名.
func (ix *Index) GetColumnsName() string {
	results := make([]string, 0, len(ix.Columns))
	for _, col := range ix.Columns {
		name := col.GetLowerName()
		results = append(results, name)
	}

	return strings.Join(results, ", ")
}

// GetColumnsExpressionValues 获取列名.
func (ix *Index) GetColumnsExpressionValues() string {
	expressionValues := make([]string, 0, len(ix.Columns))
	for _, col := range ix.Columns {
		expressionValues = append(expressionValues, "data."+col.GetUpperStartName())
	}

	return strings.Join(expressionValues, ", ")
}

// GetColumnCacheName 获取缓存使用的列名.
func (ix *Index) GetColumnCacheName() string {
	results := make([]string, 0, len(ix.Columns))
	for _, col := range ix.Columns {
		name := col.GetUpperStartName()
		results = append(results, name)
	}

	return strings.Join(results, "")
}
