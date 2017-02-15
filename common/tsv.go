package common

import (
	"encoding/csv"
	"fmt"
	"os"
)

func TsvFromArrays(filename string, data [][]string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	w := csv.NewWriter(fp)
	w.Comma = '\t'

	for _, rec := range data {
		if err := w.Write(rec); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func ReadTsv(filename string) ([][]string, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer fp.Close()
	r := csv.NewReader(fp)
	r.Comma = '\t'
	return r.ReadAll()
}

func PrependUniqueIds(data [][]string) [][]string {
	processed := [][]string{}
	for k, v := range data {
		processed = append(processed, append([]string{fmt.Sprintf("%d", k+1)}, v...))
	}
	return processed
}
