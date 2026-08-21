package policy

import "strings"

// TagFilter 标签匹配。
type TagFilter struct {
	Require []string
	Exclude []string
}

// Match 是否匹配标签集。
func (f TagFilter) Match(tags []string) bool {
	set := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		set[strings.ToLower(t)] = struct{}{}
	}
	for _, r := range f.Require {
		if _, ok := set[strings.ToLower(r)]; !ok {
			return false
		}
	}
	for _, x := range f.Exclude {
		if _, ok := set[strings.ToLower(x)]; ok {
			return false
		}
	}
	return true
}

// MergeTags 去重合并。
func MergeTags(a, b []string) []string {
	set := make(map[string]struct{})
	var out []string
	for _, xs := range [][]string{a, b} {
		for _, t := range xs {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			k := strings.ToLower(t)
			if _, ok := set[k]; ok {
				continue
			}
			set[k] = struct{}{}
			out = append(out, t)
		}
	}
	return out
}
