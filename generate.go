// +build generate

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
package main

import (
	"fmt"
)

// template block starts
const top = `package rencode

import (
	"fmt"
	"math"
	"math/big"
)

func (r *Encoder) Encode(data interface{}) error {
	if data == nil {
		return r.EncodeNone()
	}
	switch data.(type) {
		case big.Int:
			x := data.(big.Int)
			s := x.String()
			if len(s) > MAX_INT_LENGTH {
				return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
			}
			return r.EncodeBigNumber(s)
		case bool:
			return r.EncodeBool(data.(bool))
		case float32:
			return r.EncodeFloat32(data.(float32))
		case float64:
			return r.EncodeFloat64(data.(float64))
		case []byte:
			return r.EncodeBytes(data.([]byte))
		case string:
			// all strings will be treated as byte arrays
			return r.EncodeBytes([]byte(data.(string)))
		case int8:
			return r.EncodeInt8(data.(int8))`

// template block ends

var (
	intTypes           map[string]int
	supportedListTypes = []string{"bool", "float32", "float64", "string", "int8"}
)

func init() {
	intTypes = map[string]int{"uint8": 8, "uint16": 16, "int16": 15, "uint32": 32, "int32": 31, "int64": 63} // NOTE: uint64 is not supported as it can overflow int64

	if ^uint(0) == uint(^uint32(0)) {
		intTypes["uint"] = 32
		intTypes["int"] = 31
	} else if ^uint(0) == uint(^uint64(0)) {
		// same here, 'uint' is not being defined on purpose
		intTypes["int"] = 63

		// add now to supported list types, since it's not referenced in intTypes map
		supportedListTypes = append(supportedListTypes, "uint")
	} else {
		panic("unrecognized default uint bitsize")
	}

	// add integer types to supported list types
	for k, _ := range intTypes {
		// array of uint8 == array of byte, do not support as lists
		if k == "uint8" {
			continue
		}
		supportedListTypes = append(supportedListTypes, k)
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
		fmt.Printf(`	case %s:
		x := data.(%s)`+"\n", t, t)

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
	fmt.Printf("\t\tcase %s:\n", caseStr)
	fmt.Println(`			s := fmt.Sprintf("%d", data)
			if len(s) > MAX_INT_LENGTH {
				return fmt.Errorf("Number is longer than %d characters", MAX_INT_LENGTH)
			}
			return r.EncodeBigNumber(s)`)

	// support lists of all supported types
	for _, t := range supportedListTypes {
		fmt.Printf(`		case []%s:
			x := data.([]%s)
			if len(x) < LIST_FIXED_COUNT {
				_, err := r.buffer.Write([]byte{byte(LIST_FIXED_START + len(x))})
				if err != nil {
					return err
				}
				for _, v := range x {
					err = r.Encode(v)
					if err != nil {
						return err
					}
				}
				return nil
			}
			_, err := r.buffer.Write([]byte{byte(CHR_LIST)})
			if err != nil {
				return err
			}

			for _, v := range x {
				err = r.Encode(v)
				if err != nil {
					return err
				}
			}

			_, err = r.buffer.Write([]byte{byte(CHR_TERM)})
			return err`+"\n", t, t)
	}

	// re-add byte to supported map types
	supportedListTypes = append(supportedListTypes, "byte")

	if 1 == 2 {
		// support maps of all supported types
		for _, keyType := range supportedListTypes {
			for _, valueType := range supportedListTypes {
				for _, listOrNot := range []string{"", "[]"} {
					mapType := fmt.Sprintf("map[%s]%s%s", keyType, listOrNot, valueType)
					// generate the map type matching case
					fmt.Printf(`		case %s:
						x := data.(%s)
						if len(x) < DICT_FIXED_COUNT {
							_, err := r.buffer.Write([]byte{byte(DICT_FIXED_START + len(x))})
							if err != nil {
								return err
							}
							for k, v := range x {
								err = r.Encode(k)
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
						_, err := r.buffer.Write([]byte{byte(CHR_DICT)})
						if err != nil {
							return err
						}

						for k, v := range x {
							err = r.Encode(k)
							if err != nil {
								return err
							}
							err = r.Encode(v)
							if err != nil {
								return err
							}
						}

						_, err = r.buffer.Write([]byte{byte(CHR_TERM)})
						return err`+"\n", mapType, mapType)
				}
			}
		}
	} // temporarily disable maps because of huge source (13k)

	// tail default case
	fmt.Println(`default:
		return fmt.Errorf("could not encode data of type %T", data)
	}
	panic("unexpected fallthrough")
}`)
}
