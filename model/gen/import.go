package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genImport(table *schemas.Table) (string, error) {
	output, err := template.With("import").
		Parse(template.Imports).
		Execute(table)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
