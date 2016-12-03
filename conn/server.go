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

type SimpleServer struct {}

func (server *SimpleServer) NewListener(contexts []Context) (l ObjectListener, err error) {
	if len(contexts) < 2 {

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

	var sl StreamListener
	ctx, ok := ctx.(StreamCreationContext)
	if ok {
		sl, err = ctx.(StreamCreationContext).Listen()
		if err != nil {
			return
		}

		var i int
		ctxs := []StreamContext{}
		for i, ctx = range contexts[:] {
			if i == 0 {
				continue
			}

			if _, ok := ctx.(StreamContext); !ok {
				break
			} else {
				ctxs = append(ctxs, ctx.(StreamContext))
			}
		}

		if len(ctxs) > 0 {
			sl = &WrapStreamListener{Base:sl, Contexts:ctxs}
		}

		return server.newMessageListener(sl, contexts[i:])

	} else {
		return server.newMessageListener(nil, contexts[:])
	}

}

func (server *SimpleServer) newMessageListener(sl StreamListener, contexts []Context) (cl ObjectListener, err error) {
	if len(contexts) < 2 {
		return nil, ErrInvalidArgs
	}

	ctx := contexts[0]
	var ml MessageListener

	if sl != nil {
		ctx, ok := ctx.(StreamToMessageContext)
		if !ok {
			return nil, ErrInvalidCtx
		}

		ml = &transStreamToMessageListener{Base:sl, Context:ctx}
		if err != nil {
			return
		}
	} else {
		ctx := contexts[0]
		ctx, ok := ctx.(MessageCreationContext)
		if !ok {
			return nil, ErrInvalidCtx
		}

		ml, err = ctx.(MessageCreationContext).Listen()
		if err != nil {
			return
		}
	}

	var i int = -1
	ctxs := []MessageContext{}

	for i, ctx = range contexts[:] {
		if i == 0 {
			continue
		}
		if _, ok := ctx.(MessageContext); !ok {
			break
		} else {
			ctxs = append(ctxs, ctx.(MessageContext))
		}
	}


	if i != len(contexts) - 1 {
		return nil, ErrInvalidCtx
	}

	if len(ctxs) > 0 {
		ml = &wrapMessageListener{Base:ml, Contexts:ctxs}
	}

	ctx = contexts[i]
	if _, ok := ctx.(MessageToObjectContext); !ok {
		return nil, ErrInvalidCtx
	}

	cl = &transMessageToObjectListener{Context:ctx.(MessageToObjectContext), Base:ml}
	return
}

