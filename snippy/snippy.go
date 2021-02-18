package snippy

import (
	"embed"
	"github.com/stoewer/go-strcase"
	"log"
	"os"
	"path"
)

type Doer interface {
	Do() error
}

type Snippy struct {
	WD           string
	TemplatePath string
	OutputDir    string
	OutputFile   string
	Module       string
	Data         interface{}
	KV           map[string]string
	ExporteData  map[string][]string
	CamelCase    func(s string) string
	LowCamelCase func(s string) string
	SnakeCase    func(s string) string
	fs           embed.FS
}

func (s Snippy) Name() string {
	return s.TemplatePath
}

func NewSnippy(fs embed.FS, outputDir string) Snippy {
	wd, _ := os.Getwd()
	return Snippy{
		WD:           wd,
		OutputDir:    outputDir,
		ExporteData:  map[string][]string{},
		KV:           map[string]string{},
		CamelCase:    strcase.UpperCamelCase,
		LowCamelCase: strcase.LowerCamelCase,
		SnakeCase:    strcase.SnakeCase,
		fs:           fs,
	}
}

func (s Snippy) Prepare(templatePath string, withValue interface{}, toPath ...string) Doer {
	s.TemplatePath = templatePath
	s.Data = withValue
	s.OutputFile = path.Join(toPath...)
	return s
}

func (s Snippy) Do() error {
	return s.processTemplate()
}

func Run(templates ...Doer) error {
	for i, r := range templates {
		err := r.Do()
		if err != nil {
			log.Fatal(i, err)
			return err
		}
	}
	return nil
}
