package sql_to_go

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
)

type SQL2Go struct {
	parser   Parser
	renderer Renderer
}

func NewSQL2Go(parser Parser, renderer Renderer) *SQL2Go {
	return &SQL2Go{
		parser:   parser,
		renderer: renderer,
	}
}

func NewDefaultSQL2Go() *SQL2Go {
	return NewSQL2Go(NewDefaultSQLParser(), NewDefaultTemplateRenderer())
}

func (sg *SQL2Go) Convert(r io.Reader, w io.Writer) error {
	ss, err := sg.parser.Parse(r)
	if err != nil {
		return err
	}

	for i, s := range ss {
		buf := bytes.NewBuffer(nil)
		err = sg.renderer.Render(buf, s)
		if err != nil {
			return err
		}

		b, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Println(string(buf.Bytes()))
			return err
		}

		fmt.Fprint(w, string(b))
		if len(ss)-1 != i {
			fmt.Fprintln(w)
		}
	}
	return nil
}
