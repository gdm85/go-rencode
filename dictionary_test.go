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

package rencode

import "testing"

func TestToSnakeCase(t *testing.T) {
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
