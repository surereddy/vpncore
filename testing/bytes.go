/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package testing

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"bytes"
)

func PrintBytes(b []byte, n int, title string) {
	fmt.Printf("%s: %v\n", title, b[:n])
}


func ProtoMessageEqual(msg1 proto.Message, msg2 proto.Message) bool {
	// Now test and newTest contain the same data.
	b1, err := proto.Marshal(msg1)
	if err != nil {
		return false
	}

	b2, err := proto.Marshal(msg2)
	if err != nil {
		return false
	}

	if !bytes.Equal(b1, b2) {
		fmt.Printf("b1:%v\n, b2:%v\n", b1, b2)
	}

	return bytes.Equal(b1, b2)
}