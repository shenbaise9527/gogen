package template

var Insert = `
// Insert insert the record
func (m *default{{.UpperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.UpperStartCamelObject}}) error {
	{{if .WithTracing}}var err error
	span := tracing.ChildOfSpanFromContext(ctx, "{{.LowerStartCamelObject}}model")
	defer span.Finish()
	ext.DBStatement.Set(span, "Insert")
	defer func() {
		if err != nil {
			ext.Error.Set(span, true)
			span.LogKV("error", err.Error())
		}
	}()

	err = m.conn.Create(data).Error
	{{else}}err := m.conn.Create(data).Error{{end}}

	return err
}
`

var InsertMethod = `// Insert insert the record
Insert(ctx context.Context, data *{{.UpperStartCamelObject}}) error`
