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

package noiseik

import (
	"github.com/FTwOoO/vpncore/conn"
	"github.com/FTwOoO/noise"
)

type NoiseIKMessageContext struct {
	Hs            *noise.HandshakeState
	IsInitiator   bool
	Pattern       noise.HandshakePattern

	CSEnc         *noise.CipherState
	CSDec         *noise.CipherState
	handshakeStep int
}

func NewNoiseIKMessageContext(cs noise.CipherSuite, pg []byte, staticI noise.DHKey, staticR noise.DHKey, isInitiator bool) (*NoiseIKMessageContext, error) {

	this := new(NoiseIKMessageContext)
	if cs == nil {
		cs = noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
	}

	//staticI := cs.GenerateKeypair(nil)
	//staticR := cs.GenerateKeypair(nil)

	if isInitiator {
		this.Hs = noise.NewHandshakeState(noise.Config{
			CipherSuite: cs,
			Pattern: noise.HandshakeIK,
			Initiator: true,
			Prologue: pg,
			StaticKeypair: staticI,
			PeerStatic: staticR.Public},
		)
	} else {

		this.Hs = noise.NewHandshakeState(noise.Config{
			CipherSuite: cs,
			Pattern: noise.HandshakeIK,
			Initiator:false,
			Prologue: pg,
			StaticKeypair: staticR},
		)
	}

	this.IsInitiator = isInitiator
	this.Pattern = noise.HandshakeIK
	return this, nil
}

func (this *NoiseIKMessageContext) Valid() (bool, error) {
	return true, nil
}

func (this *NoiseIKMessageContext) Layer() conn.Layer {
	return conn.AUTH_LAYER
}

func (this *NoiseIKMessageContext) IsHandshakeCompleted() bool {
	return this.handshakeStep >= len(this.Pattern.Messages)
}

func (this *NoiseIKMessageContext) IsWriteStep() bool {
	return (this.handshakeStep % 2 == 1 && this.IsInitiator) || (this.handshakeStep % 2 == 0 && !this.IsInitiator)
}

func (this *NoiseIKMessageContext) IsReadStep() bool {
	return (this.handshakeStep % 2 == 1 && !this.IsInitiator) || (this.handshakeStep % 2 == 0 && this.IsInitiator)
}

func (this *NoiseIKMessageContext) IsFinalStep() bool {
	return this.handshakeStep == len(this.Pattern.Messages)
}

func (this *NoiseIKMessageContext) StepOne() {
	if !this.IsHandshakeCompleted() {
		this.handshakeStep += 1
	}
}

func (this *NoiseIKMessageContext) HandshakeDone(csenc *noise.CipherState, csdec *noise.CipherState) {
	this.CSEnc = csenc
	this.CSDec = csdec

}

func (this *NoiseIKMessageContext) Encode(b []byte) (packet []byte, err error) {

	if !this.IsHandshakeCompleted() {
		this.StepOne()

		if this.IsWriteStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1 = this.Hs.WriteMessage(nil, b)

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.IsFinalStep() {
				this.HandshakeDone(cs1, cs0)
			}

		} else {
			return nil, conn.ErrInValidHandshakeStep
		}

	} else {
		packet = this.CSEnc.EncryptWithAd(nil, nil, b)
	}

	return

}

func (this *NoiseIKMessageContext) Decode(b []byte) (packet []byte, err error) {
	if !this.IsHandshakeCompleted() {
		this.StepOne()

		if this.IsReadStep() {
			var cs0, cs1 *noise.CipherState
			packet, cs0, cs1, err = this.Hs.ReadMessage(nil, b)
			if err != nil {
				return
			}

			// final msg for noise handshake pattern we can extract the
			// cipher
			if this.IsFinalStep() {
				this.HandshakeDone(cs0, cs1)
			}

			return

		} else {
			return nil, conn.ErrInValidHandshakeStep
		}

	} else {
		packet, err = this.CSDec.DecryptWithAd(nil, nil, b)
		if err != nil {
			return
		}

	}
	return
}
