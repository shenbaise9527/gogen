package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genVars(table *schemas.Table) (string, error) {
	output, err := template.With("vars").
		Parse(template.Vars).
		Execute(table)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
