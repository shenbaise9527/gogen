package template

var FindOne = `
// FindOne query the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) FindOne(ctx context.Context, {{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error) {
	var err error
	{{if .WithTracing}}{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "FindOne")
	span.SetTag("key", {{.GetPrimaryIndexLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	{{if .WithCached}}{{if not .WithTracing}}{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}}){{end}}
	var resp {{.UpperStartCamelObject}}
	err = m.QueryRow(&resp, {{.GetPrimaryIndexLowerName}}Key, func(conn DBConn, v interface{}) error {
		return conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(v).Error
	})
	{{else}}var resp {{.UpperStartCamelObject}}
	err = m.conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(&resp).Error
	{{end}}

	return &resp, err
}

{{range .UniqueIndex}}
// FindBy{{.GetSuffixName}} query the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) FindBy{{.GetSuffixName}}(ctx context.Context, {{.GetColumnsNameAndType}}) (*{{$.UpperStartCamelObject}}, error) {
	var err error
	{{if $.WithTracing}}{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsName}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{$.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "FindBy{{.GetSuffixName}}")
	span.SetTag("key", {{.GetLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	var data {{$.UpperStartCamelObject}}
	{{if $.WithCached}}{{if not $.WithTracing}}{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsName}}){{end}}
	var primaryKey {{$.LowerStartCamelObject}}Primary
	var found bool
	err = m.cache.TakeWithExpire(&primaryKey, {{.GetLowerName}}Key, func(val interface{}, expire time.Duration) error {
		err := m.conn.Where("{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).First(&data).Error
		if err != nil {
			return err
		}

		found = true
		primaryKey = {{$.LowerStartCamelObject}}Primary{{$.GetPrimaryExprValuesByPrefixWrap "data."}}
		key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExprValuesByPrefix "data."}})
		return m.cache.SetCacheWithExpire(key, &data, expire+CacheSafeGapBetweenIndexAndPrimary)
	})
	if err != nil {
		return nil, err
	}

	if found {
		return &data, nil
	}

	key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExprValuesByPrefix "primaryKey."}})
	err = m.QueryRow(&data, key, func(conn DBConn, v interface{}) error {
		return conn.Where("{{$.GetPrimaryKeyAndMark}}", {{$.GetPrimaryExprValuesByPrefix "primaryKey."}}).First(v).Error
	})
	{{else}}err = m.conn.Where("{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).First(&data).Error{{end}}

	return &data, err
}
{{end}}
`

var FindOneMethod = `
// FindOne query the record by the primary key
FindOne(ctx context.Context, {{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error)

{{range .UniqueIndex}}
// FindBy{{.GetSuffixName}} query the record by the unique key-{{.Name}}
FindBy{{.GetSuffixName}}(ctx context.Context, {{.GetColumnsNameAndType}}) (*{{$.UpperStartCamelObject}}, error)
{{end}}
`
