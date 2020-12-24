package sql_to_go

import "io"

type Parser interface {
	Parse(io.Reader) ([]*Struct, error)
}
