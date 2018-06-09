package template

import "github.com/deadcheat/goblet"

// Assets struct for assign values to template
type Assets struct {
	PackageName string
	VarName     string
	DirMap      map[string][]string
	FileMap     map[string]*goblet.File
	Paths       []string
}

// AssetFileTemplate template for asset file
var AssetFileTemplate = `package {{.PackageName}}

import(
	"time"

	"github.com/deadcheat/goblet"
)
{{ $FileMap := .FileMap}}{{ $DirMap := .DirMap}}{{ $VarName := .VarName }}
// {{ $VarName }} a generated file system
var {{ $VarName }} = goblet.NewFS(
	map[string][]string{
		{{- range $p := .Paths }}{{ with (index $DirMap $p)}}
		"{{ $p }}": []string{
			{{ range $s := . }}"{{ $s }}", {{ end }}
		},{{ end }}
		{{- end }}
	},
	map[string]*goblet.File {
		{{- range $p := .Paths }}{{ with (index $FileMap $p)}}
		"{{$p}}": goblet.NewFile("{{$p}}", {{if not .Data }}nil{{ else }}[]byte(_{{ $VarName }}{{ sha1 $p }}){{ end }}, {{ printf "%#v" .FileMode }}, time.Unix({{ .ModifiedAt.Unix }}, {{ .ModifiedAt.UnixNano }})),{{ end }}
		{{- end }}
	},
)

// binary data
var (
	{{- range $p := .Paths }}{{ with (index $FileMap $p) }}
	{{if .Data }}_{{ $VarName }}{{ sha1 $p}} = {{ printData .Data }}{{ end }}{{ end }}
	{{- end }}
)
`
