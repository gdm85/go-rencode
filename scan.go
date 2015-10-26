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
	"errors"
	"fmt"
)

var (
	// ErrConversionOverflow is returned when the scanned integer would overflow the destination integer
	ErrConversionOverflow = errors.New("conversion would overflow integer size")
)

// Scan will scan the decoder data to fill in the specified target objects; if possible,
// a conversion will be performed. If targets have not pointer types or if the conversion is
// not possible, an error will be returned.
func (d *Decoder) Scan(targets ...interface{}) error {
	for i, target := range targets {
		src, err := d.DecodeNext()
		if err != nil {
			return err
		}

		err = convertAssign(src, target)
		if err != nil {
			return fmt.Errorf("scan element %d: %v", i, err)
		}
	}

	return nil
}

// Scan will scan the list to fill in the specified target objects; if possible,
// a conversion will be performed. If targets have not pointer types or if the conversion is
// not possible, an error will be returned.
func (l *List) Scan(targets ...interface{}) error {
	if len(targets) > l.Length() {
		return errors.New("not enough elements in list")
	}
	for i, target := range targets {
		err := convertAssign(l.values[i], target)
		if err != nil {
			return fmt.Errorf("scan element %d: %v", i, err)
		}
	}

	return nil
}

func convertAssign(src, dest interface{}) error {
	switch src.(type) {
	case bool:
		switch dest.(type) {
		case *bool:
			d := dest.(*bool)
			*d = src.(bool)
			return nil
		}
	case List:
		switch dest.(type) {
		case *List:
			d := dest.(*List)
			*d = src.(List)
			return nil
		}
	case Dictionary:
		switch dest.(type) {
		case *Dictionary:
			d := dest.(*Dictionary)
			*d = src.(Dictionary)
			return nil
		}
	case float32:
		switch dest.(type) {
		case *float32:
			d := dest.(*float32)
			*d = src.(float32)
			return nil
		case *float64:
			d := dest.(*float64)
			*d = float64(src.(float32))
			return nil
		}
	case float64:
		switch dest.(type) {
		case *float64:
			d := dest.(*float64)
			*d = src.(float64)
			return nil
		}
	case []byte:
		switch dest.(type) {
		case *[]byte:
			d := dest.(*[]byte)
			*d = src.([]byte)
			return nil
		case *string:
			d := dest.(*string)
			*d = string(src.([]byte))
			return nil
		}
	case string:
		switch dest.(type) {
		case *[]byte:
			d := dest.(*[]byte)
			*d = []byte(src.(string))
			return nil
		case *string:
			d := dest.(*string)
			*d = src.(string)
			return nil
		}
	}

	return convertAssignInteger(src, dest)
}
