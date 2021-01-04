package template

var FindOne = `
// FindOne query the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) FindOne({{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error) {
	var resp {{.UpperStartCamelObject}}
	err := m.conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(&resp).Error

	return &resp, err
}

{{range .UniqueIndex}}
// UpdateBy{{.GetSuffixName}} update the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) FindBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) (*{{$.UpperStartCamelObject}}, error) {
	var resp {{$.UpperStartCamelObject}}
	err := m.conn.Where("{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).First(&resp).Error

	return &resp, err
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
