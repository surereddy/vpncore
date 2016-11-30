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

package noise

import (
	"github.com/FTwOoO/vpncore/conn"
	"github.com/flynn/noise"
)

type noiseIKMsg struct {
	base conn.Message
	ctx  *NoiseIKMessageContext
}

func (this *noiseIKMsg) Unmarshal(packet []byte) (err error) {

	if !this.ctx.IsHandshakeCompleted() {
		this.ctx.StepOne()

		if this.ctx.IsReadStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1, err = this.ctx.Hs.ReadMessage(nil, packet)
			if err != nil {
				return
			}

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.ctx.IsFinalStel() {
				this.ctx.HandshakeDone(cs0, cs1)
			}

			_, err = this.base.Unmarshal(packet)
			return

		} else {
			return conn.ErrInValidHandshakeStep
		}

	}

	packet = this.ctx.CS0.Encrypt(nil, nil, packet)
	_, err = this.base.Unmarshal(packet)
	return
}

func (this *noiseIKMsg) Marshal() (packet []byte, err error) {

	if !this.ctx.IsHandshakeCompleted() {
		this.ctx.StepOne()

		if this.ctx.IsWriteStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1 = this.ctx.Hs.WriteMessage(nil, packet)

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.ctx.IsFinalStel() {
				this.ctx.HandshakeDone(cs0, cs1)
			}

		} else {
			return nil, conn.ErrInValidHandshakeStep
		}

	} else {
		packet = this.ctx.CS0.Encrypt(nil, nil, packet)
	}

	packet, err = this.base.Marshal()
	return
}
