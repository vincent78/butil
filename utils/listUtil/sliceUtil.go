package listUtil

import "sort"

func unique(src []string) []string {
	var uniq []string

	set := make(map[string]struct{})
	for _, v := range src {
		if _, exist := set[v]; !exist {
			set[v] = struct{}{}
			uniq = append(uniq, v)
		}
	}

	sort.Strings(uniq)

	return uniq
}
