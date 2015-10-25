/*
 * go-rencode v0.1.0 - Go implementation of rencode - fast (basic)
                    object serialization similar to bencode
 * Copyright (C) 2015 gdm85 - https://github.com/gdm85/go-rencode/

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/
package rencode

import (
	"errors"
)

var (
	ErrKeyNotFound = errors.New("key not found in dictionary")
)

type List struct {
	values []interface{}
}

func (l *List) Add(value interface{}) {
	l.values = append(l.values, value)
}

func (l *List) Values() []interface{} {
	return l.values
}

func (l *List) Get(i int) (interface{}, error) {
	if i < 0 || i >= len(l.values) {
		return nil, ErrKeyNotFound
	}

	return l.values[i], nil
}

func (l *List) Length() int {
	return len(l.values)
}
