package template

var Delete = `
// Delete delete the record
func (m *default{{.UpperStartCamelObject}}Model) Delete(ctx context.Context, data *{{.UpperStartCamelObject}}) error {
	{{if .WithCachedAndUniqueIndex}}{{if .WithTracing}}var err error
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryExprValuesByPrefix "data."}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "Delete")
	span.SetTag("key", {{.GetPrimaryIndexLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	err = m.delete(data){{else}}err := m.delete(data){{end}}
	{{else}}err := m.DeleteBy{{.GetPrimaryIndexSuffixName}}(ctx, {{.GetPrimaryExprValuesByPrefix "data."}}){{end}}

	return err
}

// Delete delete the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) DeleteBy{{.GetPrimaryIndexSuffixName}}(ctx context.Context, {{.GetPrimaryKeyAndType}}) error {
	{{if .WithTracing}}var err error
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "DeleteBy{{.GetPrimaryIndexSuffixName}}")
	span.SetTag("key", {{.GetPrimaryIndexLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	{{if .WithCached}}{{if not .WithTracing}}var err error{{end}}
	{{if .HasUniqueIndex}}data, err := m.FindOne(ctx, {{.GetPrimaryKey}})
	if err != nil {
		return err
	}

	err = m.delete(data)
	{{else}}{{if not .WithTracing}}{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}}){{end}}
	_, err = m.Exec(func(conn *DBConn) (int64, error) {
		db := conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}})

		return db.RowsAffected, db.Error
	}, {{.GetPrimaryIndexLowerName}}Key)
	{{end}}
	{{else}}{{if .WithTracing}}err = m.conn.DoWithAcceptable({{else}}err := m.conn.DoWithAcceptable({{end}}
		func() error {
			err := m.conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).Error

			return err
		}, m.conn.Acceptable)
	{{end}}

	return err
}

{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) DeleteBy{{.GetSuffixName}}(ctx context.Context, {{.GetColumnsNameAndType}}) error {
	{{if $.WithTracing}}var err error
	{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsName}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{$.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "DeleteBy{{.GetSuffixName}}")
	span.SetTag("key", {{.GetLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	{{if $.WithCached}}data, err := m.FindBy{{.GetSuffixName}}(ctx, {{.GetColumnsName}})
	if err != nil {
		return err
	}

	err = m.delete(data)
	{{else}}{{if $.WithTracing}}err = m.conn.DoWithAcceptable({{else}}err := m.conn.DoWithAcceptable({{end}}
		func() error {
			err := m.conn.Delete({{$.UpperStartCamelObject}}{}, "{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).Error

			return err
		}, m.conn.Acceptable)
	{{end}}

	return err
}
{{end}}

{{if .WithCachedAndUniqueIndex}}
func (m *default{{.UpperStartCamelObject}}Model) delete(data *{{.UpperStartCamelObject}}) error {
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryExprValuesByPrefix "data."}})
	{{range .UniqueIndex}}
	{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsExprValuesByPrefix "data."}})
	{{end}}
	_, err := m.Exec(func(conn *DBConn) (int64, error) {
		db := conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryExprValuesByPrefix "data."}})

		return db.RowsAffected, db.Error
	}, {{.GetPrimaryIndexLowerName}}Key, {{.GetUniqueIndexKey}})

	return err
}
{{end}}
`

var DeleteMethod = `
// Delete delete the record
Delete(ctx context.Context, data *{{.UpperStartCamelObject}}) error

// Delete delete the record by the primary key
DeleteBy{{.GetPrimaryIndexSuffixName}}(ctx context.Context, {{.GetPrimaryKeyAndType}}) error
{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
DeleteBy{{.GetSuffixName}}(ctx context.Context, {{.GetColumnsNameAndType}}) error
{{end}}
`
