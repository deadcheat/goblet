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
{{ $FileMap := .FileMap}}{{ $DirMap := .DirMap}}{{ $VarName := .VarName }}
// {{ $VarName }} a generated file system
var {{ $VarName }} = awsset.NewFS(
	map[string][]string{
		{{- range $p := .Paths }}{{ with (index $DirMap $p)}}
		"{{ $p }}": []string{
			{{ range $s := . }}"{{ $s }}", {{ end }}
		},{{ end }}
		{{- end }}
	},
	map[string]*awsset.File {
		{{- range $p := .Paths }}{{ with (index $FileMap $p)}}
		"{{$p}}": awsset.NewFile("{{$p}}", {{if not .Data }}nil{{ else }}[]byte(_{{ $VarName }}{{ sha1 $p }}){{ end }}, {{ printf "%#v" .FileMode }}, time.Unix({{ .ModifiedAt.Unix }}, {{ .ModifiedAt.UnixNano }})),{{ end }}
		{{- end }}
	},
)

// binary data
var (
	{{- range $p := .Paths }}{{ with (index $FileMap $p) }}
	{{if .Data }}_{{ $VarName }}{{ sha1 $p}} = {{ printf "%s" .Data |  printf "%#v" | safeHTML }}{{ end }}{{ end }}
	{{- end }}
)
`
