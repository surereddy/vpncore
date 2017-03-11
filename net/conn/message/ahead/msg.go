/*
 * Author: FTwOoO <booobooob@gmail.com>
 * Created: 2017-03
 */

package ahead

import (
)

//go:generate msgp

type Foo struct {
    Bar string  `msg:"bar"`
    Baz float64 `msg:"baz"`
}
