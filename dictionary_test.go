//
// go-rencode v0.1.8 - Go implementation of rencode - fast (basic)
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

import "testing"

func TestToSnakeCase(t *testing.T) {
	t.Parallel()

	table := []struct {
		Input  string
		Output string
	}{
		{"", ""},
		{"small", "small"},
		{"TestCase", "test_case"},
		{"testCase", "test_case"},
		{"ETA", "eta"},
		{"JSONData", "json_data"},
		{"entityID", "entity_id"},
		{"AAArgh", "aa_argh"},
		{"zZ", "z_z"},
		{"already_converted", "already_converted"},
	}

	for _, testCase := range table {
		value := ToSnakeCase(testCase.Input)

		if value != testCase.Output {
			t.Errorf(
				"For input '%s' got '%s' (expected '%s')",
				testCase.Input, value, testCase.Output)
		}
	}
}

func TestToStruct(t *testing.T) {
	t.Parallel()

	var s struct {
		Alpha int
		Beta  string
		Gamma uint8
	}
	var d Dictionary
	d.Add("alpha", int(54123))
	d.Add("beta", "test")
	d.Add("gamma", uint8(42))

	err := d.ToStruct(&s, "")
	if err != nil {
		t.Errorf("expected succcess but got %v", err)
	}
}

func TestExtraFieldsFailure(t *testing.T) {
	t.Parallel()

	var s struct {
		Alpha int
		Beta  string
	}
	var d Dictionary
	d.Add("alpha", int(54123))
	d.Add("beta", "test")
	d.Add("gamma", uint8(42))

	err := d.ToStruct(&s, "")
	if err == nil {
		t.Error("expected failure")
	}
}

func TestExcludeTag(t *testing.T) {
	t.Parallel()

	var s struct {
		Alpha int
		Beta  string
		Gamma float64 `rencode:"exclude-me"`
	}
	var d Dictionary
	d.Add("alpha", int(54123))
	d.Add("beta", "test")

	err := d.ToStruct(&s, "exclude-me")
	if err != nil {
		t.Errorf("mapping failed: %v", err)
	}
}

func TestNestedExcludeTag(t *testing.T) {
	t.Parallel()

	var s struct {
		Alpha int
		Beta  string
		Gamma float64 `rencode:"exclude-me"`
		Delta []struct {
			Epsilon bool
			Zeta    int8 `rencode:"exclude-me"`
		}
	}

	var d2 Dictionary
	d2.Add("epsilon", true)

	var l List
	l.Add(d2)

	var d Dictionary
	d.Add("alpha", int(54123))
	d.Add("beta", "test")
	d.Add("delta", l)

	err := d.ToStruct(&s, "exclude-me")
	if err != nil {
		t.Errorf("mapping failed: %v", err)
	}
}
