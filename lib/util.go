package lib

func PrintIf(f func(), condition bool) {
	if condition {
		f()
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
