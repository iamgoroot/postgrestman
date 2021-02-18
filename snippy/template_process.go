package snippy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

func (s *Snippy) processTemplate() error {
	file, err := s.fs.ReadFile(s.TemplatePath)
	if err != nil {
		return err
	}
	output, err := templateItem(string(file), s)
	if err != nil {
		return err
	}
	dir := path.Join(s.WD, s.OutputDir)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		fmt.Println(err)
	}

	outFilePath := path.Join(dir, s.OutputFile)

	err = os.MkdirAll(path.Dir(outFilePath), 0700)
	if err != nil {
		fmt.Println(err)
	}

	return ioutil.WriteFile(outFilePath, output, 0664)
}

func templateItem(templateSrc string, item Exported) ([]byte, error) {
	tmpl, err := template.New(item.Name()).Parse(templateSrc)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	item.ResetExported()
	err = tmpl.Execute(buf, &item)
	if item.Exported() != nil {
		buf.Reset()
		err = tmpl.Execute(buf, &item)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	output := buf.Bytes()
	return output, nil
}
