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

/*
Package rencode is a Go implementation of https://github.com/aresch/rencode

The rencode logic is similar to bencode (https://en.wikipedia.org/wiki/Bencode).
For complex, heterogeneous data structures with many small elements, r-encodings take up significantly less space than b-encodings.

Usage

You can use either specific methods to encode one of the supported types, or the interface-generic Encode() method.

The DecodeNext() method can be used to decode the next value from the rencode stream.

*/
package rencode
