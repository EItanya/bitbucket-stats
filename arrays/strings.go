package arrays

// --------------------------------------
//
//  Helper functions for working with arrays
//
// --------------------------------------

// FilterSTR array map function
func FilterSTR(s []string, f func(string) bool) []string {
	m := make([]string, 0)
	for _, v := range s {
		if f(v) {
			m = append(m, v)
		}
	}
	return m
}

// FindSTR array map function
func FindSTR(s []string, f func(string) bool) string {
	for _, v := range s {
		if f(v) {
			return v
		}
	}
	return ""
}

// IndexOfSTR array indexOf function for strings
func IndexOfSTR(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}
	return -1
}
