package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genTypes(table *schemas.Table, methods string) (string, error) {
	fieldsString, err := genFields(table)
	if err != nil {
		return "", err
	}

	output, err := template.With("types").
		Parse(template.Types).
		Execute(map[string]interface{}{
			"upperStartCamelObject": table.UpperStartCamelObject(),
			"method":                methods,
			"fields":                fieldsString,
			"comment":               table.TableComment,
		})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
