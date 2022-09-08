package tools

import (
	"flag"
	"fmt"
)

type stringEnum struct {
	value string
}

func (s *stringEnum) Set(v string) error {
	switch v {
	case "http", "https":
		s.value = v
		return nil
	default:
		return fmt.Errorf("must be one of %s", "[http,https]")
	}
}

func (s *stringEnum) String() string {
	return s.value
}

func StringEnumVar(name string, value string, usage string) *string {
	s := stringEnum{value: value}
	flag.CommandLine.Var(&s, name, usage)
	return &(s.value)
}
