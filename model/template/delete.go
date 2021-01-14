package template

var Delete = `
// Delete delete the record
func (m *default{{.UpperStartCamelObject}}Model) Delete(data *{{.UpperStartCamelObject}}) error {
	{{if .WithCachedAndUniqueIndex}}return m.delete(data){{else}}return m.DeleteBy{{.GetPrimaryIndexSuffixName}}({{.GetPrimaryExprValuesByPrefix "data."}}){{end}}
}

// Delete delete the record by the primary key
func (m *default{{.UpperStartCamelObject}}Model) DeleteBy{{.GetPrimaryIndexSuffixName}}({{.GetPrimaryKeyAndType}}) error {
	{{if .WithCached}}{{if .HasUniqueIndex}}data, err := m.FindOne({{.GetPrimaryKey}})
	if err != nil {
		return err
	}

	err = m.delete(data)
	{{else}}{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryKey}})
	_, err := m.Exec(func(conn DBConn) (int64, error) {
		db := conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}})

		return db.RowsAffected, db.Error
	}, {{.GetPrimaryIndexLowerName}}Key)
	{{end}}
	{{else}}
	err := m.conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryKey}}).Error
	{{end}}

	return err
}

{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
func (m *default{{$.UpperStartCamelObject}}Model) DeleteBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) error {
	{{if $.WithCached}}data, err := m.FindBy{{.GetSuffixName}}({{.GetColumnsName}})
	if err != nil {
		return err
	}

	err = m.delete(data)
	{{else}}
	err := m.conn.Delete({{$.UpperStartCamelObject}}{}, "{{.GetColumnsNameAndMark}}", {{.GetColumnsName}}).Error
	{{end}}

	return err
}
{{end}}

{{if .WithCachedAndUniqueIndex}}
func (m *default{{.UpperStartCamelObject}}Model) delete(data *{{.UpperStartCamelObject}}) error {
	{{.GetPrimaryIndexLowerName}}Key := fmt.Sprintf("{{.GetPrimaryIndexKeyFmt}}", cache{{.UpperStartCamelObject}}PKPrefix, {{.GetPrimaryExprValuesByPrefix "data."}})
	{{range .UniqueIndex}}
	{{.GetLowerName}}Key := fmt.Sprintf("{{.GetColumnKeyFmt}}", cache{{$.UpperStartCamelObject}}{{.GetSuffixName}}Prefix, {{.GetColumnsExprValuesByPrefix "data."}})
	{{end}}
	_, err := m.Exec(func(conn DBConn) (int64, error) {
		db := conn.Delete({{.UpperStartCamelObject}}{}, "{{.GetPrimaryKeyAndMark}}", {{.GetPrimaryExprValuesByPrefix "data."}})

		return db.RowsAffected, db.Error
	}, {{.GetPrimaryIndexLowerName}}Key, {{.GetUniqueIndexKey}})

	return err
}
{{end}}
`

var DeleteMethod = `
// Delete delete the record
Delete(data *{{.UpperStartCamelObject}}) error

// Delete delete the record by the primary key
DeleteBy{{.GetPrimaryIndexSuffixName}}({{.GetPrimaryKeyAndType}}) error
{{range .UniqueIndex}}
// DeleteBy{{.GetSuffixName}} delete the record by the unique key-{{.Name}}
DeleteBy{{.GetSuffixName}}({{.GetColumnsNameAndType}}) error
{{end}}
`
