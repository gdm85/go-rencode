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
	"bytes"
	"testing"
)

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

func TestListScanAndShift(t *testing.T) {
	l1 := NewList(100, true, "hello world", NewList(42, "nesting is awesome"), 3.14, Dictionary{})

	var (
		i  int
		b  bool
		s  string
		l2 List
		f  float64
		d  Dictionary
	)

	err := l1.Scan(&i, &b, &s, &l2, &f, &d)
	if err != nil {
		t.Fatal(err)
	}

	if i != 100 || !b || s != "hello world" {
		t.Error("invalid data scanned")
	}

	var s2 string
	shifted := l1.Shift(2)
	if shifted != 2 {
		t.Fatalf("could not shift %d items, shifted %d instead", 2, shifted)
	}

	err = l1.Scan(&s2)
	if err != nil {
		t.Fatal(err)
	}

	if s2 != "hello world" {
		t.Fatal("could not read third value of the list after shift")
	}
}
