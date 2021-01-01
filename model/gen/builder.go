package gen

import "github.com/shenbaise9527/gogen/model/template"

func genBuilder(pkg string) (string, error) {
	output, err := template.With("builder").
		GoFmt(true).
		Parse(template.Builder).
		Execute(map[string]interface{}{
			"pkg": pkg,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
