package template

var Insert = `
// Insert insert the record
func (m *default{{.UpperStartCamelObject}}Model) Insert(data *{{.UpperStartCamelObject}}) (sql.Result, error) {
	db := m.conn.Create(data)
	if db.Error != nil {
		return nil, db.Error
	}
	{{if .IsContainAutoIncrement}}
	res := newSqlResult(db.RowsAffected, data.{{.GetAutoUpperStartName}}, db.Error)
	{{else}}
	res := newSqlResult(db.RowsAffected, 0, db.Error)
	{{end}}
	
	return res, db.Error
}
`

var InsertMethod = `// Insert insert the record
Insert(data *{{.UpperStartCamelObject}}) (sql.Result, error)`
