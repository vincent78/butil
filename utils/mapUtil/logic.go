package mapUtil

import "errors"

/*
 *
 * 从m中选出小于s的 最大值
 *
 */
func Location[T any](s int, m map[int]T) (T, error) {
	if len(m) == 0 {
		t := new(T)
		return *t, errors.New("the map is empty")
	}
	min := 0
	max := 0

	for k, _ := range m {
		if k >= s {
			if max > k {
				max = k
			}
		}

		if k <= s {
			if min < k {
				min = k
			}
		}
	}

	if min == 0 {
		min = max
	}

	if max == 0 {
		max = min
	}

	return m[min], nil
}
