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

func wrapStream(contexts []StreamContext, origin StreamIO) (final StreamIO) {

	final = origin
	for _, ctx := range contexts[:] {

		if _, ok := ctx.(StreamContext); !ok {
			break
		}
		final = ctx.(StreamContext).Pipe(final)
	}
	return
}

type stackStreamListener struct {
	StreamListener
	Contexts []StreamContext
}

func (l *stackStreamListener) Accept() (StreamIO, error) {
	c, err := l.StreamListener.Accept()
	if err != nil {
		return nil, err
	} else {

		return wrapStream(l.Contexts, c), nil
	}
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

type stackMessageListener struct {
	MessageListener
	Contexts []MessageContext
}

func (l *stackMessageListener) Accept() (MessageIO, error) {
	c, err := l.MessageListener.Accept()
	if err != nil {
		return nil, err
	} else {

		newC := stackMessageIO{Base:c, Contexts:l.Contexts}
		return newC, nil
	}
}

type stackMessageIO struct {
	Base     MessageIO
	Contexts []MessageContext
}

func (this *stackMessageIO) Read() (packet []byte, err error) {
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

func (this *stackMessageIO) Write(packet []byte) error {

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

func (this *stackMessageIO) Close() error {
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

		oc := transMessageToObjectIO{Base:c, Context:this.Context}
		return oc, nil
	}
}