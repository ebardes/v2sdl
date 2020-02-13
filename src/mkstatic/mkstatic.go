// +build none
package main

import (
	"os"
	"sort"
	"text/template"

	"github.com/rs/zerolog/log"
)

func main() {
	t := template.New("")
	t, err := t.Parse(temp)
	if err != nil {
		log.Error().Err(err).Msg("Parse")
		return
	}
	files := os.Args[1:]
	sort.Strings(files)

	fileinfo := make([]os.FileInfo, len(files))
	for i, fn := range files {
		fileinfo[i], err = os.Stat(fn)
		if err != nil {
			panic(err)
		}
	}

	err = t.Execute(os.Stdout, fileinfo)
	if err != nil {
		panic(err)
	}
}

var temp = `

_staticlen:
.quad {{ len . }}

{{- range $i,$n := . }}
name_{{$i}}:
.asciz "{{$n.Name}}"
{{- end }}

{{- range $i,$n := . }}
data_{{$i}}_start:
.incbin "{{$n.Name}}"
data_{{$i}}_end:
{{- end }}

_staticnames:
{{- range $i,$n := . }}
.quad name_{{$i}}
.quad data_{{$i}}_start
.quad {{$n.Size}}
.quad {{$n.ModTime.UTC.Unix}}
{{- end }}

`
