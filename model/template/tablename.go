package template

var TableName = `
// TableName is {{.Name}}
func ({{.UpperStartCamelObject}}) TableName() string {
	return "{{.Name}}"
}
`
