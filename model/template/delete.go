package template

var Delete = `
// Delete delete the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) Delete({{.GetPrimaryKeyAndType}}) error {
	query := fmt.Sprintf("delete from %s where {{.GetPrimaryKeyAndMark}}", m.table)
	_, err := m.conn.Exec(query, {{.GetPrimaryKey}})

	return err
}
`

var DeleteMethod = `Delete({{.GetPrimaryKeyAndType}}) error`
