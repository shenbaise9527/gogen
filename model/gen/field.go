package gen

import (
	"strings"

	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genFields(table *schemas.Table) (string, error) {
	var results []string
	for _, col := range table.Columns {
		tag, err := genTag(col)
		if err != nil {
			return "", err
		}

		colType, err := col.ConvertType()
		if err != nil {
			return "", err
		}

		comment := strings.ReplaceAll(col.ColumnComment, "\r", " ")
		comment = strings.ReplaceAll(col.ColumnComment, "\n", " ")
		output, err := template.With("field").
			Parse(template.Field).
			Execute(map[string]interface{}{
				"name":       col.GetUpperStartName(),
				"type":       colType,
				"tag":        tag,
				"hasComment": comment != "",
				"comment":    comment,
			})
		if err != nil {
			return "", nil
		}

		results = append(results, output.String())
	}

	return strings.Join(results, "\n"), nil
}
