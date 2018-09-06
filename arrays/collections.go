package arrays

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

func Find(s []string, f func(string) bool) string {
	for _, v := range s {
		if f(v) {
			return v
		}
	}
	return ""
}

func IndexOf(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return -1
}
