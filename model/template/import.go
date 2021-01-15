package template

var Imports = `import (
"context"
"fmt"
"strings"
{{if .IsContainTimeType}}"time"{{end}}
)
`
