package pg_info

type Column struct {
	tableName       struct{} `pg:"information_schema.columns"`
	TableCatalog    string
	TableSchema     string
	TblName         string `pg:"tbl_name"`
	ColumnName      string
	OrdinalPosition string
	ColumnDefault   string
	IsNullable      string
	DataType        string
}
