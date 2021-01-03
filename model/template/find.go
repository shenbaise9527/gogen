package template

var FindOne = `
// FindOne query the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) FindOne({{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error) {
	var resp {{.UpperStartCamelObject}}
	has, err := m.conn.Where("{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).First(&resp)
	if !has && err == nil {
		return nil, sql.ErrNoRows
	}

	return &resp, err
}
`

var FindOneMethod = `FindOne({{.GetPrimaryKeyAndType}}) (*{{.UpperStartCamelObject}}, error)`
