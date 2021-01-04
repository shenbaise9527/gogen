package gen

import "github.com/shenbaise9527/gogen/model/template"

func genDBConn(pkg string) (string, error) {
	output, err := template.With("dbconn").
		GoFmt(true).
		Parse(template.DBConn).
		Execute(map[string]interface{}{
			"pkg": pkg,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
