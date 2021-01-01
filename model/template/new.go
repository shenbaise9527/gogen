package template

var New = `
// New{{.UpperStartCamelObject}}Model new model object
func New{{.UpperStartCamelObject}}Model(conn *gorm.DB) {{.UpperStartCamelObject}}Model {
	return &default{{.UpperStartCamelObject}}Model{
		conn: conn,
		table: "{{.Name}}",
	}
}
`
