package helper

import "strings"

// EnumerateAssignments returns all truth assignments for the provided variables.
// The assignments are generated in deterministic order (increasing bit mask).
func EnumerateAssignments(vars []string) []map[string]bool {
	if len(vars) == 0 {
		return []map[string]bool{{}}
	}

	total := 1 << len(vars)
	result := make([]map[string]bool, 0, total)

	for mask := 0; mask < total; mask++ {
		assign := make(map[string]bool, len(vars))
		for idx, name := range vars {
			assign[name] = (mask>>idx)&1 == 1
		}
		result = append(result, assign)
	}

	return result
}

// AssignmentStringify renders an assignment according to the provided variable order.
func AssignmentStringify(vars []string, assign map[string]bool) string {
	if len(vars) == 0 {
		return ""
	}
	parts := make([]string, len(vars))
	for i, name := range vars {
		val := "F"
		if assign[name] {
			val = "T"
		}
		parts[i] = name + "=" + val
	}
	return strings.Join(parts, ", ")
}
