// Code generated by "stringer -type=InputKind -linecomment"; DO NOT EDIT.

package shared

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[INPUT_TEXT-0]
	_ = x[INPUT_MARKDOWN-1]
	_ = x[INPUT_HTML-2]
	_ = x[INPUT_WYSIWYG-3]
	_ = x[INPUT_TEXTAREA-4]
	_ = x[INPUT_HIDDEN-5]
	_ = x[INPUT_PASSWORD-6]
	_ = x[INPUT_URL-7]
	_ = x[INPUT_EMAIL-8]
	_ = x[INPUT_COLOR-9]
	_ = x[INPUT_CHECKBOX-10]
	_ = x[INPUT_RADIO-11]
	_ = x[INPUT_DATETIME-12]
	_ = x[INPUT_FILE-13]
	_ = x[INPUT_IMAGE-14]
	_ = x[INPUT_NUMBER-15]
	_ = x[INPUT_RANGE-16]
	_ = x[INPUT_SELECT-17]
	_ = x[INPUT_SEARCH-18]
	_ = x[INPUT_UNKOWN-19]
}

const _InputKind_name = "textmarkdownhtmlwysiwygtextareahiddenpasswordurlemailcolorcheckboxradiodatetimefileimagenumberrangeselectsearchunknown"

var _InputKind_index = [...]uint8{0, 4, 12, 16, 23, 31, 37, 45, 48, 53, 58, 66, 71, 79, 83, 88, 94, 99, 105, 111, 118}

func (i InputKind) String() string {
	if i < 0 || i >= InputKind(len(_InputKind_index)-1) {
		return "InputKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _InputKind_name[_InputKind_index[i]:_InputKind_index[i+1]]
}
