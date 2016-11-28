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




type transMessageListener struct {
	StreamListener
	Context MessageTransitionContext
}

func (l *transMessageListener) Accept() (MessageIO, error) {
	c, err := l.StreamListener.Accept()
	if err != nil {
		return nil, err
	} else {
		return l.Context.Pipe(c), nil
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

func (this *stackMessageIO) Read(msg Message) error {
	lastMsg := msg
	for _, ctx := range this.Contexts {
		lastMsg = ctx.PipeMessage(lastMsg)
	}

	return this.Base.Read(lastMsg)
}

func (this *stackMessageIO) Write(msg Message) error {
	lastMsg := msg
	for _, ctx := range this.Contexts {
		lastMsg = ctx.PipeMessage(lastMsg)
	}

	return this.Base.Write(lastMsg)
}

func (this *stackMessageIO) Close() error {
	this.Base.Close()
}