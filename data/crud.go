package data

type Crud struct {
	Root
	Name       string
	Entities   []RequestEntity
	GoType     func(s string) TypeID
	Filterable []string
}

type Root struct {
	Module    string
	Names     []string
	OpenApiUI string
}

type RequestEntity struct {
	TableCatalog    string
	TableSchema     string
	TableName       string
	ColumnName      string
	OrdinalPosition string
	ColumnDefault   string
	IsNullable      string
	DataType        string
}

type TypeID struct {
	Name     string
	Pkg      string
	PkgShort string
}
