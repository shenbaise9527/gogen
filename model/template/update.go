package template

var Update = `
// Update update the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) Update(ctx context.Context, data *{{.UpperStartCamelObject}}) error {
	{{if .WithTracing}}var err error
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryExprValuesByPrefix "data."}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "Update")
	span.SetTag("key", {{.GetPrimaryIndexLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	{{if .WithCached}}{{if not .WithTracing}}var err error
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryExprValuesByPrefix "data."}}){{end}}
	_, err = m.Exec(func(conn *DBConn) (int64, error) {
		{{if .IsContainAutoIncrement}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPA){{else}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPK){{end}}
		db := conn.Exec(query, {{.GetPKUpdateExpressionValues}})

		return db.RowsAffected, db.Error
	}, {{.GetPrimaryIndexLowerName}}Key)
	{{else}}{{if .WithTracing}}err = m.conn.DoWithAcceptable({{else}}err := m.conn.DoWithAcceptable({{end}}
		func() error {
			{{if .IsContainAutoIncrement}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPA){{else}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPK){{end}}
			err := m.conn.Exec(query, {{.GetPKUpdateExpressionValues}}).Error

			return err
		}, m.conn.Acceptable)
	{{end}}

	return err
}

{{range .UniqueIndex}}
// UpdateBy{{.GetSuffixName}} update the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) UpdateBy{{.GetSuffixName}}(ctx context.Context, data *{{$.UpperStartCamelObject}}) error {
	{{if $.WithTracing}}var err error
	{{$.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExprValuesByPrefix "data."}})
	span := tracing.ChildOfSpanFromContext(ctx, "{{$.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "UpdateBy{{.GetSuffixName}}")
	span.SetTag("key", {{$.GetPrimaryIndexLowerName}}Key)
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()
	{{end}}
	{{if $.WithCached}}{{if not $.WithTracing}}var err error
	{{$.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{$.GetPrimaryIndexKeyFmt}}", cache{{$.UpperStartCamelObject}}PKPrefix, {{$.GetPrimaryExprValuesByPrefix "data."}}){{end}}
	_, err = m.Exec(func(conn *DBConn) (int64, error) {
		query := fmt.Sprintf("update %s set %s where {{.GetColumnsNameAndMark}}", m.table, {{$.LowerStartCamelObject}}Rows{{.GetSuffixName}}NoPA)
		db := conn.Exec(query, {{$.GetUKUpdateExpressionValues .Name}})

		return db.RowsAffected, db.Error
	}, {{$.GetPrimaryIndexLowerName}}Key)
	{{else}}{{if $.WithTracing}}err = m.conn.DoWithAcceptable({{else}}err := m.conn.DoWithAcceptable({{end}}
		func() error {
			query := fmt.Sprintf("update %s set %s where {{.GetColumnsNameAndMark}}", m.table, {{$.LowerStartCamelObject}}Rows{{.GetSuffixName}}NoPA)
			err := m.conn.Exec(query, {{$.GetUKUpdateExpressionValues .Name}}).Error

			return err
		}, m.conn.Acceptable)
	{{end}}

	return err
}
{{end}}
`

var UpdateMethod = `
// Update update the record by the primary key
Update(ctx context.Context, data *{{.UpperStartCamelObject}}) error
{{range .UniqueIndex}}
// UpdateBy{{.GetSuffixName}} update the record by the unique key-{{.Name}}
UpdateBy{{.GetSuffixName}}(ctx context.Context, data *{{$.UpperStartCamelObject}}) error
{{end}}
`
