package template

var Vars = `
var (
	{{.LowerStartCamelObject}}FieldNames = fieldNames(&{{.UpperStartCamelObject}}{})
	{{.LowerStartCamelObject}}Rows       = strings.Join({{.LowerStartCamelObject}}FieldNames, ",")
	{{.LowerStartCamelObject}}RowsPKWithPlaceHolder = strings.Join(removeField({{.LowerStartCamelObject}}FieldNames, {{.GetPrimaryKeyUpperName}}), "=?,") + "=?"
	{{if .IsContainAutoIncrement}}{{.LowerStartCamelObject}}RowsAutoWithPlaceHolder = strings.Join(removeField({{.LowerStartCamelObject}}FieldNames, {{.GetAutoKeyUpperName}}), "=?,") + "=?"{{end}}
	{{range .UniqueIndex}}{{$.LowerStartCamelObject}}Rows{{.GetSuffixName}}WithPlaceHolder = strings.Join(removeField({{$.LowerStartCamelObject}}FieldNames, {{.GetColumnsUpperNameByDq}}), "=?,") + "=?"{{end}}
)
`
