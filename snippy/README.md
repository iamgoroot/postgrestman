[Snippy](#Snippy)

- [Motivation](#motivation)
- [Status](#status)
- [Getting started](#getting-started)

# Snippy

## Motivation

Execute go templates in few lines

### Status

Project is in pre-alpha status.

## Getting started

One template:

```
snippy.
    NewSnippy("output/dir").
    Prepare("templates/input.tmpl", data, "output.file").
    Do()
```

Batch:

```
	if err := snippy.Run(
		snp.Prepare("templates/root/openapi.tmpl", item, "api/openapi.yaml"),
		snp.Prepare("templates/root/gearbox_main.tmpl", item, "cmd/gearbox/main.go"),
		snp.Prepare("templates/root/app.tmpl", item, "cmd/gearbox/app.go"),
	); err != nil {
		log.Panicln(err)
	}
```
