//
// go-rencode v0.1.8 - Go implementation of rencode - fast (basic)
//                  object serialization similar to bencode
// Copyright (C) 2015~2019 gdm85 - https://github.com/gdm85/go-rencode/

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

/*
Package rencode is a Go implementation of https://github.com/aresch/rencode

The rencode logic is similar to bencode (https://en.wikipedia.org/wiki/Bencode).
For complex, heterogeneous data structures with many small elements, r-encodings take up significantly less space than b-encodings.

Usage

Example of encoder construction and use:

	b := bytes.Buffer{}
	e := rencode.NewEncoder(&b)

	err := e.Encode(100, true, "hello world", rencode.NewList(42, "nesting is awesome"), 3.14, rencode.Dictionary{})

You can use either specific methods to encode one of the supported types, or the interface-generic Encode() method.

Example of decoder construction:

	e := rencode.NewDecoder(&b)


The DecodeNext() method can be used to decode the next value from the rencode stream; however this method returns an interface{}
while it is usually the norm that there is an expected type instead; in such cases, it is advised to use the Scan() method instead,
which accepts a pointer to any of the supported types.

Example:

	var i int
	var b bool
	var s string
	var l rencode.List
	err := e.Scan(&i, &b, &s, &l)


Supported types

Only the following types are supported:

 - rencode.List
 - rencode.Dictionary
 - big.Int (any integer with more than 63 bits of information)
 - bool
 - float32, float64
 - []byte, string (all strings are stored as byte slices anyway)
 - int8, int16, int32, int64, int
 - uint8, uint16, uint32, uint64, uint

Accessory types

The rencode.List and rencode.Dictionary implement Python-alike features and can store values and keys of
the simpler types enumerated above.

*/
package rencode
