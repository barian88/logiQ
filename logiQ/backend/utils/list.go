package utils

func AreSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	// 用 map 构建一个集合
	set := make(map[int]bool)
	for _, v := range a {
		set[v] = true
	}

	// 检查 b 中的每个元素是否都在 set 中
	for _, v := range b {
		if !set[v] {
			return false
		}
	}

	return true
}
