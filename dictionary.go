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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrKeyAlreadyExists is the error returned when the specified key is already defined within the dictionary
	ErrKeyAlreadyExists = errors.New("key already exists in dictionary")
)

// Dictionary is a rencode-specific dictionary that allows any type of key to be mapped to any type of value
type Dictionary struct {
	values []interface{}
	keys   []interface{}
}

// Length returns the total count of elements
func (d *Dictionary) Length() int {
	return len(d.values)
}

// Keys returns all defined keys
func (d *Dictionary) Keys() []interface{} {
	return d.keys
}

// Values returns all stored values
func (d *Dictionary) Values() []interface{} {
	return d.values
}

// Add appends a new (key, value) pair; does not check if key already exists.
func (d *Dictionary) Add(key, value interface{}) {
	d.keys = append(d.keys, key)
	d.values = append(d.values, value)
}

// Zip returns a map with strings as keys or an error if a duplicate key exists.
func (d *Dictionary) Zip() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for i, k := range d.keys {
		var sv string
		v, ok := k.([]uint8)
		if !ok {
			sv, ok = k.(string)
			if !ok {
				return nil, fmt.Errorf("found key with type %T, expected []uint8 or string", k)
			}
		} else {
			sv = string(v)
		}
		if _, ok := result[sv]; ok {
			return nil, ErrKeyAlreadyExists
		}
		result[sv] = d.values[i]
	}

	return result, nil
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

// ToSnakeCase will convert a 'CamelCase' string to the corresponding 'under_score' representation.
func ToSnakeCase(name string) string {
	if len(name) == 0 {
		return ""
	}
	delim := '_'
	b := bytes.NewBufferString("")

	var prev rune
	var curr rune
	for _, next := range name {
		if curr == delim {
			if prev != delim {
				b.WriteRune(delim)
			}
		} else if isUpper(curr) {
			if isLower(prev) || (isUpper(prev) && isLower(next)) {
				b.WriteRune(delim)
			}
			b.WriteRune(curr)
		} else if curr != 0 {
			b.WriteRune(curr)
		}
		prev = curr
		curr = next
	}

	if isUpper(curr) && isLower(prev) && prev != 0 {
		b.WriteRune(delim)
	}
	b.WriteRune(curr)

	return strings.ToLower(b.String())
}

func (d *Dictionary) ToStruct(dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer to struct, got %v", v.Type())
	}

	// get a temporary map with zipped fields
	tmp, err := d.Zip()
	if err != nil {
		return err
	}

	iv := reflect.Indirect(v)
	t := iv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// destination field
		ivf := iv.Field(i)
		name := ToSnakeCase(f.Name)

		// see if this field is available
		v, ok := tmp[name]
		if !ok {
			return fmt.Errorf("field %q: cannot be satisfied", f.Name)
		}

		// special behaviour for slices
		if ivf.Kind() == reflect.Slice {
			// get value as list
			var l List
			err = convertAssign(v, &l)
			if err != nil {
				return fmt.Errorf("slice field %q: value %v: %v", f.Name, v, err)
			}
			// create a new slice
			ns := reflect.MakeSlice(ivf.Type(), l.Length(), l.Length())
			// get element type
			elemType := ivf.Type().Elem()
			for i, v := range l.Values() {
				// all pointed fields are expected to be structs
				if elemType.Kind() == reflect.Struct {
					d, ok := v.(Dictionary)
					if !ok {
						return fmt.Errorf("slice field %q: expected value to be dictionary", f.Name)
					}

					obj := reflect.New(elemType)

					err = d.ToStruct(obj.Interface())
					if err != nil {
						return fmt.Errorf("slice field %q: %v", f.Name, err)
					}

					ns.Index(i).Set(reflect.Indirect(obj))
				} else {
					err = convertAssign(v, ns.Index(i).Addr().Interface())
					if err != nil {
						return fmt.Errorf("slice field %q: value %v: %v", f.Name, v, err)
					}
				}
			}
			ivf.Set(ns)
		} else {
			err = convertAssign(v, ivf.Addr().Interface())
			if err != nil {
				return fmt.Errorf("field %q: value %v: %v", f.Name, v, err)
			}
		}

		// start removing fields that have been used
		delete(tmp, name)
	}

	if len(tmp) != 0 {
		return fmt.Errorf("%d fields left after parsing: %v", len(tmp), tmp)
	}

	return nil
}
