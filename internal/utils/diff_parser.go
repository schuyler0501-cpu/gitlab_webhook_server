package utils

import (
	"strings"
)

// ParseDiffStats 从 diff 字符串中解析添加和删除的行数
// diff 格式为标准 unified diff，以 + 开头的是添加行，以 - 开头的是删除行
// 注意：忽略文件头（+++ 和 ---）以及上下文行（以空格开头）
func ParseDiffStats(diff string) (added, removed int) {
	if diff == "" {
		return 0, 0
	}

	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		// 跳过空行
		if len(line) == 0 {
			continue
		}

		// 跳过文件头（+++ 和 ---）
		if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			continue
		}

		// 跳过上下文行（以空格开头）和 hunk 头（以 @@ 开头）
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "@@") {
			continue
		}

		// 统计添加的行（以 + 开头）
		if strings.HasPrefix(line, "+") {
			added++
		}
		// 统计删除的行（以 - 开头）
		if strings.HasPrefix(line, "-") {
			removed++
		}
	}

	return added, removed
}

