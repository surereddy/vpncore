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

package conn_test

import (
	"testing"
	mt "github.com/FTwOoO/vpncore/testing"
	"github.com/FTwOoO/vpncore/crypto"
	"fmt"
	"bytes"
	mrand "math/rand"
	"sync"
	"github.com/FTwOoO/vpncore/conn/stream/crypt"
	"github.com/FTwOoO/vpncore/conn/stream/transport"
	"github.com/FTwOoO/vpncore/conn"
	"time"
	"github.com/FTwOoO/vpncore/conn/message/fragment"
	"github.com/FTwOoO/vpncore/conn/message/protobuf"
	"reflect"
	"github.com/golang/protobuf/proto"
)

func TestStreamIO(t *testing.T) {
	p := conn.PROTO_TCP
	password := "123456"
	testCiphers := []crypto.Cipher{crypto.AES128CFB, crypto.AES256CFB, /*enc.SALSA20,*/
		crypto.NONE}
	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}

	for _, cipher := range testCiphers {
		for _, testDatalen := range testDatalens {
			port := mrand.Intn(100) + 30000

			fmt.Printf("Test PROTOCOL[%s] with ENCRYPTION[%s] PASS[%s]\n", p, cipher, password)

			context1 := &transport.TransportStreamContext{
				Protocol:p,
				ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
				RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}
			context2 := &crypt.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{Cipher:cipher, Password:password}}

			listener, err := context1.Listen()
			if err != nil {
				t.Fatal(err)
			}
			listener = &conn.WrapStreamListener{Base:listener, Contexts:[]conn.StreamContext{context2}}

			connection, err := context1.Dial()
			connection = conn.WrapStream([]conn.StreamContext{context2}, connection)

			testStreamIOReadWrite(t, listener, connection, testDatalen)
		}
	}
}

func testStreamIOReadWrite(t *testing.T, listener conn.StreamListener, connection conn.StreamIO, testDatalen int) {
	defer listener.Close()
	defer connection.Close()

	fmt.Printf("Test data length :%d\n", testDatalen)

	testData := mt.RandomBytes(testDatalen)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		receiveData := make([]byte, testDatalen)

		c, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		nread := 0
		for {

			n, err := c.Read(receiveData[nread:])
			if err != nil {
				t.Fatal(err)
			}

			mt.PrintBytes(receiveData[nread:], 0x10, fmt.Sprintf("Read %d bytes", n))

			nread += n
			if nread >= testDatalen {
				if !bytes.Equal(receiveData, testData) {
					t.Fatal("Bytes does not equal!")
				}
				return
			}

		}

	}()

	<-time.After(1 * time.Second)

	go func() {
		defer wg.Done()
		nwrite := 0
		for {
			n, err := connection.Write(testData[nwrite:])
			if err != nil {
				t.Fatal(err)
			}

			mt.PrintBytes(testData[nwrite:], 0x10, fmt.Sprintf("Write %d bytes", n))

			nwrite += n
			if nwrite >= testDatalen {
				return
			}
		}
	}()

	wg.Wait()

}

func TestAllStack(t *testing.T) {
	p := conn.PROTO_TCP
	password := "123456"
	port := mrand.Intn(100) + 30000
	cipher := crypto.AES256CFB

	context1 := &transport.TransportStreamContext{
		Protocol:p,
		ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
		RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}
	context2 := &crypt.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{
		Cipher:cipher, Password:password,
	},
	}

	context3 := new(fragment.FragmentContext)
	context4, err := protobuf.NewProtobufMessageContext([]reflect.Type{reflect.TypeOf(&mt.TestPacket{})})
	if err != nil {
		t.Fatal(err)
	}

	contexts := []conn.Context{context1, context2, context3, context4}
	server := new(conn.SimpleServer)
	client := new(conn.SimpleClient)

	listener, err := server.NewListener(contexts)
	connection, err := client.Dial(contexts)

	sendMsg1 := &mt.TestPacket{
		Mark: false,
		Sid:  999,
		Sessions: map[string]uint64{"a":1, "b":2},

	}

	sendMsg2 := &mt.TestPacket{
		Mark: true,
		Sid:  18,
		Sessions: map[string]uint64{"xxx":10, "yyy":20},

	}

	testMessageIOReadWrite(t, listener, connection, []proto.Message{sendMsg1, sendMsg2})

}

func testMessageIOReadWrite(t *testing.T, listener conn.ObjectListener, connection conn.ObjectIO, msgs []proto.Message) {
	defer listener.Close()
	defer connection.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		c, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("Accept a conn !\n")

		for i, msg := range msgs {
			if i % 2 == 1 {
				err = c.Write(msg)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Printf("[L] Write msg %v\n", msg)

			} else {
				recvMsg, err := c.Read()
				if err != nil {
					t.Fatal(err)
				}

				fmt.Printf("[L] Read msg %v\n", msg)

				if !mt.ProtoMessageEqual(recvMsg.(proto.Message), msg) {
					t.Fatal()
				}
			}
		}

	}()

	<-time.After(1 * time.Second)

	go func() {
		defer wg.Done()
		c := connection

		for i, msg := range msgs {

			if i % 2 == 0 {
				err := c.Write(msg)
				if err != nil {
					t.Fatal(err)
				}

				fmt.Printf("[C] Write msg %v\n", msg)

			} else {
				recvMsg, err := c.Read()
				if err != nil {
					t.Fatal(err)
				}

				fmt.Printf("[C] Read msg %v\n", msg)

				if !mt.ProtoMessageEqual(recvMsg.(proto.Message), msg) {
					t.Fatal()
				}
			}
		}
	}()

	wg.Wait()
}

