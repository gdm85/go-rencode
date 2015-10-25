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
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

// Decoder implements a rencode decoder
type Decoder struct {
	r io.Reader
}

var (
	maxUint64 big.Int
)

func init() {
	maxUint64.SetUint64(^uint64(0))
}

// NewDecoder returns a rencode decoder that sources all bytes from the specified reader
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (r *Decoder) readByte() (b byte, err error) {
	data := []byte{0}
	_, err = r.r.Read(data)
	if err != nil {
		return
	}
	b = data[0]
	return
}

func (r *Decoder) readSlice(delim byte) (data []byte, err error) {
	var b byte
	for {
		b, err = r.readByte()
		if err != nil {
			return
		}

		if b == delim {
			break
		}

		data = append(data, b)
	}
	return
}

// DecodeNext returns the next available object stored in the rencode stream.
// If no more objects are available, an io.EOF error will be returned.
func (r *Decoder) DecodeNext() (interface{}, error) {
	typeCode, err := r.readByte()
	if err != nil {
		return nil, err
	}

	return r.decode(typeCode)
}

func (r *Decoder) decode(typeCode byte) (v interface{}, err error) {
	switch typeCode {
	case CHR_TRUE:
		v = true
	case CHR_FALSE:
		v = false
	case CHR_NONE:
		// leave v as nil
	case CHR_INT1:
		data := []byte{0}
		_, err = r.r.Read(data)
		if err != nil {
			return
		}
		v = int8(data[0])
	case CHR_INT2:
		var data int16
		err = binary.Read(r.r, binary.BigEndian, &data)
		v = data
	case CHR_INT4:
		var data int32
		err = binary.Read(r.r, binary.BigEndian, &data)
		v = data
	case CHR_INT8:
		var data int64
		err = binary.Read(r.r, binary.BigEndian, &data)
		v = data
	case CHR_INT:
		var collected []byte
		collected, err = r.readSlice(CHR_TERM)
		if err != nil {
			return
		}

		var i big.Int
		_, err = fmt.Sscan(string(collected), &i)
		if err != nil {
			return
		}

		// if this is simply an uint64, return it as such
		if i.Cmp(&maxUint64) == 0 {
			v = maxUint64
		} else {
			v = i
		}
	case CHR_FLOAT32:
		var data float32
		err = binary.Read(r.r, binary.BigEndian, &data)
		v = data
	case CHR_FLOAT64:
		var data float64
		err = binary.Read(r.r, binary.BigEndian, &data)
		v = data
	case CHR_LIST:
		v, err = r.decodeList()
		return
	case CHR_DICT:
		v, err = r.decodeDict()
		return
	default:
		if INT_POS_FIXED_START <= typeCode && typeCode < INT_POS_FIXED_START+INT_POS_FIXED_COUNT {
			v = int8(typeCode) - INT_POS_FIXED_START
			return
		}
		if INT_NEG_FIXED_START <= typeCode && typeCode < INT_NEG_FIXED_START+INT_NEG_FIXED_COUNT {
			i := (int(typeCode) - INT_NEG_FIXED_START + 1) * -1
			v = int8(i)
			return
		}
		if STR_FIXED_START <= typeCode && typeCode < STR_FIXED_START+STR_FIXED_COUNT {
			b := typeCode - STR_FIXED_START
			data := make([]byte, b)

			_, err = r.r.Read(data)
			if err != nil {
				return
			}
			v = data
			return
		}
		if '1' <= typeCode && typeCode <= '9' {
			var collected []byte
			collected, err = r.readSlice(':')
			if err != nil {
				return
			}

			// use the typeCode as first digit
			n := []byte{typeCode}
			n = append(n, collected...)

			var stringSz int
			stringSz, err = strconv.Atoi(string(n))
			if err != nil {
				return
			}

			data := make([]byte, stringSz)
			_, err = r.r.Read(data)
			if err != nil {
				return
			}

			v = data
		}

		if LIST_FIXED_START <= typeCode && typeCode <= (LIST_FIXED_START+LIST_FIXED_COUNT-1) {
			var l List
			var value interface{}
			var i byte
			size := typeCode - LIST_FIXED_START

			for i = 0; i < size; i++ {
				// get next value
				value, err = r.DecodeNext()
				if err != nil {
					return
				}

				// add, never update existing key
				l.Add(value)
			}
			v = l
		}
		if DICT_FIXED_START <= typeCode && typeCode < DICT_FIXED_START+DICT_FIXED_COUNT {
			var d Dictionary
			var key, value interface{}
			var i byte
			size := typeCode - DICT_FIXED_START

			for i = 0; i < size; i++ {
				// get next key
				key, err = r.DecodeNext()
				if err != nil {
					return
				}

				// get next value
				value, err = r.DecodeNext()
				if err != nil {
					return
				}

				// add, never update existing key
				err = d.Add(key, value)
				if err != nil {
					return
				}
			}
			v = d
		}
	} // end of switch

	// AOK
	return
}

func (r *Decoder) decodeDict() (d Dictionary, err error) {
	var key, value interface{}
	var typeCode byte

	for {
		typeCode, err = r.readByte()
		if err != nil {
			return
		}
		if typeCode == CHR_TERM {
			// no more (key, value) pairs
			break
		}

		// get next key
		key, err = r.decode(typeCode)
		if err != nil {
			return
		}

		typeCode, err = r.readByte()
		if err != nil {
			return
		}
		if typeCode == CHR_TERM {
			err = fmt.Errorf("incomplete key-value pair in dictionary data")
			break
		}

		// get next value
		value, err = r.decode(typeCode)
		if err != nil {
			return
		}

		// add, never update existing key
		err = d.Add(key, value)
		if err != nil {
			return
		}
	}

	return
}

func (r *Decoder) decodeList() (l List, err error) {
	var value interface{}
	var typeCode byte

	for {
		typeCode, err = r.readByte()
		if err != nil {
			return
		}
		if typeCode == CHR_TERM {
			// no more values
			break
		}

		// get next value
		value, err = r.decode(typeCode)
		if err != nil {
			return
		}

		l.Add(value)
	}

	return
}
