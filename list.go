package rencode

//
// go-rencode v0.1.1 - Go implementation of rencode - fast (basic)
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

import (
	"errors"
)

var (
	// ErrKeyNotFound is the error returned when specified key does not exist in List or Dictionary
	ErrKeyNotFound = errors.New("key not found")
)

// List is a rencode-specific list that allows any type of value to be concatenated
type List struct {
	values []interface{}
}

// NewList returns a new list with the values specified as arguments
func NewList(values ...interface{}) List {
	return List{values}
}

// Add appends one or more values to the list
func (l *List) Add(values ...interface{}) {
	l.values = append(l.values, values...)
}

// Values returns all values in the list
func (l *List) Values() []interface{} {
	return l.values
}

// Length returns the total count of elements
func (l *List) Length() int {
	return len(l.values)
}
