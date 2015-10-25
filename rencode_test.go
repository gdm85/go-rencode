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
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestFixedPosInts(t *testing.T) {
	for _, value := range []int8{10, -10} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeChar(t *testing.T) {
	for _, value := range []int8{100, -100} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestSingleByteArray(t *testing.T) {
	e := Encoder{}
	err := e.Encode([]byte{62})
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}

	value := found.([]byte)

	if value[0] != 62 {
		t.Fatalf("expected %v but %v found", 62, found)
	}
}

func TestDecodeShort(t *testing.T) {
	for _, value := range []int16{27123, -27123} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeInt(t *testing.T) {
	for _, value := range []int32{7483648, -7483648} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeLongLong(t *testing.T) {
	for _, value := range []int64{8223372036854775808, -8223372036854775808} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeBigNumber(t *testing.T) {
	var value big.Int

	value.SetUint64(^uint64(0))

	value.Mul(&value, big.NewInt(32))

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	i := found.(big.Int)

	if i.Cmp(&value) != 0 {
		t.Fatalf("expected %v but %v found", value, found)
	}
}

func TestDecodeFloat32(t *testing.T) {
	value := float32(1234.56)

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(float32)

	if value != f {
		t.Fatalf("expected %v but %v found", value, found)
	}
}

func TestDecodeFloat64(t *testing.T) {
	value := float64(1234.56)

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(float64)

	if value != f {
		t.Fatalf("expected %v but %v found", value, found)
	}
}

func TestDecodeFixedString(t *testing.T) {
	value := "foobarbaz"

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := string(found.([]byte))

	if value != f {
		t.Fatalf("expected %v but %v found", []byte(value), []byte(f))
	}
}

func TestDecodeString(t *testing.T) {
	value := strings.Repeat("f", 255)

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := string(found.([]byte))

	if value != f {
		t.Fatalf("expected %v but %v found", []byte(value), []byte(f))
	}
}

func TestDecodeUnicode(t *testing.T) {
	value := "fööbar"

	e := Encoder{}

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := string(found.([]byte))

	if value != f {
		t.Fatalf("expected %v but %v found", []byte(value), []byte(f))
	}
}

func TestDecodeNone(t *testing.T) {
	e := Encoder{}

	err := e.Encode(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}

	if nil != found {
		t.Fatalf("expected %v but %v found", nil, found)
	}
}

func TestDecodeBool(t *testing.T) {
	for _, value := range []bool{true, false} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}

		if found != value {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeStringBytes(t *testing.T) {
	for _, value := range [][]byte{
		[]byte{202, 132, 100, 114, 97, 119, 1, 0, 0, 63, 1, 242, 63},
		[]byte{202, 132, 100, 114, 97, 119, 1, 0, 0, 63, 1, 242, 63, 1, 60, 132, 120, 50, 54, 52, 49, 51, 48, 58, 0, 0, 0, 1, 65, 154, 35, 215, 48, 204, 4, 35, 242, 3, 122, 218, 67, 192, 127, 40, 241, 127, 2, 86, 240, 63, 135, 177, 23, 119, 63, 31, 226, 248, 19, 13, 192, 111, 74, 126, 2, 15, 240, 31, 239, 48, 85, 238, 159, 155, 197, 241, 23, 119, 63, 2, 23, 245, 63, 24, 240, 86, 36, 176, 15, 187, 185, 248, 242, 255, 0, 126, 123, 141, 206, 60, 188, 1, 27, 254, 141, 169, 132, 93, 220, 252, 121, 184, 8, 31, 224, 63, 244, 226, 75, 224, 119, 135, 229, 248, 3, 243, 248, 220, 227, 203, 193, 3, 224, 127, 47, 134, 59, 5, 99, 249, 254, 35, 196, 127, 17, 252, 71, 136, 254, 35, 196, 112, 4, 177, 3, 63, 5, 220},
	} {
		e := Encoder{}

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(hex.Dump(e.Bytes()))

		d := NewDecoder(bytes.NewReader(e.Bytes()))

		found, err := d.DecodeNext()
		if err != nil {
			t.Fatal(err)
		}
		f := found.([]byte)

		if bytes.Compare(value, f) != 0 {
			t.Fatalf("expected %v but %v found", value, found)
		}
	}
}

func TestDecodeFixedList(t *testing.T) {
	var l List

	l.Add(int8(100))
	l.Add(false)
	l.Add([]byte("foobar"))
	l.Add([]byte("bäz"))

	e := Encoder{}

	err := e.Encode(l)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(List)

	for i, v := range l.Values() {
		fv, err := f.Get(i)
		if err != nil {
			t.Fatal(err)
		}
		switch v.(type) {
		case []byte:
			if bytes.Compare(v.([]byte), fv.([]byte)) != 0 {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		default:
			if v != fv {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		}
	}
}

func TestDecodeList(t *testing.T) {
	var l List

	for i := 0; i < 80; i++ {
		l.Add(int8(100))
		l.Add(false)
		l.Add([]byte("foobar"))
		l.Add([]byte("bäz"))
	}

	e := Encoder{}

	err := e.Encode(l)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(List)

	for i, v := range l.Values() {
		fv, err := f.Get(i)
		if err != nil {
			t.Fatal(err)
		}
		switch v.(type) {
		case []byte:
			if bytes.Compare(v.([]byte), fv.([]byte)) != 0 {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		default:
			if v != fv {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		}
	}
}

func TestDecodeFixedDict(t *testing.T) {
	var dict Dictionary

	dict.Add("abcdefghijk", int16(1234))
	dict.Add(false, []byte("bäz"))

	e := Encoder{}

	err := e.Encode(dict)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(Dictionary)

	keys := dict.Keys()
	for i, v := range dict.Values() {
		fv, err := f.Get(keys[i])
		if err != nil {
			t.Fatal(err)
		}
		switch v.(type) {
		case []byte:
			if bytes.Compare(v.([]byte), fv.([]byte)) != 0 {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		default:
			if v != fv {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		}
	}
}

func TestDecodeDictionary(t *testing.T) {
	var dict Dictionary
	var nestedDict Dictionary
	var nestedList List

	nestedDict.Add("abcdefghijk", int16(1234))
	nestedDict.Add(false, []byte("bäz"))
	nestedList.Add(true)
	nestedList.Add("carrot")

	for i := 0; i < 120; i++ {
		dict.Add(fmt.Sprintf("abcde %d", i), []byte("foo"))
		dict.Add(fmt.Sprintf("fghijk %d", i), nestedDict)
		dict.Add(fmt.Sprintf("z %d", i), nestedList)
	}

	e := Encoder{}

	err := e.Encode(dict)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.Dump(e.Bytes()))

	d := NewDecoder(bytes.NewReader(e.Bytes()))

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(Dictionary)

	keys := dict.Keys()
	for i, v := range dict.Values() {
		fv, err := f.Get(keys[i])
		if err != nil {
			t.Fatal(err)
		}
		switch v.(type) {
		case []byte:
			if bytes.Compare(v.([]byte), fv.([]byte)) != 0 {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		case Dictionary:
			d1 := v.(Dictionary)
			d2 := fv.(Dictionary)
			if !d1.Compare(&d2) {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		case List:
			l1 := v.(List)
			l2 := fv.(List)
			if !l1.Compare(&l2) {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		default:
			if v != fv {
				t.Fatalf("index %d: expected %q (type %T) but %q (type %T) found", i, v, v, fv, fv)
			}
		}
	}

	// check that we have a carrot in one of the nested list values
	v, err := f.Get("z 10")
	if err != nil {
		t.Fatal(err)
	}

	l := v.(List)

	fv, err := l.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	if string(fv.([]byte)) != "carrot" {
		t.Fatal("carrot not found")
	}
}
