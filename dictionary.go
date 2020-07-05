//
// go-rencode v0.1.6 - Go implementation of rencode - fast (basic)
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

package rencode

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
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

// Get returns the value in the dictionary corresponding to the specified key.
// If the key is not found then 'nil, false' is returned instead.
// Keys of type 'string' and '[]byte' are both compared as if they were strings.
//NOTE: slice keys cannot be used with this method.
func (d *Dictionary) Get(key interface{}) (interface{}, bool) {
	// normalize a byte array key to string
	if keyAsByteArray, ok := key.([]byte); ok {
		return d.internalGet(string(keyAsByteArray))
	}

	// any other key type
	// notice that here slice keys will not be treated in any particular way
	return d.internalGet(key)
}

func (d *Dictionary) internalGet(key interface{}) (interface{}, bool) {
	for i, k := range d.keys {
		// convert byte array keys to string
		if kAsByteArray, ok := k.([]byte); ok {
			k := string(kAsByteArray)

			if k == key {
				return d.values[i], true
			}
			continue
		}

		// generic inteface comparison
		if k == key {
			return d.values[i], true
		}
	}

	return nil, false
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

// ToSnakeCase will convert a 'CamelCase' string to the corresponding 'snake_case' representation.
// Acronyms are converted to lower-case and preceded by an underscore.
func ToSnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return string(out)
}

type RemainingFieldsError struct {
	error error
	remainingFields map[string]interface{}
}

func (e *RemainingFieldsError) Error() string {
	return e.error.Error()
}

func (e *RemainingFieldsError) Fields() map[string]interface{} {
	return e.remainingFields
}

// ToStruct will map a Dictionary into a struct, recursively.
// All dictionary keys must map to a field or an error will be returned.
// It is possible to exclude fields with a specific annotation.
func (d *Dictionary) ToStruct(dest interface{}, excludeAnnotationTag string) error {
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
	l := t.NumField()
	rf := make(map[string]interface{})
	for i := 0; i < l; i++ {
		f := t.Field(i)
		// destination field
		ivf := iv.Field(i)
		name := ToSnakeCase(f.Name)
		skippable := false
		exclude := false

		if rencodeTag, ok := f.Tag.Lookup("rencode"); ok {
			tags := strings.Split(rencodeTag, ",")
			for _, tag := range tags {
				if excludeAnnotationTag != "" && tag == excludeAnnotationTag {
					exclude = true
				} else if tag == "skippable" {
					skippable = true
				}
			}
		}
		if exclude {
			// skip this field
			delete(tmp, name)
			continue
		}

		// see if this field is available
		v, ok := tmp[name]
		if !ok {
			if skippable { continue }
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

					err = d.ToStruct(obj.Interface(), excludeAnnotationTag)
					if err != nil {
						if cvtErr, ok := err.(*RemainingFieldsError); ok {
							rf[fmt.Sprintf("%s_%d", name, i)] = cvtErr.Fields()
							err = nil
						} else {
							return fmt.Errorf("slice field %q: %v", f.Name, err)
						}
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
		for k, v := range tmp {
			rf[k] = v
		}
		return &RemainingFieldsError{
			error:           fmt.Errorf("%d fields left after parsing: %v", len(tmp), tmp),
			remainingFields: rf,
		}
	}

	return nil
}
