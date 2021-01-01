package template

var Types = `
type (
	// {{.upperStartCamelObject}}Model model interface
	{{.upperStartCamelObject}}Model interface{
		{{.method}}
	}

	// default{{.upperStartCamelObject}}Model model object
	default{{.upperStartCamelObject}}Model struct {
		conn *gorm.DB
		table string
	}

	// {{.upperStartCamelObject}} {{.comment}}.
	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
