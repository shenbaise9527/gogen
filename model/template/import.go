package template

var Imports = `import (
"context"
"fmt"
"strings"
{{if .IsContainTimeType}}"time"{{end}}
{{if .WithTracing}}
"github.com/shenbaise9527/tracing"
"github.com/opentracing/opentracing-go/ext"
{{end}}
)
`
