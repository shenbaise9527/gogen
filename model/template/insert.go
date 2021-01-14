package template

var Insert = `
// Insert insert the record
func (m *default{{.UpperStartCamelObject}}Model) Insert(data *{{.UpperStartCamelObject}}) error {
	return m.conn.Create(data).Error
}
`

var InsertMethod = `// Insert insert the record
Insert(data *{{.UpperStartCamelObject}}) error`
