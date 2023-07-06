package utils

// Modified version of experimental slices library -> https://pkg.go.dev/golang.org/x/exp@v0.0.0-20230626212559-97b1e661b5df/slices#Delete
func DeleteIndex[S ~[]E, E any](s S, i int, v ...E) S {
	_ = s[i : i+1] // bounds check

	return append(s[:i], s[i+1:]...)
}
