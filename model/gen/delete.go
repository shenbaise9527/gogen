package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genDelete(table *schemas.Table) (string, string, error) {
	output, err := template.With("delete").
		Parse(template.Delete).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	methodOut, err := template.With("deleteMethod").
		Parse(template.DeleteMethod).
		Execute(table)
	if err != nil {
		return "", "", err
	}

	return output.String(), methodOut.String(), nil
}
