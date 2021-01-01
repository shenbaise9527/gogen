package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genUpdate(table *schemas.Table) (string, string, error) {
	output, err := template.With("update").
		Parse(template.Update).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	methodOut, err := template.With("updateMethod").
		Parse(template.UpdateMethod).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	return output.String(), methodOut.String(), nil
}
