package sql_to_go

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types/parser_driver"
	"io"
	"io/ioutil"
)

var _ Parser = &SQLParser{}

type SQLParser struct {
	config SQLParserConfig
}

func NewSQLParser(config SQLParserConfig) *SQLParser {
	return &SQLParser{
		config: config,
	}
}

func NewDefaultSQLParser() *SQLParser {
	return NewSQLParser(DefaultSQLParserConfig)
}

func (sp *SQLParser) Parse(r io.Reader) ([]*Struct, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	p := parser.New()
	s, warns, err := p.Parse(string(b), "", "")
	if err != nil {
		return nil, err
	}

	for _, warn := range warns {
		if sp.config.IsStrictMode {
			return nil, warn
		}
	}

	var ss []*Struct
	for _, i := range s {
		v := sp.newColumnVisitor()
		i.Accept(v)
		ss = append(ss, v.newStruct)
	}

	return ss, nil
}

func (sp *SQLParser) newColumnVisitor() *columnVisitor {
	return &columnVisitor{
		newStruct: &Struct{},
		config:    sp.config,
	}
}

type columnVisitor struct {
	newStruct *Struct
	config    SQLParserConfig
}

func (v *columnVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if s, ok := in.(*ast.CreateTableStmt); ok {
		v.newStruct.Name = v.config.StructNameMap(s.Table.Name.O)

		for _, col := range s.Cols {
			opts := getColumnOptions(col.Options)
			field := Field{
				Column:  col.Name.Name.O,
				Name:    v.config.FieldNameMap(col.Name.Name.O),
				Type:    v.config.FieldTypeMap(col.Tp.Tp, opts),
				Tag:     v.config.FieldTagMap(col.Name.Name.O),
				Comment: opts.Comment,
			}
			v.newStruct.Fields = append(v.newStruct.Fields, field)
		}
	}
	return in, false
}

func (v *columnVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type ColumnOptions struct {
	Comment    string
	IsNullable bool
}

func getColumnOptions(options []*ast.ColumnOption) ColumnOptions {
	cos := ColumnOptions{}
	for _, opt := range options {
		switch opt.Tp {
		case ast.ColumnOptionComment:
			if v, ok := opt.Expr.(*driver.ValueExpr); ok {
				s, _ := v.ToString()
				cos.Comment = s
			}
		case ast.ColumnOptionNull:
			cos.IsNullable = true
		}
	}
	return cos
}
