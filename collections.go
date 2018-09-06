package main

// Filter array map function
func Filter(s []string, f func(string) bool) []string {
	m := make([]string, 0)
	for _, v := range s {
		if f(v) {
			m = append(m, v)
		}
	}
	return m
}
