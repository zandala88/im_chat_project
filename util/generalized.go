package util

// Map
// 去除结构体中的一列
func Map[T any, V any](arr []T, f func(T) V) []V {
	result := make([]V, len(arr))
	for i, v := range arr {
		result[i] = f(v)
	}
	return result
}
