package template

import "github.com/deadcheat/awsset"

// Assets struct for assign values to template
type Assets struct {
	PackageName string
	VarName     string
	DirMap      map[string][]string
	FileMap     map[string]*awsset.File
	Paths       []string
}

// AssetFileTemplate template for asset file
var AssetFileTemplate = `package {{.PackageName}}

import(
	"time"

	"github.com/deadcheat/awsset"
)

// {{ .VarName }} a generated file system
var {{ .VarName }} = awsset.NewFS(
	map[string][]string{
		{{- range $k, $v := .DirMap}}
		"{{ $k }}": []string{
			{{ range $s := $v }}"{{ $s }}", {{ end }}
		},
		{{- end }}
	},
	map[string]*File {
		{{- range $k, $v := .FileMap}}
		"{{$k}}": NewFile("{{$k}}", _{{ sha1 $k }}, {{ printf "%#v" $v.FileMode }}, time.Unix({{ $v.ModifiedAt.Unix }}, {{ $v.ModifiedAt.UnixNano }})),
		{{- end }}
	},
)

// binary data
var (
	{{ $m := .FileMap}}
	{{- range $p := .Paths }}
	_{{ sha1 $p}} = {{ with (index $m $p) }}{{if not .Data }}nil{{ else }}{{ printf "%#v" .Data }}{{ end }}{{ end }}
	{{- end }}
)
`
