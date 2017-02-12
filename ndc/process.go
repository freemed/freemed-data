package main

import (
	"encoding/csv"
	"fmt"
	"github.com/freemed/freemed-data/common"
	http "github.com/jbuchbinder/gosimplehttp"
	"strings"
)

const (
	NdcZipUrl      = "http://www.accessdata.fda.gov/cder/ndctext.zip"
	NdcProductFile = "product.txt"
)

func main() {
	fmt.Printf("HTTP GET : %s\n", NdcZipUrl)
	code, file, _, err := http.SimpleGet(NdcZipUrl)
	if err != nil {
		panic(err)
	}
	if code > 299 {
		fmt.Printf("HTTP request got result code %d\n", code)
		return
	}

	fmt.Printf("Extract product file %s from archive\n", NdcProductFile)
	contents, err := common.FileFromZipArchive(file, NdcProductFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decompressed %d bytes from ZIP archive\n", len(contents))
	r := csv.NewReader(strings.NewReader(string(contents)))
	r.Comma = '\t'
	r.LazyQuotes = true
	rec, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", rec[1:])
}
