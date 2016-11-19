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

package crypt

import (
	"net"
	"github.com/FTwOoO/vpncore/enc"
)

type cryptConn struct {
	net.Conn
	R *enc.CryptionReadWriter
}

func NewCryptConn(conn net.Conn, stream enc.StreamCipher) (*cryptConn, error) {
	connection := new(cryptConn)
	connection.Conn = conn
	connection.R = enc.NewCryptionReadWriter(stream, conn)

	return connection, nil
}

func (c *cryptConn) Read(b []byte) (n int, err error) {
	return c.R.Read(b)
}

func (c *cryptConn) Write(b []byte) (n int, err error) {
	return c.R.Write(b)
}
