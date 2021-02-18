package main

import (
	"fmt"
	"github.com/getlantern/deepcopy"
	"github.com/iamgoroot/postgrestman/data"
	"github.com/iamgoroot/postgrestman/pg_info"
	"github.com/mpvl/unique"
)

const ui = "rapidoc" //TODO: make swagger alternative

type PgCrawler struct {
	pgInfoRepo pg_info.PgInfoRepo
	DB         string
	Schema     string
	GoModule   string
}

func (r PgCrawler) Root() (data.Root, error) {
	tables, err := r.pgInfoRepo.GetUserTables(r.DB, r.Schema)
	if err != nil {
		return data.Root{}, fmt.Errorf("failed get table list %w", err)
	}
	return data.Root{Names: tables, OpenApiUI: ui, Module: r.GoModule}, err
}

func (r PgCrawler) Read(name string) (data.Crud, error) {
	root, err := r.Root()
	if err != nil {
		fmt.Println("failed get cols for table", name, err)
		return data.Crud{}, err
	}
	columns, err := r.pgInfoRepo.GetUserColumns(r.DB, r.Schema, name)
	if err != nil {
		fmt.Println("failed get cols for table", name, err)
		return data.Crud{Root: root}, err
	}
	entities := make([]data.RequestEntity, len(columns))
	for i, column := range columns {
		entities[i] = colToRequestEntity(column)
	}
	keys := r.pgInfoRepo.GetIndexesForTable(r.Schema, name)
	var filterable []string
	for _, key := range keys {
		filterable = append(filterable, key.Filterable()...)
	}
	unique.Strings(&filterable)
	return data.Crud{
		Name:       name,
		Root:       root,
		Entities:   entities,
		GoType:     convPgToGoType,
		Filterable: filterable,
	}, nil
}

func colToRequestEntity(c pg_info.Column) (data data.RequestEntity) {
	_ = deepcopy.Copy(&data, c)
	data.TableName = c.TblName
	return data
}

func convPgToGoType(s string) data.TypeID {
	switch s {
	case "text":
		return data.TypeID{Name: "string"}
	case "bigint", "integer", "smallint":
		return data.TypeID{Name: "int64"}
	case "boolean":
		return data.TypeID{Name: "bool"}
	case "timestamp", "timestamp with time zone":
		return data.TypeID{Name: "time.Time", Pkg: "time"}

	}
	//TODO:
	return data.TypeID{Name: s}
}
