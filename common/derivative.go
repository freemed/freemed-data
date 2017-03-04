package common

import (
	"reflect"
	"strconv"
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

func HasElement(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {

			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func CoerceSliceStringToInt(sl []string) []int {
	out := make([]int, 0)
	for _, v := range sl {
		iv, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		out = append(out, iv)
	}
	return out
}

func MinIntSlice(v []int) (m int) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] < m {
			m = v[i]
		}
	}
	return
}

func MaxIntSlice(v []int) (m int) {
	if len(v) > 0 {
		m = v[0]
	}
	for i := 1; i < len(v); i++ {
		if v[i] > m {
			m = v[i]
		}
	}
	return
}
