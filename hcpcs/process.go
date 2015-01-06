package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	fixed "github.com/jbuchbinder/gofixedfield"
	http "github.com/jbuchbinder/gosimplehttp"
	"io"
	"strings"
)

const (
	HcpcsZipUrl = "http://www.cms.gov/Medicare/Coding/HCPCSReleaseCodeSets/Downloads/14anweb.zip"
	HcpcsFile   = "HCPC2014_ANWEB.txt"
)

func main() {
	fmt.Printf("HTTP GET : %s\n", HcpcsZipUrl)
	code, file, _, err := http.SimpleGet(HcpcsZipUrl)
	if err != nil {
		panic(err)
	}
	if code > 299 {
		fmt.Printf("HTTP request got result code %d\n", code)
		return
	}

	fmt.Printf("Extract file %s from archive\n", HcpcsFile)
	contents, err := FileFromZipArchive(file, HcpcsFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decompressed %d bytes from ZIP archive\n", len(contents))
	rec := strings.Split(string(contents), fixed.EOL_DOS)
	//rec, err := fixed.RecordsFromFile("HCPC2013_A-N.txt", fixed.EOL_DOS)
	if err != nil {
		panic(err)
	}

	p := make([]HcpcsRecord, 0)
	for i := range rec {
		if len(rec[i]) > 100 {
			fmt.Printf("Processing record %d\n", i)
			var out HcpcsRecord
			fixed.Unmarshal(rec[i], &out)
			p = append(p, out)
		}
	}

}

func FileFromZipArchive(data []byte, filename string) ([]byte, error) {
	br := bytes.NewReader(data)
	zr, err := zip.NewReader(br, int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		fmt.Printf("Found file '%s'\n", f.Name)
		if f.Name == filename {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			contents := make([]byte, f.UncompressedSize)
			_, err = io.ReadFull(rc, contents)
			if err != nil {
				return nil, err
			}
			return contents, nil
		}
	}
	return nil, nil
}
