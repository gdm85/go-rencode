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
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

type Decoder struct {
	r io.Reader
}

var (
	maxUint64 big.Int
)

func init() {
	maxUint64.SetUint64(^uint64(0))
}

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

func (r *Decoder) DecodeNext() (v interface{}, err error) {
	var typeCode byte
	typeCode, err = r.readByte()
	if err != nil {
		return
	}

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

		i := new(big.Int)
		_, err = fmt.Sscan(string(collected), i)
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
		panic("list decoding not yet implemented")
		//v, err = r.decodeList()
		return
	case CHR_DICT:
		panic("dict decoding not yet implemented")
		//v, err = r.decodeDict()
		return
	default:
		if INT_POS_FIXED_START <= typeCode && typeCode < INT_POS_FIXED_START+INT_POS_FIXED_COUNT {
			var b byte
			b, err = r.readByte()
			if err != nil {
				return
			}

			v = int8(b) - INT_POS_FIXED_START
			return
		}
		if INT_NEG_FIXED_START <= typeCode && typeCode < INT_NEG_FIXED_START+INT_NEG_FIXED_COUNT {
			var b byte
			b, err = r.readByte()
			if err != nil {
				return
			}

			i := (int(b) - INT_NEG_FIXED_START + 1) * -1
			v = int8(i)
			return
		}
		if STR_FIXED_START <= typeCode && typeCode < STR_FIXED_START+STR_FIXED_COUNT {
			var b byte
			b, err = r.readByte()
			if err != nil {
				return
			}

			b = b - STR_FIXED_START + 1
			data := make([]byte, b)

			_, err = r.r.Read(data)
			if err != nil {
				return
			}
			v = string(data)
			return
		}
		if 49 <= typeCode && typeCode <= 57 {
			var collected []byte
			collected, err = r.readSlice(':')
			if err != nil {
				return
			}

			var stringSz int
			stringSz, err = strconv.Atoi(string(collected))
			if err != nil {
				return
			}

			data := make([]byte, stringSz)
			_, err = r.r.Read(data)
			if err != nil {
				return
			}

			v = string(data)
		}

		if LIST_FIXED_START <= typeCode && typeCode <= (LIST_FIXED_START+LIST_FIXED_COUNT-1) {
			panic("list decoding not yet implemented")
		}
		if DICT_FIXED_START <= typeCode && typeCode < DICT_FIXED_START+DICT_FIXED_COUNT {
			panic("map decoding not yet implemented")
		}
	} // end of switch

	// AOK
	err = nil
	return
}

/*func (r *Decoder) PushOffset() error {
	offset, err := r.r.Seek(1, 0)
	if err != nil {
		return err
	}
	r.offsets = append(r.offsets, offset)
	return nil
}

func (r *Decoder) PopOffset() error {
	if len(r.offsets) == 0 {
		return fmt.Errorf("no more pushed offsets available")
	}
	offset = r.offsets[len(r.offsets)-1]
	r.offsets = r.offsets[:len(r.offsets)-1]

	_, err := r.r.Seek(0, offset)
	return err
}*/

/*func (r *Decoder) DecodeInt8(v *int8) error {
	typeCode := []byte{0}
	_, err := r.r.Read(typeCode)
	if err != nil {
		return err
	}

	if typeCode != CHR_INT1 {
		return fmt.Errorf("expected type code %d but %d found instead", CHR_INT1, typeCode)
	}

	data := []byte{0}
	_, err = r.r.Read(data)
	if err != nil {
		return err
	}
	*v = int8(data[0])

	return nil
}*/
