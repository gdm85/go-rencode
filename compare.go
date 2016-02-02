package rencode

//
// go-rencode v0.1.0 - Go implementation of rencode - fast (basic)
//                  object serialization similar to bencode
// Copyright (C) 2015 gdm85 - https://github.com/gdm85/go-rencode/

// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

import (
	"bytes"
)

// this hack allows fetching keys by either string or byte slice type
// no other special matching is performed with other slice types
func deepEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case string:
		switch bv := b.(type) {
		case []byte:
			return bytes.Compare([]byte(av), bv) == 0
		case string:
			return av == bv
		default:
			return false
		}
	case []byte:
		switch bv := b.(type) {
		case []byte:
			return bytes.Compare(av, bv) == 0
		case string:
			return bytes.Compare([]byte(bv), av) == 0
		default:
			return false
		}
	case Dictionary:
		switch bv := b.(type) {
		case Dictionary:
			return av.Equals(&bv)
		}
		return false
	case List:
		switch bv := b.(type) {
		case List:
			return av.Equals(&bv)
		}
		return false
	}

	return a == b
}
