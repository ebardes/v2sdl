// +build none
package main

import (
	"os"
	"sort"
	"text/template"

	"github.com/rs/zerolog/log"
)

type FullInfo struct {
	os.FileInfo
	Path string
}

func main() {
	t := template.New("")
	t, err := t.Parse(temp)
	if err != nil {
		log.Error().Err(err).Msg("Parse")
		return
	}
	files := os.Args[1:]
	sort.Strings(files)

	fileinfo := make([]FullInfo, len(files))
	for i, fn := range files {
		fi, err := os.Stat(fn)
		if err != nil {
			panic(err)
		}
		fileinfo[i].FileInfo = fi
		fileinfo[i].Path = fn
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
.asciz "{{$n.Path}}"
{{- end }}

{{- range $i,$n := . }}
data_{{$i}}_start:
.incbin "{{$n.Path}}"
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
