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

func wrapStream(contexts []StreamContext, origin StreamReadWriteCloser) (final StreamReadWriteCloser) {

	final = origin
	for _, ctx := range contexts[:] {

		if _, ok := ctx.(StreamContext); !ok {
			break
		}
		final = ctx.(StreamContext).NewPipe(final)
	}
	return
}

type stackStreamListener struct {
	StreamListener
	Contexts []StreamContext
}

func (l *stackStreamListener) Accept() (StreamReadWriteCloser, error) {
	conn, err := l.StreamListener.Accept()
	if err != nil {
		return nil, err
	} else {

		return wrapStream(l.Contexts, conn), nil
	}
}

func wrapMessage(contexts []MessageContext, origin MessagegReadWriteCloser) (final MessagegReadWriteCloser) {

	final = origin
	for _, ctx := range contexts[:] {

		if _, ok := ctx.(MessageContext); !ok {
			break
		}
		final = ctx.(MessageContext).NewPipe(final)
	}
	return
}

type stackMessageListener struct {
	MessageListener
	Contexts []MessageContext
}

func (l *stackMessageListener) Accept() (MessagegReadWriteCloser, error) {
	conn, err := l.MessageListener.Accept()
	if err != nil {
		return nil, err
	} else {
		return wrapMessage(l.Contexts, conn), nil
	}
}

type transMessageListener struct {
	StreamListener
	Context MessageTransitionContext
}

func (l *transMessageListener) Accept() (MessagegReadWriteCloser, error) {
	conn, err := l.StreamListener.Accept()
	if err != nil {
		return nil, err
	} else {
		return l.Context.NewPipe(conn), nil
	}
}
