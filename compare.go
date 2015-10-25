package rencode

import (
	"bytes"
)

// this hack allows fetching keys by either string or byte slice type
// no other special matching is performed with other slice types
func deepEqual(a, b interface{}) bool {
	switch a.(type) {
	case string:
		switch b.(type) {
		case []byte:
			return bytes.Compare([]byte(a.(string)), b.([]byte)) == 0
		case string:
			return a == b
		default:
			return false
		}
	case []byte:
		switch b.(type) {
		case []byte:
			return bytes.Compare(a.([]byte), b.([]byte)) == 0
		case string:
			return bytes.Compare([]byte(b.(string)), a.([]byte)) == 0
		default:
			return false
		}
	case Dictionary:
		switch b.(type) {
		case Dictionary:
			d1 := a.(Dictionary)
			d2 := b.(Dictionary)
			return d1.Equals(&d2)
		}
		return false
	case List:
		switch b.(type) {
		case List:
			l1 := a.(List)
			l2 := b.(List)
			return l1.Equals(&l2)
		}
		return false
	}

	return a == b
}
