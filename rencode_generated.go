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
	"fmt"
	"math"
	"math/big"
)

func (r *Encoder) encodeSingle(data interface{}) error {
	if data == nil {
		return r.EncodeNone()
	}
	switch x := data.(type) {
	case big.Int:
		s := x.String()
		if len(s) > MAX_INT_LENGTH {
			return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
		}
		return r.EncodeBigNumber(s)
	case List:
		if x.Length() < LIST_FIXED_COUNT {
			_, err := r.w.Write([]byte{byte(LIST_FIXED_START + x.Length())})
			if err != nil {
				return err
			}
			for _, v := range x.Values() {
				err = r.Encode(v)
				if err != nil {
					return err
				}
			}
			return nil
		}
		_, err := r.w.Write([]byte{byte(CHR_LIST)})
		if err != nil {
			return err
		}

		for _, v := range x.Values() {
			err = r.Encode(v)
			if err != nil {
				return err
			}
		}

		_, err = r.w.Write([]byte{byte(CHR_TERM)})
		return err
	case Dictionary:
		if x.Length() < DICT_FIXED_COUNT {
			_, err := r.w.Write([]byte{byte(DICT_FIXED_START + x.Length())})
			if err != nil {
				return err
			}
			keys := x.Keys()
			for i, v := range x.Values() {
				err = r.Encode(keys[i])
				if err != nil {
					return err
				}
				err = r.Encode(v)
				if err != nil {
					return err
				}
			}
			return nil
		}
		_, err := r.w.Write([]byte{byte(CHR_DICT)})
		if err != nil {
			return err
		}
		keys := x.Keys()
		for i, v := range x.Values() {
			err = r.Encode(keys[i])
			if err != nil {
				return err
			}
			err = r.Encode(v)
			if err != nil {
				return err
			}
		}

		_, err = r.w.Write([]byte{byte(CHR_TERM)})
		return err
	case bool:
		return r.EncodeBool(x)
	case float32:
		return r.EncodeFloat32(x)
	case float64:
		return r.EncodeFloat64(x)
	case []byte:
		return r.EncodeBytes(x)
	case string:
		// all strings will be treated as byte arrays
		return r.EncodeBytes([]byte(x))
	case int8:
		return r.EncodeInt8(x)
	case int:
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if math.MinInt32 <= x && x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
		return r.EncodeInt64(int64(x))
	case uint8:
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
	case uint16:
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
	case int16:
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		return r.EncodeInt16(int16(x))
	case uint32:
		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
	case int32:
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		return r.EncodeInt32(int32(x))
	case int64:
		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}
		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}
		if math.MinInt32 <= x && x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}
		return r.EncodeInt64(int64(x))
	case uint64, uint:
		s := fmt.Sprintf("%d", data)
		if len(s) > MAX_INT_LENGTH {
			return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
		}
		return r.EncodeBigNumber(s)
	default:
		return fmt.Errorf("could not encode data of type %T", data)
	}
	panic("unexpected fallthrough")
}
func convertAssignInteger(src, dest interface{}) error {
	switch sv := src.(type) {
	case big.Int:
		switch dv := dest.(type) {
		case *big.Int:
			*dv = sv
			return nil
		}
	case int64:
		switch dv := dest.(type) {
		case *int64:
			*dv = sv
			return nil
		case *int16:
			if sv > math.MaxInt16 || sv < math.MinInt16 {
				return ConversionOverflow{"int64", "int16"}
			}
			*dv = int16(sv)
			return nil
		case *int32:
			if sv > math.MaxInt32 || sv < math.MinInt32 {
				return ConversionOverflow{"int64", "int32"}
			}
			*dv = int32(sv)
			return nil
		case *int:
			*dv = int(sv)
			return nil
		case *int8:
			if sv > math.MaxInt8 || sv < math.MinInt8 {
				return ConversionOverflow{"int64", "int8"}
			}
			*dv = int8(sv)
			return nil
		}
	case int:
		switch dv := dest.(type) {
		case *int:
			*dv = sv
			return nil
		case *int8:
			if sv > math.MaxInt8 || sv < math.MinInt8 {
				return ConversionOverflow{"int", "int8"}
			}
			*dv = int8(sv)
			return nil
		case *int16:
			if sv > math.MaxInt16 || sv < math.MinInt16 {
				return ConversionOverflow{"int", "int16"}
			}
			*dv = int16(sv)
			return nil
		case *int32:
			if sv > math.MaxInt32 || sv < math.MinInt32 {
				return ConversionOverflow{"int", "int32"}
			}
			*dv = int32(sv)
			return nil
		case *int64:
			*dv = int64(sv)
			return nil
		}
	case int8:
		switch dv := dest.(type) {
		case *int8:
			*dv = sv
			return nil
		case *int64:
			*dv = int64(sv)
			return nil
		case *int:
			*dv = int(sv)
			return nil
		case *int16:
			*dv = int16(sv)
			return nil
		case *int32:
			*dv = int32(sv)
			return nil
		}
	case uint8:
		switch dv := dest.(type) {
		case *uint8:
			*dv = sv
			return nil
		case *uint16:
			*dv = uint16(sv)
			return nil
		case *uint32:
			*dv = uint32(sv)
			return nil
		}
	case uint16:
		switch dv := dest.(type) {
		case *uint16:
			*dv = sv
			return nil
		case *uint8:
			if sv > math.MaxUint8 {
				return ConversionOverflow{"uint16", "uint8"}
			}
			*dv = uint8(sv)
			return nil
		case *uint32:
			*dv = uint32(sv)
			return nil
		}
	case int16:
		switch dv := dest.(type) {
		case *int16:
			*dv = sv
			return nil
		case *int32:
			*dv = int32(sv)
			return nil
		case *int64:
			*dv = int64(sv)
			return nil
		case *int:
			*dv = int(sv)
			return nil
		case *int8:
			if sv > math.MaxInt8 || sv < math.MinInt8 {
				return ConversionOverflow{"int16", "int8"}
			}
			*dv = int8(sv)
			return nil
		}
	case uint32:
		switch dv := dest.(type) {
		case *uint32:
			*dv = sv
			return nil
		case *uint8:
			if sv > math.MaxUint8 {
				return ConversionOverflow{"uint32", "uint8"}
			}
			*dv = uint8(sv)
			return nil
		case *uint16:
			if sv > math.MaxUint16 {
				return ConversionOverflow{"uint32", "uint16"}
			}
			*dv = uint16(sv)
			return nil
		}
	case int32:
		switch dv := dest.(type) {
		case *int32:
			*dv = sv
			return nil
		case *int:
			*dv = int(sv)
			return nil
		case *int8:
			if sv > math.MaxInt8 || sv < math.MinInt8 {
				return ConversionOverflow{"int32", "int8"}
			}
			*dv = int8(sv)
			return nil
		case *int16:
			if sv > math.MaxInt16 || sv < math.MinInt16 {
				return ConversionOverflow{"int32", "int16"}
			}
			*dv = int16(sv)
			return nil
		case *int64:
			*dv = int64(sv)
			return nil
		}
	}
	return fmt.Errorf("cannot convert from %T into %T", src, dest)
}
