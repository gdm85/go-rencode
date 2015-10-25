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

//go:generate sh -c "go run --tags=generate generate.go > rencode_generated.go"

import (
	"bytes"
)

const (
    DEFAULT_FLOAT_BITS = 32 // Default number of bits for serialized floats, either 32 or 64 (also a parameter for dumps()).
    MAX_INT_LENGTH = 64 // Maximum length of integer when written as base 10 string.
    // The bencode 'typecodes' such as i, d, etc have been extended and relocated on the base-256 character set.
    CHR_LIST    = 59
    CHR_DICT    = 60
    CHR_INT     = 61
    CHR_INT1    = 62
    CHR_INT2    = 63
    CHR_INT4    = 64
    CHR_INT8    = 65
    CHR_FLOAT32 = 66
    CHR_FLOAT64 = 44
    CHR_TRUE    = 67
    CHR_FALSE   = 68
    CHR_NONE    = 69
    CHR_TERM    = 127
    // Positive integers with value embedded in typecode.
    INT_POS_FIXED_START = 0
    INT_POS_FIXED_COUNT = 44
    // Dictionaries with length embedded in typecode.
    DICT_FIXED_START = 102
    DICT_FIXED_COUNT = 25
    // Negative integers with value embedded in typecode.
    INT_NEG_FIXED_START = 70
    INT_NEG_FIXED_COUNT = 32
    // Strings with length embedded in typecode.
    STR_FIXED_START = 128
    STR_FIXED_COUNT = 64
    // Lists with length embedded in typecode.
    LIST_FIXED_START = STR_FIXED_START+STR_FIXED_COUNT
    LIST_FIXED_COUNT = 64
)

type Rencode struct {
	buffer bytes.Buffer
}

func (r *Rencode) encodeInt8(x int8) error {
	if x >=0 && x < INT_POS_FIXED_COUNT {
		_, err := r.buffer.Write([]byte{byte(INT_POS_FIXED_START + x)})
		return err
	}
	if x < 0 && x >= -INT_NEG_FIXED_COUNT {
		_, err := r.buffer.Write([]byte{byte(INT_NEG_FIXED_START - 1 - x)})
		return err
	}
	if x >= -128 && x <= 127 {
		_, err := r.buffer.Write([]byte{CHR_INT1, byte(x)})
		return err
	}
	panic("impossible just happened")
}

/*cdef encode_char(char **buf, unsigned int *pos, signed char x):
    if 0 <= x < INT_POS_FIXED_COUNT:
        write_buffer_char(buf, pos, INT_POS_FIXED_START + x)
    elif -INT_NEG_FIXED_COUNT <= x < 0:
        write_buffer_char(buf, pos, INT_NEG_FIXED_START - 1 - x)
    elif -128 <= x < 128:
        write_buffer_char(buf, pos, CHR_INT1)
        write_buffer_char(buf, pos, x)*/
        




/*
cdef encode(char **buf, unsigned int *pos, data):
    t = type(data)
    if t == int or t == long:
        if -128 <= data < 128:
            encode_char(buf, pos, data)
        elif -32768 <= data < 32768:
            encode_short(buf, pos, data)
        elif MIN_SIGNED_INT <= data < MAX_SIGNED_INT:
            encode_int(buf, pos, data)
        elif MIN_SIGNED_LONGLONG <= data < MAX_SIGNED_LONGLONG:
            encode_long_long(buf, pos, data)
        else:
            s = str(data)
            if py3:
                s = s.encode("ascii")
            if len(s) >= MAX_INT_LENGTH:
                raise ValueError("Number is longer than %d characters" % MAX_INT_LENGTH)
            encode_big_number(buf, pos, s)
    elif t == float:
        if _float_bits == 32:
            encode_float32(buf, pos, data)
        elif _float_bits == 64:
            encode_float64(buf, pos, data)
        else:
            raise ValueError('Float bits (%d) is not 32 or 64' % _float_bits)

    elif t == bytes:
        encode_str(buf, pos, data)

    elif t == unicode:
        u = data.encode("utf8")
        encode_str(buf, pos, u)

    elif t == type(None):
        encode_none(buf, pos)

    elif t == bool:
        encode_bool(buf, pos, data)

    elif t == list or t == tuple:
        encode_list(buf, pos, data)

    elif t == dict:
        encode_dict(buf, pos, data)

    else:
        raise Exception("type %s not handled" % t)
*/
