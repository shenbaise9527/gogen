package gen

import (
	"strings"

	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genFields(table *schemas.Table) (string, error) {
	var results []string
	for _, col := range table.Columns {
		output, err := template.With("field").
			Parse(template.Field).
			Execute(col)
		if err != nil {
			return "", nil
		}

		results = append(results, output.String())
	}

	return strings.Join(results, "\n"), nil
}
