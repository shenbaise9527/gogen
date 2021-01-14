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
	WithCached          bool               // 是否缓存.
}

// NewTable 新建表对象.
func NewTable(name, comment string, withCached bool) *Table {
	return &Table{
		Name:         name,
		TableComment: comment,
		Columns:      make([]*Column, 0),
		Name2Column:  make(map[string]*Column),
		UniqueIndex:  make([]*Index, 0),
		NormalIndex:  make([]*Index, 0),
		WithCached:   withCached,
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

// GetPrimaryKeyName 获取主键列,字段名带双引号.
func (t *Table) GetPrimaryKeyName() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnsNameByDq(), nil
}

// IsContainAutoIncrement 是否包含自增字段.
func (t *Table) IsContainAutoIncrement() bool {
	return t.AutoIncrementColumn != nil
}

// GetAutoKeyName 获取自增列.
func (t *Table) GetAutoKeyName() string {
	if t.AutoIncrementColumn != nil {
		return t.AutoIncrementColumn.GetNameByDoubleQuotation()
	}

	return ""
}

// GetPrimaryAndAutoKeyName 获取主键列和自增列.
func (t *Table) GetPrimaryAndAutoKeyName() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	pk := t.PrimaryIndex.GetColumnNameSliceByDq()
	if t.AutoIncrementColumn != nil {
		ak := t.AutoIncrementColumn.GetNameByDoubleQuotation()
		var flag bool
		for i := range pk {
			if pk[i] == ak {
				flag = true
				break
			}
		}

		if !flag {
			pk = append(pk, ak)
		}
	}

	return strings.Join(pk, ","), nil
}

// GetPKUpdateExpressionValues 获取主键字段列表.
func (t *Table) GetPKUpdateExpressionValues() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	pk := t.PrimaryIndex.GetColumns()
	expressionValues := make([]string, 0, len(t.Columns))
	for _, col := range t.Columns {
		if _, ok := pk[col.Name]; ok {
			continue
		}

		if col.IsAutoIncrement() {
			continue
		}

		expressionValues = append(expressionValues, "data."+col.GetUpperStartName())
	}

	return strings.Join(expressionValues, ", ") + ", " + t.PrimaryIndex.GetColumnsExprValuesByPrefix("data."), nil
}

// GetUKUpdateExpressionValues 根据唯一索引更新时的字段列表.
func (t *Table) GetUKUpdateExpressionValues(ukname string) (string, error) {
	for _, ix := range t.UniqueIndex {
		if ukname != ix.Name {
			continue
		}

		uk := ix.GetColumns()
		expressionValues := make([]string, 0, len(t.Columns))
		for _, col := range t.Columns {
			if _, ok := uk[col.Name]; ok {
				continue
			}

			if col.IsAutoIncrement() || col.IsPrimaryKey() {
				continue
			}

			expressionValues = append(expressionValues, "data."+col.GetUpperStartName())
		}

		return strings.Join(expressionValues, ", ") + ", " + ix.GetColumnsExprValuesByPrefix("data."), nil
	}

	return "", fmt.Errorf("the table(%s) does not have the unique index(%s)", t.Name, ukname)
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

// GetAutoUpperStartName 获取自增字段.
func (t *Table) GetAutoUpperStartName() (string, error) {
	if t.AutoIncrementColumn == nil {
		return "", fmt.Errorf("%s has not autoincrement key.", t.Name)
	}

	return t.AutoIncrementColumn.GetUpperStartName(), nil
}

// GetPrimaryCacheKey 获取主键缓存用的列名.
func (t *Table) GetPrimaryCacheKey() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnCacheName(), nil
}

// GetPrimaryExprValuesByPrefix 获取主键字段列表.
func (t *Table) GetPrimaryExprValuesByPrefix(prefix string) (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	pk := t.PrimaryIndex.GetColumnsExprValuesByPrefix(prefix)

	return pk, nil
}

// GetPrimaryExprValuesByPrefixWrap 获取主键字段列表.
func (t *Table) GetPrimaryExprValuesByPrefixWrap(prefix string) (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	pk := "{" + t.PrimaryIndex.GetColumnsExprValuesByPrefix(prefix) + "}"

	return pk, nil
}

// GetPrimaryIndexLowerName 获取主键索引名字.
func (t *Table) GetPrimaryIndexLowerName() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetLowerName(), nil
}

// GetPrimaryIndexSuffixName 获取主键索引名字.
func (t *Table) GetPrimaryIndexSuffixName() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetSuffixName(), nil
}

// GetPrimaryIndexKeyFmt 获取主键的序列,%s#%v#%v...
func (t *Table) GetPrimaryIndexKeyFmt() (string, error) {
	if t.PrimaryIndex == nil {
		return "", fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.GetColumnKeyFmt(), nil
}

// GetPrimaryIndexColumns 获取主键索引名字.
func (t *Table) GetPrimaryIndexColumns() ([]*Column, error) {
	if t.PrimaryIndex == nil {
		return nil, fmt.Errorf("%s has not primarykey.", t.Name)
	}

	return t.PrimaryIndex.Columns, nil
}

// HasUniqueIndex 是否存在唯一索引.
func (t *Table) HasUniqueIndex() bool {
	return len(t.UniqueIndex) != 0
}

// GetUniqueIndexKey 获取唯一索引key的列表.
func (t *Table) GetUniqueIndexKey() string {
	keys := make([]string, 0, len(t.UniqueIndex))
	for i := range t.UniqueIndex {
		keys = append(keys, t.UniqueIndex[i].GetLowerName()+"Key")
	}

	return strings.Join(keys, ", ")
}

// WithCachedAndUniqueIndex 是否有缓存和存在唯一索引.
func (t *Table) WithCachedAndUniqueIndex() bool {
	return t.WithCached && len(t.UniqueIndex) != 0
}
