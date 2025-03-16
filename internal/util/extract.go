package util

func ExtractPart[T any](parts []T, position int, fallback T) T {
	if position >= 0 && position < len(parts) {
		return parts[position]
	}
	return fallback
}
