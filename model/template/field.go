package template

var Field = "{{.GetUpperStartName}} {{.ConvertType}} `json:\"{{.GetLowerName}}\" gorm:\"{{.GetTags}}\"` {{if .HasComment}}// {{.GetComment}}{{end}}"
