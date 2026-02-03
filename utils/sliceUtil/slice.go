package sliceUtil

// Union 数组并集操作
func Union[T comparable](a, b []T) []T {
	m := make(map[T]struct{})

	// 加入 a
	for _, v := range a {
		m[v] = struct{}{}
	}

	// 加入 b
	for _, v := range b {
		m[v] = struct{}{}
	}

	// map 转 slice
	res := make([]T, 0, len(m))
	for v := range m {
		res = append(res, v)
	}

	return res
}
