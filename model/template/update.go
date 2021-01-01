package template

var Update = `
// Update update the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) Update(data *{{.UpperStartCamelObject}}) error {
	query := fmt.Sprintf("update %s set %s where {{.GetPrimaryKeyAndMark}}", m.table, {{.LowerStartCamelObject}}RowsPKWithPlaceHolder)
	_, err := m.conn.Exec(query, {{.GetPrimaryExpressionValues}})

	return err
}
`

var UpdateMethod = `Update(data *{{.UpperStartCamelObject}}) error`
