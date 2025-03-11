package util

func GetValueOrDefault[T comparable](value T, fallback T) T {
	var zero T
	if value != zero {
		return value
	}
	return fallback
}

func GetPointerValueOrDefault[T comparable](value *T, fallback T) T {
	if value != nil {
		return *value
	}
	return fallback
}
