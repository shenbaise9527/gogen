package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genFind(table *schemas.Table) (string, string, error) {
	output, err := template.With("find").
		Parse(template.FindOne).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	methodOut, err := template.With("findMethod").
		Parse(template.FindOneMethod).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	return output.String(), methodOut.String(), nil
}
