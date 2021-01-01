package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genInsert(table *schemas.Table) (string, string, error) {
	output, err := template.With("insert").
		Parse(template.Insert).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	methodOut, err := template.With("insertMethod").
		Parse(template.InsertMethod).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	return output.String(), methodOut.String(), nil
}
