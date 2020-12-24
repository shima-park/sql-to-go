package sql_to_go

import "os"

var DefaultTemplateRendererConfig = TemplateRendererConfig{
	Package: func() string {
		pkg := os.Getenv("PACKAGE")
		if pkg != "" {
			return pkg
		}
		return "unknown"
	}(),
	TemplateFilepath: os.Getenv("TEMPLATE_FILEPATH"),
	TemplateText: `
package {{.Package}}

type {{.Struct.Name}} struct{
    {{range $index, $field := .Struct.Fields}} {{$field.Name}} {{$field.Type}} {{$field.Tag}} // {{$field.Comment}}
    {{end}}
}
    `,
}

type TemplateRendererConfig struct {
	Package          string
	TemplateFilepath string
	TemplateText     string
}
