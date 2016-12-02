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

package conn

import "net"

func WrapStream(contexts []StreamContext, origin StreamIO) (final StreamIO) {

	final = origin
	for _, ctx := range contexts[:] {

		if _, ok := ctx.(StreamContext); !ok {
			break
		}
		final = ctx.(StreamContext).Pipe(final)
	}
	return
}

type WrapStreamListener struct {
	Base     StreamListener
	Contexts []StreamContext
}

func (this *WrapStreamListener) Accept() (StreamIO, error) {
	c, err := this.Base.Accept()
	if err != nil {
		return nil, err
	} else {

		return WrapStream(this.Contexts, c), nil
	}
}

func (this *WrapStreamListener) Close() error {
	return this.Base.Close()
}

func (this *WrapStreamListener) Addr() net.Addr {
	return this.Base.Addr()
}

type transStreamToMessageListener struct {
	Base    StreamListener
	Context StreamToMessageContext
}

func (this *transStreamToMessageListener) Accept() (MessageIO, error) {
	c, err := this.Base.Accept()
	if err != nil {
		return nil, err
	} else {
		return this.Context.Pipe(c), nil
	}
}

func (this *transStreamToMessageListener) Close() error {
	return this.Base.Close()
}
func (this *transStreamToMessageListener) Addr() net.Addr {
	return this.Base.Addr()
}

type wrapMessageListener struct {
	Base     MessageListener
	Contexts []MessageContext
}

func (this *wrapMessageListener) Accept() (MessageIO, error) {
	c, err := this.Base.Accept()
	if err != nil {
		return nil, err
	} else {

		newC := &wrapMessageIO{Base:c, Contexts:this.Contexts}
		return newC, nil
	}
}

func (this *wrapMessageListener) Close() error {
	return this.Base.Close()
}

func (this *wrapMessageListener) Addr() net.Addr {
	return this.Base.Addr()
}

type wrapMessageIO struct {
	Base     MessageIO
	Contexts []MessageContext
}

func (this *wrapMessageIO) Read() (packet []byte, err error) {
	packet, err = this.Base.Read()
	if err != nil {
		return
	}

	for _, ctx := range this.Contexts {
		packet, err = ctx.Decode(packet)
		if err != nil {
			return
		}
	}

	return
}

func (this *wrapMessageIO) Write(packet []byte) error {

	i := len(this.Contexts) - 1
	for {
		if i < 0 {
			break
		}

		ctx := this.Contexts[i]
		packet = ctx.Encode(packet)
		i -= 1
	}

	return this.Base.Write(packet)
}

func (this *wrapMessageIO) Close() error {
	return this.Base.Close()
}

type transMessageToObjectIO struct {
	Base    MessageIO
	Context MessageToObjectContext
}

func (this *transMessageToObjectIO) Read() (interface{}, error) {
	packet, err := this.Base.Read()
	if err != nil {
		return nil, err
	}

	return this.Context.Decode(packet)

}

func (this *transMessageToObjectIO) Write(obj interface{}) error {
	packet, err := this.Context.Encode(obj)
	if err != nil {
		return err
	}

	return this.Base.Write(packet)

}

func (this *transMessageToObjectIO) Close() error {
	return this.Base.Close()
}

type transMessageToObjectListener struct {
	Base    MessageListener
	Context MessageToObjectContext
}

func (this *transMessageToObjectListener) Accept() (ObjectIO, error) {
	c, err := this.Base.Accept()
	if err != nil {
		return nil, err
	} else {

		oc := &transMessageToObjectIO{Base:c, Context:this.Context}
		return oc, nil
	}
}

func (this *transMessageToObjectListener) Close() error {
	return this.Base.Close()
}
func (this *transMessageToObjectListener) Addr() net.Addr {
	return this.Base.Addr()
}