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
	ErrKeyAlreadyExists = errors.New("key already exists in dictionary")
)

type Dictionary struct {
	List
	keys []interface{}
}

func (d *Dictionary) Keys() []interface{} {
	return d.keys
}

func (d *Dictionary) Get(key interface{}) (interface{}, error) {
	for i, k := range d.keys {
		if deepEqual(k, key) {
			return d.values[i], nil
		}
	}
	return nil, ErrKeyNotFound
}

func (d *Dictionary) Set(key, value interface{}) bool {
	for i, k := range d.keys {
		if deepEqual(k, key) {
			d.values[i] = value
			return true
		}
	}

	d.keys = append(d.keys, key)
	d.values = append(d.values, value)
	return false
}

func (d *Dictionary) Add(key, value interface{}) error {
	for _, k := range d.keys {
		if deepEqual(k, key) {
			return ErrKeyAlreadyExists
		}
	}

	d.keys = append(d.keys, key)
	d.values = append(d.values, value)
	return nil
}

func (d *Dictionary) Compare(b *Dictionary) bool {
	keys2 := b.Keys()
	for i, k1 := range d.keys {
		if !deepEqual(k1, keys2[i]) {
			return false
		}

		// compare values as well
		if !deepEqual(d.values[i], b.values[i]) {
			return false
		}
	}

	return true
}
