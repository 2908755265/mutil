package mutil

func InSlice[T comparable](elem T, slice []T) bool {
	for _, t := range slice {
		if t == elem {
			return true
		}
	}
	return false
}
