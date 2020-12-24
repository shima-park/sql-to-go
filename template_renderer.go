package sql_to_go

import (
	"io"
	"text/template"
)

type TemplateRenderer struct {
	config TemplateRendererConfig
}

func NewTemplateRenderer(config TemplateRendererConfig) *TemplateRenderer {
	return &TemplateRenderer{
		config: config,
	}
}

func NewDefaultTemplateRenderer() *TemplateRenderer {
	return NewTemplateRenderer(DefaultTemplateRendererConfig)
}

func (t TemplateRenderer) Render(w io.Writer, s *Struct) error {
	tpl := template.New("")

	var err error
	if t.config.TemplateFilepath != "" {
		tpl, err = tpl.ParseFiles(t.config.TemplateFilepath)
		if err != nil {
			return err
		}
		tpl = tpl.Lookup(t.config.TemplateFilepath)
	} else if t.config.TemplateText != "" {
		tpl, err = tpl.Parse(t.config.TemplateText)
	}
	if err != nil {
		return err
	}

	return tpl.Execute(w, map[string]interface{}{
		"Package": t.config.Package,
		"Struct":  s,
	})
}
