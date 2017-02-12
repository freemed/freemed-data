package common

import (
	"fmt"
	"strings"
)

func InsertsFromArrays(table string, fields []string, data [][]string) {
	for _, rec := range data {
		fmt.Printf("INSERT INTO `%s` ( %s ) VALUES ", table, strings.Join(fields, ", "))
		for j, f := range rec {
			// Quote field
			fmt.Printf(`"%s"`, strings.Replace(f, `"`, `\"`, -1))

			// Add comma between fields
			if j < len(rec)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(" );\n")
	}
}
