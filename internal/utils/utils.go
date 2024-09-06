package utils

func Filter[T any](items []T, pred func(T) bool) []T {
	var res []T
	for _, v := range items {
		if pred(v) {
			res = append(res, v)
		}
	}
	return res
}

func Find[T any](items []T, pred func(T) bool) *T {
	for _, v := range items {
		if pred(v) {
			return &v
		}
	}
	return nil
}
