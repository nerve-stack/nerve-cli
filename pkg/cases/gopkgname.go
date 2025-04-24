package cases

import (
	"bytes"
	"unicode"
)

// ToGoPkgName converts a valid directory name to a valid package name.
// It leaves only letters and digits. All letters are mapped to lower case.
func ToGoPkgName(dirName string) string {
	var buf bytes.Buffer

	for _, ch := range dirName {
		switch {
		case unicode.IsLetter(ch):
			buf.WriteRune(unicode.ToLower(ch))

		case buf.Len() != 0 && unicode.IsDigit(ch):
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}
