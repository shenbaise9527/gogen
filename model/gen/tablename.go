package gen

import (
	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genTableName(table *schemas.Table) (string, error) {
	output, err := template.With("tablename").
		Parse(template.TableName).
		Execute(table)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
