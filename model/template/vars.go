package template

var Vars = `
var (
	{{.LowerStartCamelObject}}FieldNames = fieldNames(&{{.UpperStartCamelObject}}{})

	{{if .IsContainAutoIncrement}}{{.LowerStartCamelObject}}RowsNoPA = strings.Join(removeField({{.LowerStartCamelObject}}FieldNames, {{.GetPrimaryAndAutoKeyName}}), "=?,") + "=?"
	{{else}}{{.LowerStartCamelObject}}RowsNoPK = strings.Join(removeField({{.LowerStartCamelObject}}FieldNames, {{.GetPrimaryKeyName}}), "=?,") + "=?"{{end}}
	{{range .UniqueIndex}}
	{{$.LowerStartCamelObject}}Rows{{.GetSuffixName}}NoPA = strings.Join(removeField({{$.LowerStartCamelObject}}FieldNames, {{.GetColumnsNameByDq}}, {{$.GetPrimaryAndAutoKeyName}}), "=?,") + "=?"
	{{end}}
)
`
