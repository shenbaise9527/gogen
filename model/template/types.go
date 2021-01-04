package template

var Types = `
type (
	// {{.upperStartCamelObject}}Model model interface
	{{.upperStartCamelObject}}Model interface{
		{{.method}}
	}

	// default{{.upperStartCamelObject}}Model model object
	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}CachedDBConn{{else}}conn DBConn{{end}}
		table string
	}

	// {{.upperStartCamelObject}} {{.comment}}.
	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
