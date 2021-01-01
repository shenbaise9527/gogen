package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genNew(table *schemas.Table) (string, error) {
	output, err := template.With("new").
		Parse(template.New).
		Execute(table)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
