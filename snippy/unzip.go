package snippy

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func (s Snippy) UnzipHttp(url string, targetDir string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatal(err)
	}
	for _, zipFile := range zipReader.File {
		_ = unZipFile(zipFile, path.Join(s.OutputDir, targetDir))
	}

}

func unZipFile(zf *zip.File, targetDir string) error {
	f, err := zf.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	outFile, err := os.Create(path.Join(targetDir, zf.Name))
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, f)
	return err
}
