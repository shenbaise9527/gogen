package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genTag(col *schemas.Column) (string, error) {
	output, err := template.With("tag").
		Parse(template.Tag).
		Execute(col)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
