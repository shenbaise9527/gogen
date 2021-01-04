package template

var Update = `
// Update update the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) Update(data *{{.UpperStartCamelObject}}) error {
	{{if .IsContainAutoIncrement}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPA)
	{{else}}query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsNoPK)
	{{end}}
	err := m.conn.Exec(query, {{.GetPKUpdateExpressionValues}}).Error

	return err
}

{{range .UniqueIndex}}
// UpdateBy{{.GetSuffixName}} update the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) UpdateBy{{.GetSuffixName}}(data *{{$.UpperStartCamelObject}}) error {
	query := fmt.Sprintf("update %s set %s where {{.GetColumnsNameAndMark}}", m.table, {{$.LowerStartCamelObject}}Rows{{.GetSuffixName}}NoPA)
	err := m.conn.Exec(query, {{$.GetUKUpdateExpressionValues .Name}}).Error

	return err
}
{{end}}
`

var UpdateMethod = `
// Update update the record by the primary key
Update(data *{{.UpperStartCamelObject}}) error
{{range .UniqueIndex}}
// UpdateBy{{.GetSuffixName}} update the record by the unique key-{{.Name}}
UpdateBy{{.GetSuffixName}}(data *{{$.UpperStartCamelObject}}) error
{{end}}
`
