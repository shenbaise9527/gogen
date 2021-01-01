package schemas

import (
	"fmt"
	"strings"
)

// Table 表.
type Table struct {
	Name                string             // 表名.
	TableComment        string             // 表注释.
	Columns             []*Column          // 所有字段.
	Name2Column         map[string]*Column // 字段名->字段.
	PrimaryIndex        *Index             // 主键索引.
	UniqueIndex         []*Index           // 唯一索引.
	NormalIndex         []*Index           // 普通索引.
	AutoIncrementColumn *Column            // 自增字段.
	IsContainTime       bool               // 是否包含time.Time类型.
}

// NewTable 新建表对象.
func NewTable(name, comment string) *Table {
	return &Table{
		Name:         name,
		TableComment: comment,
		Columns:      make([]*Column, 0),
		Name2Column:  make(map[string]*Column),
		UniqueIndex:  make([]*Index, 0),
		NormalIndex:  make([]*Index, 0),
	}
}

// AddColumn 添加列.
func (t *Table) AddColumn(col *Column) {
	t.Columns = append(t.Columns, col)
	t.Name2Column[col.Name] = col
	if col.IsAutoIncrement() {
		t.AutoIncrementColumn = col
	}

	if !t.IsContainTime && col.IsTime() {
		t.IsContainTime = true
	}
}

// GetColumn 获取列.
func (t *Table) GetColumn(columnName string) *Column {
	col, ok := t.Name2Column[columnName]
	if ok {
		return col
	}

	return nil
}

// AddIndex 添加索引.
func (t *Table) AddIndex(ix *Index) {
	indexType := ix.Type()
	switch indexType {
	case PrimaryKeyIndex:
		t.PrimaryIndex = ix
	case UniqueKeyIndex:
		t.UniqueIndex = append(t.UniqueIndex, ix)
	case NormalKeyIndex:
		t.NormalIndex = append(t.NormalIndex, ix)
	default:
	}
}

// IsContainTimeType 是否包含time.Time字段.
func (t *Table) IsContainTimeType() bool {
	return t.IsContainTime
}

// LowerStartCamelObject 表名小写.
func (t *Table) LowerStartCamelObject() string {
	return LowerStartCamel(t.Name)
}

// UpperStartCamelObject 表名第一字母大写.
func (t *Table) UpperStartCamelObject() string {
	return UpperStartCamel(t.Name)
}

// GetPrimaryKeyUpperName 获取主键列.
func (t *Table) GetPrimaryKeyUpperName() string {
	return t.PrimaryIndex.GetColumnsUpperNameByDq()
}

// IsContainAutoIncrement 是否包含自增字段.
func (t *Table) IsContainAutoIncrement() bool {
	return t.AutoIncrementColumn != nil
}

// GetAutoKeyUpperName 获取自增列.
func (t *Table) GetAutoKeyUpperName() string {
	if t.AutoIncrementColumn != nil {
		return t.AutoIncrementColumn.GetUpperNameByDoubleQuotation()
	}

	return ""
}

// GetExpression 获取待插入字段列表.
func (t *Table) GetExpression() string {
	expressions := make([]string, 0, len(t.Columns))
	for _, col := range t.Columns {
		if col.IsAutoIncrement() {
			continue
		}

		expressions = append(expressions, "?")
	}

	return strings.Join(expressions, ", ")
}

// GetExpressionValues 获取待插入字段列表.
func (t *Table) GetExpressionValues() string {
	expressionValues := make([]string, 0, len(t.Columns))
	for _, col := range t.Columns {
		if col.IsAutoIncrement() {
			continue
		}

		expressionValues = append(expressionValues, "data."+col.GetUpperStartName())
	}

	return strings.Join(expressionValues, ", ")
}

// GetPrimaryExpressionValues 获取主键字段列表.
func (t *Table) GetPrimaryExpressionValues() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnsExpressionValues(), nil
}

// GetPrimaryKeyAndType 获取主键字段名和类型,逗号分隔.
func (t *Table) GetPrimaryKeyAndType() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnsNameAndType()
}

// GetPrimaryKeyAndMark 获取主键字段名和问号,and分隔.
func (t *Table) GetPrimaryKeyAndMark() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnsNameAndMark(), nil
}

// GetPrimaryKey 获取主键字段名,逗号分隔.
func (t *Table) GetPrimaryKey() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnsName(), nil
}
