package slices

// Find the element that matches the given expression.
// The s = slice of collection elements.
// The f = expression function that should be return a boolean.
func Find[T any](s []T, f func(T) bool) T {
	def := new(T)
	for _, v := range s {
		if found := f(v); found {
			return v
		}
	}
	return *def
}

func FindPT[T any](s []T, f func(T) bool) *T {
	// def := new(T)
	for _, v := range s {
		if found := f(v); found {
			return &v
		}
	}
	return nil
}

// Return -1 if the given expression is false
func FindIdx[T any](s []T, f func(T) bool) (int, T) {
	def := new(T)
	idx := -1
	for i, v := range s {
		if found := f(v); found {
			return i, v
		}
	}
	return idx, *def
}

func First[T any](s []T) T {
	if len(s) == 0 {
		return *new(T)
	}
	return s[0]
}

func Last[T any](s []T) T {
	if len(s) == 0 {
		return *new(T)
	}
	return s[len(s)-1]
}

func ForEach[T any](s []T, f func(*T)) {
	for i := 0; i < len(s); i++ {
		f(&s[i])
	}
}

func Copy[T any](s []T) []T {
	n := make([]T, len(s))
	copy(n, s)
	return n
}

func Contains[T comparable](s []T, v T) bool {
	f := FindPT[T](s, func(t T) bool { return t == v })
	return f != nil
}
