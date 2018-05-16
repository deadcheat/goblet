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
{{ $FileMap := .FileMap}}
{{ $DirMap := .DirMap}}
// {{ .VarName }} a generated file system
var {{ .VarName }} = awsset.NewFS(
	map[string][]string{
		{{- range $p := .Paths }}{{ with (index $DirMap $p)}}
		"{{ $p }}": []string{
			{{ range $s := . }}"{{ $s }}", {{ end }}
		},{{ end }}
		{{- end }}
	},
	map[string]*File {
		{{- range $p := .Paths }}{{ with (index $FileMap $p)}}
		"{{$p}}": awsset.NewFile("{{$p}}", _{{ .VarName }}{{ sha1 $p }}, {{ printf "%#v" .FileMode }}, time.Unix({{ .ModifiedAt.Unix }}, {{ .ModifiedAt.UnixNano }})),{{ end }}
		{{- end }}
	},
)

// binary data
var (
	{{- range $p := .Paths }}
	_{{ .VarName }}{{ sha1 $p}} = {{ with (index $FileMap $p) }}{{if not .Data }}nil{{ else }}{{ printf "%#v" .Data }}{{ end }}{{ end }}
	{{- end }}
)
`
