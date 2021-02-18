package main

import (
	"embed"
	"flag"
	"github.com/go-pg/pg"
	"github.com/iamgoroot/postgrestman/data"
	"github.com/iamgoroot/postgrestman/pg_info"
	"github.com/iamgoroot/postgrestman/snippy"
	"log"
)

//go:embed templates/request/* templates/root/*
var templates embed.FS

type Args struct {
	Addr     string
	DB       string
	Schema   string
	User     string
	Password string
	Out      string
	Module   string
}

func main() {
	args := Args{}
	flag.StringVar(&args.Addr, "addr", "localhost:5432", "-addr localhost:5432")
	flag.StringVar(&args.DB, "db", "postgres", "-db batabase")
	flag.StringVar(&args.DB, "schema", "public", "-schema schema")
	flag.StringVar(&args.User, "user", "postgres", "-user user")
	flag.StringVar(&args.Password, "MODULE", "postgres", "-password password")
	flag.StringVar(&args.Out, "out", "out", "-out out_dir")
	flag.StringVar(&args.Module, "module", "this/is/just/module/name", "-module github.com/{username}/{pkg}")
	flag.Parse()

	pgInfoRepo := pg_info.PgInfoRepo{DB: pg.Connect(&pg.Options{
		User:     args.User,
		Password: args.Password,
		Database: args.DB,
		Addr:     args.Addr,
	})}
	defer pgInfoRepo.Close()

	crawler := PgCrawler{pgInfoRepo: pgInfoRepo, DB: args.DB, Schema: args.Schema, GoModule: args.Module}

	snp := snippy.NewSnippy(templates, args.Out)
	if err := snippy.Run(
		snp.Prepare("templates/root/gomod.tmpl", args.Module, "go.mod"),
		snp.Prepare("templates/root/docker-compose.tmpl", args.Module, "docker-compose.yaml"),
	); err != nil {
		log.Panicln(err)
	}
	root := once(crawler, snp)
	perTable(root.Names, crawler, snp)
}

func once(crawler PgCrawler, snp snippy.Snippy) data.Root {
	item, err := crawler.Root()
	if err != nil {
		log.Panicln("Could not read DB tables list", err)
	}
	if err := snippy.Run(
		snp.Prepare("templates/root/openapi.tmpl", item, "api/openapi.yaml"),
		snp.Prepare("templates/root/gearbox_main.tmpl", item, "cmd/gearbox/main.go"),
		snp.Prepare("templates/root/app.tmpl", item, "cmd/gearbox/app.go"),
	); err != nil {
		log.Panicln(err)
	}
	return item
}

func perTable(tables []string, crawler PgCrawler, snp snippy.Snippy) {
	for _, table := range tables {
		item, err := crawler.Read(table)
		if err != nil {
			log.Panicln(err)
		}
		if err := snippy.Run(
			snp.Prepare("templates/request/api.tmpl", item, "api/openapi-"+table+".yaml"),
			snp.Prepare("templates/request/gearbox.tmpl", item, "internal/route", table+".go"),
			snp.Prepare("templates/request/model.tmpl", item, "internal/models", table+".go"),
		); err != nil {
			log.Panicln(err)
		}
	}
}
