package tools

type OrderStringsSlice []string

func (s OrderStringsSlice) Len() int {
	return len(s)
}

func (s OrderStringsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s OrderStringsSlice) Less(i, j int) bool {
	if len(s[i]) != len(s[j]) {
		return len(s[i]) > len(s[j])
	}
	return s[i] > s[j]
}
