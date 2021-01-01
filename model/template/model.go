package template

var Model = `package {{.pkg}}
{{.imports}}
{{.vars}}
{{.types}}
{{.tablename}}
{{.new}}
{{.insert}}
{{.find}}
{{.update}}
{{.delete}}
`
