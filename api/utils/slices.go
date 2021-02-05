package utils

import "strings"

func RemoveDuplicateStrings(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func InSliceOfStrings(slice []string, val string, caseSensitive bool) (bool, int) {
	for i, item := range slice {
		if caseSensitive {
			if item == val {
				return true, i
			}
		} else {
			if strings.ToLower(item) == strings.ToLower(val) {
				return true, i
			}
		}

	}
	return false, -1
}
