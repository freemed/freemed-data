package common

import (
	"strings"
)

func Derivatives(values [][]string, idx int, sep string) []string {
	out := make([]string, 0)
	for _, v := range values {
		if strings.Index(v[idx], sep) != -1 {
			sub := strings.Split(v[idx], sep)
			for _, sv := range sub {
				out = append(out, strings.TrimSpace(sv))
			}
			continue
		}
		out = append(out, strings.TrimSpace(v[idx]))
	}
	RemoveDuplicates(&out)
	return out
}

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
