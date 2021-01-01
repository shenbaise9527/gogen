package template

var Imports = `import (
"database/sql"
"fmt"
"strings"
{{if .IsContainTimeType}}"time"{{end}}

"gorm.io/gorm"
)
`
