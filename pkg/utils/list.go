package utils

import (
	"strings"
)

func ToList(csv string) []string {
	return toList(csv)
}

func toList(csv string) []string {
	var list []string

	strlist := strings.Split(csv, ",")
	for _, n := range strlist {
		if ok := contains(list, n); !ok {
			list = append(list, n)
		}
	}

	return list
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}

	return false
}

func find(collection []interface{}, index int) interface{} {
	for i, v := range collection {
		if i == index {
			return v
		}
	}

	return nil
}
