package sliceutils

func Coalesce[T comparable](tt ...T) T {
	var zero T
	for _, t := range tt {
		if t != zero {
			return t
		}
	}
	return zero
}
