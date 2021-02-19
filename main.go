package main

import (
	"embed"
	"flag"
	pg "github.com/go-pg/pg/v11"
	"github.com/iamgoroot/postgrestman/pg_info"
	"github.com/iamgoroot/postgrestman/snippy"
	"github.com/stoewer/go-strcase"
	"log"
)

//go:embed templates/request/* templates/root/*
var templates embed.FS

func main() {
	out := flag.String("out", "out", "-out out_dir")
	conn := flag.String("conn", "postgres://user:pass@localhost:5432/db_name", "-url postgres://user:pass@localhost:5432/db_name")
	module := flag.String("module", "generated/module", "-module github.com/{username}/{pkg}")
	assignPort := flag.Int("assignPort", 1234, "-assignPort 8080")

	flag.Parse()

	opt, err := pg.ParseURL(*conn)
	if err != nil {
		panic(err)
	}
	pgInfoRepo := pg_info.PgInfoRepo{DB: pg.Connect(opt)}
	defer pgInfoRepo.Close(nil)

	crawler := PgCrawler{pgInfoRepo: pgInfoRepo, DB: opt.Database, Schema: "public", GoModule: *module}

	snp := snippy.NewSnippy(templates, *out)
	if err := snippy.Run(

		snp.Prepare("templates/root/gomod.tmpl", *module, "go.mod"),

		snp.Prepare("templates/root/Dockerfile.tmpl", struct {
			Conn   string
			Module string
			Port   int
		}{*conn, *module, *assignPort}, "Dockerfile"),

		snp.Prepare("templates/root/docker-compose.tmpl", struct {
			Port int
		}{*assignPort}, "docker-compose.yaml"),
	); err != nil {
		log.Panicln(err)
	}
	templateRoot(crawler, snp)
	pertemplateEndpoint(crawler, snp)
}

func templateRoot(crawler PgCrawler, snp snippy.Snippy) {
	item, err := crawler.GetSetup()
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
}

func pertemplateEndpoint(crawler PgCrawler, snp snippy.Snippy) {
	cruds, err := crawler.GetSetup()
	if err != nil {
		log.Panicln("Could not read DB tables list", err)
	}
	for _, table := range cruds.Names {
		item, err := crawler.Read(table)
		if err != nil {
			log.Panicln(err)
		}
		if err := snippy.Run(
			snp.Prepare("templates/request/api.tmpl", item, "api/openapi-"+strcase.KebabCase(table)+".yaml"),
			snp.Prepare("templates/request/gearbox.tmpl", item, "internal/route", table+".go"),
			snp.Prepare("templates/request/model.tmpl", item, "internal/models", table+".go"),
		); err != nil {
			log.Panicln(err)
		}
	}
}
