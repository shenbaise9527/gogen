package template

var Imports = `import (
"fmt"
"strings"
{{if .IsContainTimeType}}"time"{{end}}
)
`
