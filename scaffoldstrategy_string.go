// Code generated by "stringer -type=ScaffoldStrategy"; DO NOT EDIT.

package crudex

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SCAFFOLD_ALWAYS-0]
	_ = x[SCAFFOLD_IF_NOT_EXISTS-1]
	_ = x[SCAFFOLD_NEVER-2]
}

const _ScaffoldStrategy_name = "SCAFFOLD_ALWAYSSCAFFOLD_IF_NOT_EXISTSSCAFFOLD_NEVER"

var _ScaffoldStrategy_index = [...]uint8{0, 15, 37, 51}

func (i ScaffoldStrategy) String() string {
	if i < 0 || i >= ScaffoldStrategy(len(_ScaffoldStrategy_index)-1) {
		return "ScaffoldStrategy(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ScaffoldStrategy_name[_ScaffoldStrategy_index[i]:_ScaffoldStrategy_index[i+1]]
}
