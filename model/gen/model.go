package gen

import (
	"strings"

	"github.com/shenbaise9527/gogen/model/schemas"
	"github.com/shenbaise9527/gogen/model/template"
)

func genModel(pkg string, table *schemas.Table) (string, error) {
	imports, err := genImport(table)
	if err != nil {
		return "", err
	}

	vars, err := genVars(table)
	if err != nil {
		return "", err
	}

	tableName, err := genTableName(table)
	if err != nil {
		return "", err
	}

	newStr, err := genNew(table)
	if err != nil {
		return "", err
	}

	insertStr, insertM, err := genInsert(table)
	if err != nil {
		return "", err
	}

	findStr, findM, err := genFind(table)
	if err != nil {
		return "", err
	}

	updateStr, updateM, err := genUpdate(table)
	if err != nil {
		return "", err
	}

	deleteStr, deleteM, err := genDelete(table)
	if err != nil {
		return "", err
	}

	var methods []string
	methods = append(methods, insertM, findM, updateM, deleteM)
	types, err := genTypes(table, strings.Join(methods, "\n"))
	if err != nil {
		return "", err
	}

	output, err := template.With("models").
		GoFmt(true).
		Parse(template.Model).
		Execute(map[string]interface{}{
			"pkg":       pkg,
			"imports":   imports,
			"vars":      vars,
			"tablename": tableName,
			"new":       newStr,
			"insert":    insertStr,
			"find":      findStr,
			"update":    updateStr,
			"delete":    deleteStr,
			"types":     types,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
