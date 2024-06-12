// Slices functions.
// Custom functions that operate on slice.

package utils

import "fmt"

// ContainsStringIndex reports whether e is present in s,
// and index of the first occurrence
func ContainsStringIndex(s []string, e string) (bool, int) {
	for i, elementS := range s {
		if elementS == e {
			return true, i
		}
	}
	return false, -1
}

// RemoveIndexString remove the element at index i from slice
// Maintains order of the slice
func RemoveIndexString(s []string, i int) ([]string, error) {
	if len(s) > 0 && i < len(s) {
		j := i + 1
		_ = s[i:j]
		return append(s[:i], s[j:]...), nil
	} else {
		return nil, fmt.Errorf("RemoveIndexString empty slice or index greater than slice size")
	}
}

func Slice2SliceString(s []interface{}) []string {
	sliceString := make([]string, len(s))
	for i, e := range s {
		sliceString[i] = fmt.Sprintf("%v", e)
	}
	return sliceString
}
