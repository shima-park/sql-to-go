package sql_to_go

import "io"

type Renderer interface {
	Render(w io.Writer, s *Struct) error
}
