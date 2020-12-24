package {{.Package}}

var {{.Struct.Name}}Columns = []string{
    {{range $index, $field := .Struct.Fields}} "{{$field.Column}}", {{end}}
}

type {{.Struct.Name}} struct{
    {{range $index, $field := .Struct.Fields}} {{$field.Name}} {{$field.Type}} {{$field.Tag}} // {{$field.Comment}}
    {{end}}
}
