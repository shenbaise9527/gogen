package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genTypes(table *schemas.Table, methods string, withCache bool) (string, error) {
	fieldsString, err := genFields(table)
	if err != nil {
		return "", err
	}

	primaryFields, err := genPrimaryFields(table)
	if err != nil {
		return "", err
	}

	output, err := template.With("types").
		Parse(template.Types).
		Execute(map[string]interface{}{
			"upperStartCamelObject": table.UpperStartCamelObject(),
			"method":                methods,
			"fields":                fieldsString,
			"withCache":             withCache,
			"comment":               table.TableComment,
			"lowerStartCamelObject": table.LowerStartCamelObject(),
			"primaryfields":         primaryFields,
			"hasUniqueIndex":        table.HasUniqueIndex(),
		})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
