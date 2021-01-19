package template

var New = `
// New{{.UpperStartCamelObject}}Model new model object
func New{{.UpperStartCamelObject}}Model(conn {{if .WithCached}}CachedDBConn{{else}}*DBConn{{end}}) {{.UpperStartCamelObject}}Model {
	return &default{{.UpperStartCamelObject}}Model{
		{{if .WithCached}}CachedDBConn{{else}}conn{{end}}:  conn,
		table: "{{.Name}}",
	}
}
`
