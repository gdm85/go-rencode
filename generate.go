// +build generate

//
// go-rencode v0.1.2 - Go implementation of rencode - fast (basic)
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

package main

import (
	"fmt"
	"strings"
)

// template block starts
const top = `//
// go-rencode v0.1.2 - Go implementation of rencode - fast (basic)
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
		return r.EncodeInt8(x)`

// template block ends

var (
	intTypes map[string]int
)

func init() {
	// NOTE: uint64 is not supported as it can overflow int64, which is
	// the maximum regular integer type for the original Python rencode
	// For values which require more bits than int64, use big.Int
	intTypes = map[string]int{"uint8": 8, "uint16": 16, "int16": 15, "uint32": 32, "int32": 31, "int64": 63}

	if ^uint(0) == uint(^uint32(0)) {
		intTypes["uint"] = 32
		intTypes["int"] = 31
	} else if ^uint(0) == uint(^uint64(0)) {
		// same here, 'uint' is not being defined on purpose
		intTypes["int"] = 63
	} else {
		panic("unrecognized default uint bitsize")
	}
}

func signedGenerate(t string, bitsize int) {
	// all signed ints can be checked against this nibble range
	fmt.Println(`		if math.MinInt8 <= x && x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}`)

	if bitsize == 15 {
		fmt.Println(`		return r.EncodeInt16(int16(x))`)
		return
	}

	if bitsize >= 15 {
		fmt.Println(`		if math.MinInt16 <= x && x <= math.MaxInt16 {
			return r.EncodeInt16(int16(x))
		}`)
	}

	if bitsize == 31 {
		fmt.Println(`		return r.EncodeInt32(int32(x))`)
		return
	}

	if bitsize >= 31 {
		fmt.Println(`		if math.MinInt32 <= x && x <= math.MaxInt32 {
			return r.EncodeInt32(int32(x))
		}`)
	}

	if bitsize == 63 {
		fmt.Println(`		return r.EncodeInt64(int64(x))`)
		return
	}

	panic("signed: using bitsize larger than 64")
}

func unsignedGenerate(t string, bitsize int) {
	// all unsigned ints can be checked against this nibble range
	fmt.Println(`		if x <= math.MaxInt8 {
			return r.EncodeInt8(int8(x))
		}`)

	if bitsize >= 16 {
		fmt.Println(`		if x <= math.MaxInt16 {
		return r.EncodeInt16(int16(x))
		}`)
	}

	if bitsize >= 32 {
		fmt.Println(`		if x <= math.MaxInt32 {
		return r.EncodeInt32(int32(x))
		}`)
		return
	}

	if bitsize == 63 {
		fmt.Println(`		return r.EncodeInt64(int64(x))`)
		return
	}
}

func main() {
	fmt.Println(top)

	for t, bitsize := range intTypes {
		fmt.Printf("	case %s:\n", t)

		if bitsize%2 == 0 {
			// unsigned integer of some bitsize
			unsignedGenerate(t, bitsize)
		} else {
			// signed integer of some bitsize
			signedGenerate(t, bitsize)
		}
	}

	// encoding for 'big numbers'
	caseStr := "uint64"
	if _, ok := intTypes["uint"]; !ok {
		caseStr += ", uint"
	}
	fmt.Printf("\tcase %s:\n", caseStr)
	fmt.Println(`		s := fmt.Sprintf("%d", data)
		if len(s) > MAX_INT_LENGTH {
			return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
		}
		return r.EncodeBigNumber(s)`)

	// tail default case
	fmt.Println(`	default:
		return fmt.Errorf("could not encode data of type %T", data)
	}
	panic("unexpected fallthrough")
}`)

	// generate integer conversion function
	fmt.Println(`func convertAssignInteger(src, dest interface{}) error {
		switch sv := src.(type) {
			case big.Int:
				switch dv := dest.(type) {
					case *big.Int:
						*dv = sv
						return nil
				}`)

	// add int8 to allowed types
	intTypes["int8"] = 7

	for st, sBitsize := range intTypes {
		fmt.Printf(`		case %s:
			switch dv := dest.(type) {
			case *%s:
				*dv = sv
				return nil`+"\n", st, st)
		for dt, dBitsize := range intTypes {
			if dt == st {
				continue
			}

			sUnsigned := sBitsize%2 == 0
			dUnsigned := dBitsize%2 == 0

			// disallow conversions between signed/unsigned
			// user should know if integer is signed/unsinged before scanning for it
			if sUnsigned != dUnsigned {
				continue
			}

			if sUnsigned {
				unsignedConvertGenerate(st, sBitsize, dt, dBitsize)
			} else {
				signedConvertGenerate(st, sBitsize, dt, dBitsize)
			}
		}
		fmt.Println(`		}`)
	}

	fmt.Println(`		}
	return fmt.Errorf("cannot convert from %T into %T", src, dest)
}`)

}

func getMaxValue(t string) string {
	return fmt.Sprintf("math.Max%s%s", strings.ToUpper(fmt.Sprintf("%c", t[0])), t[1:])
}

func getMinValue(t string) string {
	return fmt.Sprintf("math.Min%s%s", strings.ToUpper(fmt.Sprintf("%c", t[0])), t[1:])
}

func unsignedConvertGenerate(sourceType string, sourceBitsize int, destType string, destBitsize int) {
	fmt.Printf("		case *%s:\n", destType)

	// extra check in case of integer downsizing
	if sourceBitsize > destBitsize {
		fmt.Printf(`			if sv > %s {
				return ConversionOverflow{%q, %q}
			}`+"\n", getMaxValue(destType), sourceType, destType)
	}

	// assign with conversion
	fmt.Printf(`			*dv = %s(sv)
		return nil`+"\n", destType)
}

func signedConvertGenerate(sourceType string, sourceBitsize int, destType string, destBitsize int) {
	fmt.Printf("		case *%s:\n", destType)

	// extra check in case of integer downsizing
	if sourceBitsize > destBitsize {
		fmt.Printf(`			if sv > %s || sv < %s {
				return ConversionOverflow{%q, %q}
			}`+"\n", getMaxValue(destType), getMinValue(destType), sourceType, destType)
	}

	// assign with conversion
	fmt.Printf(`			*dv = %s(sv)
			return nil`+"\n", destType)
}
