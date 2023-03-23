package utils

func StringElementExists(arr []string, x string) bool {
	for _, n := range arr {
		if x == n {
			return true
		}
	}
	return false
}
