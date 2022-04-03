package utils

// Contains evaluates whether a string is in the defined slice
func Contains(slice []string, inputValue string) bool {
	for _, sliceValue := range slice {
		if sliceValue == inputValue {
			return true
		}
	}
	return false
}
