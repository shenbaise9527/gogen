package template

var Delete = `
// Delete delete the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) Delete({{.GetPrimaryKeyAndType}}) error {
	err := m.conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).Error

	return err
}

{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) DeleteBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) error {
	err := m.conn.Delete({{$.UpperStartCamelObject}}{}, "{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).Error

	return err
}
{{end}}
`

var DeleteMethod = `
// Delete delete the record by the primary key
Delete({{.GetPrimaryKeyAndType}}) error
{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
DeleteBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) error
{{end}}
`
