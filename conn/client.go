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

type MyClient struct {}

func (client *MyClient) Dial(contexts []Context) (m ObjectIO, err error) {
	if len(contexts) < 1 {
		return nil, ErrInvalidArgs
	}

	for _, ctx := range contexts[:] {
		if _, err = ctx.Valid(); err != nil {
			return nil, ErrInvalidCtx
		}
	}

	ctx := contexts[0]
	if ctx.Layer() != TRANSPORT_LAYER {
		return nil, ErrInvalidCtx
	}

	var c StreamIO
	ctx, ok := ctx.(StreamCreationContext)
	if ok {
		c, err = ctx.(StreamCreationContext).Dial()
		if err != nil {
			return
		}

		var i int
		for i, ctx = range contexts[1:] {

			if _, ok := ctx.(StreamContext); !ok {
				break
			}

			c = ctx.(StreamContext).Pipe(c)

		}

		return client.dialMessageConn(c, contexts[i + 1:])

	} else {
		return client.dialMessageConn(nil, contexts[:])
	}

}

func (client *MyClient) dialMessageConn(c StreamIO, contexts []Context) (oc ObjectIO, err error) {
	if len(contexts) < 2 {
		return nil, ErrInvalidArgs
	}

	ctx := contexts[0]
	var mc MessageIO

	if c != nil {
		ctx, ok := ctx.(StreamToMessageContext)
		if !ok {
			return nil, ErrInvalidCtx
		}

		mc = ctx.(StreamToMessageContext).Pipe(c)
	} else {
		ctx := contexts[0]
		ctx, ok := ctx.(MessageCreationContext)
		if !ok {
			return nil, ErrInvalidCtx
		}

		mc, err = ctx.(MessageCreationContext).Dial()
		if err != nil {
			return
		}
	}


	var i int = -1
	for i, ctx = range contexts[:] {
		if i == 0 {
			continue
		}
		if _, ok := ctx.(MessageContext); !ok {
			break
		}
	}

	if i != len(contexts) - 1 {
		return nil, ErrInvalidCtx
	}

	if i-1 > 0 {
		mc = stackMessageIO{Base:mc, Contexts:contexts[1:i]}
	}


	ctx = contexts[i]
	if _, ok := ctx.(MessageToObjectContext); !ok {
		return nil, ErrInvalidCtx
	}

	oc = transMessageToObjectIO{Context:ctx, Base:mc}
	return
}



/*

func CreateServer(tranProtocol TransportProtocol, address string, cipher crypto.Cipher, password string, codecProtocol link.Protocol) (*link.Server, error) {
	context1 := &transport.TransportStreamContext{
		Protocol:tranProtocol,
		ListenAddr:address,
		RemoveAddr:""}
	context2 := &crypt.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{Cipher:cipher, Password:password}}

	listener, err := NewListener([]StreamContext{context1, context2})
	if err != nil {
		return nil, err
	}

	return link.NewServer(listener, codecProtocol, 0x100), nil
}

func CreateClient(tranProtocol TransportProtocol, address string, cipher crypto.Cipher, password string, codecProtocol link.Protocol) (*link.Client, error) {
	context1 := &transport.TransportStreamContext{
		Protocol:tranProtocol,
		ListenAddr:"",
		RemoveAddr:address}
	context2 := &crypt.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{Cipher:cipher, Password:password}}

	dialer := link.DialerFunc(func() (net.Conn, error) {
		return Dial([]StreamContext{context1, context2})
	})

	client := link.NewClient(dialer, codecProtocol, 2, 50, 0)
	return client, nil
}

*/
