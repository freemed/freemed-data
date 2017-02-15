package common

import (
	"fmt"
	"strings"
)

func InsertsFromArrays(table string, fields []string, data [][]string) {
	for _, rec := range data {
		fmt.Printf("INSERT INTO `%s` ( %s ) VALUES ( ", table, strings.Join(fields, ", "))
		for j, f := range rec {
			if f == "NULL" {
				fmt.Printf("NULL")
			} else {
				// Quote field
				fmt.Printf(`"%s"`, strings.Replace(f, `"`, `\"`, -1))
			}

			// Add comma between fields
			if j < len(rec)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(" );\n")
	}
}

func OneToMultiArray(src []string) [][]string {
	dest := [][]string{}

	for iter, _ := range src {
		if src[iter] != "" {
			dest = append(dest, []string{src[iter]})
		}
	}

	return dest
}
