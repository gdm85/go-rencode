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

package rencode

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestFixedPosInts(t *testing.T) {
	t.Parallel()

	for _, value := range []int8{10, -10} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	for _, value := range []int8{100, -100} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	var b bytes.Buffer
	e := NewEncoder(&b)
	err := e.Encode([]byte{62})
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	for _, value := range []int16{27123, -27123} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	for _, value := range []int32{7483648, -7483648} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	for _, value := range []int64{8223372036854775808, -8223372036854775808} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	var value big.Int

	value.SetUint64(^uint64(0))

	value.Mul(&value, big.NewInt(32))

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	value := float32(1234.56)

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	value := float64(1234.56)

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	value := "foobarbaz"

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	value := strings.Repeat("f", 255)

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	value := "fööbar"

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(value)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

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
	t.Parallel()

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(nil)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}

	if nil != found {
		t.Fatalf("expected %v but %v found", nil, found)
	}
}

func TestEncodeNilInterface(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer
	e := NewEncoder(&b)

	type someInterface interface{
		SomeMethod()
	}

	var v someInterface

	err := e.Encode(v)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}

	if nil != found {
		t.Fatalf("expected %v but %v found", nil, found)
	}
}

func TestDecodeBool(t *testing.T) {
	t.Parallel()

	for _, value := range []bool{true, false} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	for _, value := range [][]byte{
		{202, 132, 100, 114, 97, 119, 1, 0, 0, 63, 1, 242, 63},
		{202, 132, 100, 114, 97, 119, 1, 0, 0, 63, 1, 242, 63, 1, 60, 132, 120, 50, 54, 52, 49, 51, 48, 58, 0, 0, 0, 1, 65, 154, 35, 215, 48, 204, 4, 35, 242, 3, 122, 218, 67, 192, 127, 40, 241, 127, 2, 86, 240, 63, 135, 177, 23, 119, 63, 31, 226, 248, 19, 13, 192, 111, 74, 126, 2, 15, 240, 31, 239, 48, 85, 238, 159, 155, 197, 241, 23, 119, 63, 2, 23, 245, 63, 24, 240, 86, 36, 176, 15, 187, 185, 248, 242, 255, 0, 126, 123, 141, 206, 60, 188, 1, 27, 254, 141, 169, 132, 93, 220, 252, 121, 184, 8, 31, 224, 63, 244, 226, 75, 224, 119, 135, 229, 248, 3, 243, 248, 220, 227, 203, 193, 3, 224, 127, 47, 134, 59, 5, 99, 249, 254, 35, 196, 127, 17, 252, 71, 136, 254, 35, 196, 112, 4, 177, 3, 63, 5, 220},
	} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

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
	t.Parallel()

	var l List

	l.Add(int8(100), false, []byte("foobar"), []byte("bäz"))

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(l)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(List)

	listCompareVerbose(t, &l, &f)
}

func TestDecodeList(t *testing.T) {
	t.Parallel()

	var l List

	for i := 0; i < 80; i++ {
		l.Add(int8(100), false, []byte("foobar"), []byte("bäz"))
	}

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(l)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(List)

	listCompareVerbose(t, &l, &f)
}

func TestDecodeFixedDict(t *testing.T) {
	t.Parallel()

	var dict Dictionary

	dict.Add("abcdefghijk", int16(1234))
	dict.Add(false, []byte("bäz"))
	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(dict)
	if err != nil {
		t.Fatal(err)
	}

	// start decoding
	d := NewDecoder(&b)

	// a dictionary is expected
	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	decodedDict := found.(Dictionary)

	dictCompareVerbose(t, &dict, &decodedDict)
}

func TestDecodeDictionary(t *testing.T) {
	t.Parallel()

	var dict Dictionary
	var nestedDict Dictionary
	var nestedList List

	nestedDict.Add("abcdefghijk", int16(1234))
	nestedDict.Add(false, []byte("bäz"))
	nestedList.Add(true, "carrot")

	for i := 0; i < 120; i++ {
		dict.Add(fmt.Sprintf("abcde %d", i), []byte("foo"))
		dict.Add(fmt.Sprintf("fghijk %d", i), nestedDict)
		dict.Add(fmt.Sprintf("z %d", i), nestedList)
	}

	var b bytes.Buffer
	e := NewEncoder(&b)

	err := e.Encode(dict)
	if err != nil {
		t.Fatal(err)
	}

	d := NewDecoder(&b)

	found, err := d.DecodeNext()
	if err != nil {
		t.Fatal(err)
	}
	f := found.(Dictionary)

	if !dictCompareVerbose(t, &dict, &f) {
		return
	}

	// check that we have a carrot in one of the nested list values
	v, ok := f.Get("z 10")
	if !ok {
		t.Fatal("key not found")
	}

	l := v.(List)

	fv := l.Values()[1]
	if string(fv.([]byte)) != "carrot" {
		t.Fatal("carrot not found")
	}
}

func TestDecodeIntIntoFloat(t *testing.T) {
	t.Parallel()

	for _, value := range []interface{}{int8(45), int32(7483648), int32(-7483648)} {
		var b bytes.Buffer
		e := NewEncoder(&b)

		err := e.Encode(value)
		if err != nil {
			t.Fatal(err)
		}

		d := NewDecoder(&b)

		var f32 float32
		err = d.Scan(&f32)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func listCompareVerbose(t *testing.T, a, b *List) bool {
	if a.Length() != b.Length() {
		t.Errorf("list length mismatch: %v != %v", a.Length(), b.Length())
		return false
	}

	matching := true
	for i, aV := range a.Values() {
		bV := b.Values()[i]

		// normalize both values to string if they are of []byte type
		if v, ok := aV.([]byte); ok {
			aV = string(v)
		}
		if v, ok := bV.([]byte); ok {
			bV = string(v)
		}

		if aV != bV {
			t.Errorf("index %d: expected %q (type %T) but %q (type %T) found", i, aV, aV, bV, bV)
			matching = false
		}
	}

	return matching
}

func dictCompareVerbose(t *testing.T, a, b *Dictionary) bool {
	if a.Length() != b.Length() {
		t.Errorf("dictionary length mismatch: %v != %v", a.Length(), b.Length())
		return false
	}

	matching := true
	for _, k := range a.Keys() {
		// get value on both dictionaries
		aV, ok := a.Get(k)
		if !ok {
			t.Errorf("value with key %v not found on first dictionary", k)
			return false
		}
		bV, ok := b.Get(k)
		if !ok {
			t.Errorf("value with key %v not found on second dictionary", k)
			return false
		}

		// normalize both values to string if they are of []byte type
		if v, ok := aV.([]byte); ok {
			aV = string(v)
		}
		if v, ok := bV.([]byte); ok {
			bV = string(v)
		}

		switch v := aV.(type) {
		case Dictionary:
			d2 := bV.(Dictionary)
			if !dictCompareVerbose(t, &v, &d2) {
				matching = false
			}
		case List:
			l2 := bV.(List)
			if !listCompareVerbose(t, &v, &l2) {
				matching = false
			}
		default:
			if aV != bV {
				t.Fatalf("index %q: expected %v (type %T) but %v (type %T) found", k, v, v, bV, bV)
			}
		}
	}

	return matching
}
