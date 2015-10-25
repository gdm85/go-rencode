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

const top = `package rencode

func (r *Rencode) encode(data interface{}) error {
	switch data.(type) {
		case int8:
			return r.encodeChar(data.(int8))`

/*		case int:
			if -128 <= data.(int) && data.(int) < 128 {
				return r.encodeChar(int8(data.(int)))
			}
		case int16:
		case int32:
		case int64:
		case uint:
		case uint8:
		case uint16:
		case uint32:
		case uint64:
			
		default:
	}
}
* */

var sigintTypes = []string{"int", "int16", "int32", "int64"}

func main() {
	fmt.Println(top)
	
	for _, t := range sigintTypes {
		fmt.Printf("\tcase %s:\n", t)
		fmt.Printf(`			if -128 <= data.(%s) && data.(%s) < 128 {
				return r.encodeInt8(int8(data.(%s)))
			}` + "\n", t, t, t)
	}
	
	fmt.Println(`default:
		return fmt.Errorf("could not encode data of type %T", data)
	}
}`)
}
