package tools

import "regexp"

func Filter(filter string, slice []string) []string {
	if filter == ".*" || filter == "" {
		return slice
	}

	s := slice[:0]
	r := regexp.MustCompile(filter)
	for _, v := range slice {
		if r.MatchString(v) {
			s = append(s, v)
		}
	}
	return s
}
