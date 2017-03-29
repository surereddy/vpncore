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

package encryption

import (
	"github.com/FTwOoO/vpncore/crypto"
	"github.com/FTwOoO/vpncore/net/conn"
	"errors"
)

var _ conn.StreamContext = new(CryptStreamContext)


type CryptStreamContext struct {
	*crypto.EncrytionConfig
}

func (this *CryptStreamContext) Layer() conn.Layer {
	return conn.ENCRYPTION_LAYER
}


//must called after init object
func (this *CryptStreamContext) Valid() (bool, error) {
	if this.EncrytionConfig == nil {
		return false, errors.New("Need to set up encrytion config")
	}

	return true, nil
}

func (this *CryptStreamContext) Pipe(base conn.StreamIO) (c conn.StreamIO) {
	cipher, err := crypto.NewStreamCipher(this.EncrytionConfig)
	if err != nil {
		return nil
	}

	c, err = crypto.NewCryptionReadWriter(base, cipher)
	if err != nil {
		return nil
	}

	return

}
