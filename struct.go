package sql_to_go

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Column  string
	Name    string
	Type    string
	Tag     string
	Comment string
}
