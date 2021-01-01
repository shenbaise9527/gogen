package template

var Insert = `
// Insert insert the record
func (m *default{{.UpperStartCamelObject}}Model) Insert(data *{{.UpperStartCamelObject}}) (sql.Result, error) {
	{{if .IsContainAutoIncrement}}
	query := fmt.Sprintf("insert into %s values ({{.GetExpression}})", m.table, {{.LowerStartCamelObject}}RowsAutoWithPlaceHolder)
	{{else}}
	query := fmt.Sprintf("insert into %s values ({{.GetExpression}})", m.table, {{.LowerStartCamelObject}}Rows)
	{{end}}
	ret, err := m.conn.Exec(query, {{.GetExpressionValues}})
	
	return ret, err
}
`

var InsertMethod = `Insert(data *{{.UpperStartCamelObject}}) (sql.Result, error)`
