package pg_info

import (
	"github.com/go-pg/pg"
	"log"
	"strings"
)

type PgInfoRepo struct {
	*pg.DB
}

const distinctTables = "distinct(table_name) as tbl_name "

func (r PgInfoRepo) GetUserTables(db string, schema string) ([]string, error) {
	var columns []string
	_, err := r.Query(&columns,
		"SELECT "+distinctTables+"FROM information_schema.columns "+
			columnsSelectBySchema, db, schema,
	)
	return columns, err
}

const (
	columnsSelectBySchema = "WHERE table_catalog = ? AND table_schema = ? "
)

func (r PgInfoRepo) GetUserColumns(db string, schema string, table string) ([]Column, error) {
	var columns []Column
	_, err := r.Query(&columns,
		"SELECT "+
			columnsFields+
			"FROM information_schema.columns "+
			columnsSelectBySchema+"AND table_name = ? ORDER BY ordinal_position", db, schema, table,
	)
	return columns, err
}

type Index struct {
	Name string
	Src  string
}

func (f Index) Filterable() []string {
	index := strings.LastIndex(f.Src, "(")
	if index < 0 {
		return nil
	}
	values := f.Src[index+1:]
	end := strings.Index(values, ")")
	return strings.Split(values[:end], ",")
}

func (r PgInfoRepo) GetIndexesForTable(schema string, table string) (res []Index) {
	_, err := r.Query(&res,
		`
	SELECT
		indexname as Name,
		indexdef as Src
	FROM
		pg_indexes
	WHERE
		schemaname = ?
	AND
		tablename = ?
	ORDER BY
		tablename,
		indexname;`, schema, table,
	)
	if err != nil {
		log.Fatal(err)
	}
	return
}

const (
	columnsFields = "table_catalog, " +
		"table_schema, " +
		"table_name as tbl_name, " +
		"column_name, ordinal_position, column_default, is_nullable, data_type "
)
