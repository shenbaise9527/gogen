package template

var FindOne = `
// FindOne query the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) FindOne({{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error) {
	{{if .WithCached}}{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}})
	var resp {{.UpperStartCamelObject}}
	err := m.QueryRow(&resp, {{.GetPrimaryIndexLowerName}}Key, func(conn DBConn, v interface{}) error {
		return conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(v).Error
	})
	{{else}}
	var resp {{.UpperStartCamelObject}}
	err := m.conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(&resp).Error
	{{end}}
	if err != nil {
		return nil, err
	}

	return &resp, err
}

{{range .UniqueIndex}}
// FindBy{{.GetSuffixName}} query the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) FindBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) (*{{$.UpperStartCamelObject}}, error) {
	{{if $.WithCached}}{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsName}})
	var data {{$.UpperStartCamelObject}}
	var primaryKey {{$.LowerStartCamelObject}}Primary
	var found bool
	err := m.cache.TakeWithExpire(&primaryKey, {{.GetLowerName}}Key, func(val interface{}, expire time.Duration) error {
		err := m.conn.Where("{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).First(&data).Error
		if err != nil {
			return err
		}

		found = true
		primaryKey = {{$.LowerStartCamelObject}}Primary{{$.GetPrimaryExprValuesByPrefixWrap "data."}}
		key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExpressionValues}})
		return m.cache.SetCacheWithExpire(key, &data, expire+CacheSafeGapBetweenIndexAndPrimary)
	})
	if err != nil {
		return nil, err
	}

	if found {
		return &data, nil
	}

	key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExpressionValuesByPrefix "primaryKey."}})
	err = m.QueryRow(&data, key, func(conn DBConn, v interface{}) error {
		return conn.Where("{{$.GetPrimaryKeyAndMark}}", {{$.GetPrimaryExpressionValuesByPrefix "primaryKey."}}).First(v).Error
	})
	{{else}}
	err := m.conn.Where("{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).First(&data).Error
	{{end}}

	return &data, err
}
{{end}}
`

var FindOneMethod = `
// FindOne query the record by the primary key
FindOne({{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error)

{{range .UniqueIndex}}
// FindBy{{.GetSuffixName}} query the record by the unique key-{{.Name}}
FindBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) (*{{$.UpperStartCamelObject}}, error)
{{end}}
`
