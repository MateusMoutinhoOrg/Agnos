package utils

// InsertAt places value at position inside list, appending it when position is
// negative or past the end.
func InsertAt[T any](list []T, value T, position int) []T {
	if position < 0 || position >= len(list) {
		return append(list, value)
	}
	out := make([]T, 0, len(list)+1)
	out = append(out, list[:position]...)
	out = append(out, value)
	return append(out, list[position:]...)
}

// RemoveAt drops the element at index from list.
func RemoveAt[T any](list []T, index int) []T {
	out := make([]T, 0, len(list)-1)
	out = append(out, list[:index]...)
	return append(out, list[index+1:]...)
}

// AppendUnique appends each value of extra to values, skipping duplicates.
func AppendUnique(values []string, extra []string) []string {
	for _, candidate := range extra {
		found := false
		for _, existing := range values {
			if existing == candidate {
				found = true
				break
			}
		}
		if !found {
			values = append(values, candidate)
		}
	}
	return values
}

// containsByte reports whether set holds letter.
func containsByte(set string, letter byte) bool {
	for i := 0; i < len(set); i++ {
		if set[i] == letter {
			return true
		}
	}
	return false
}

// contains reports whether list holds value.
func contains(list []string, value string) bool {
	for _, one := range list {
		if one == value {
			return true
		}
	}
	return false
}
