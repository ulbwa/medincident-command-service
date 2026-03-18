package utils

func Ptr[T any](v T) *T {
	return &v
}

func PtrEquals[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

type Equatable[T any] interface {
	Equals(T) bool
}

func PtrEqualsFunc[T Equatable[T]](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return (*a).Equals(*b)
}
