package common

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func FileFromZipArchive(data []byte, filename string, debug bool) ([]byte, error) {
	br := bytes.NewReader(data)
	zr, err := zip.NewReader(br, int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		if debug {
			fmt.Printf("## Found file '%s'\n", f.Name)
		}
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
