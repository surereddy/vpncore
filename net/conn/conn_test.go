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
	"fmt"
	"bytes"
	mrand "math/rand"
	"sync"
	"reflect"
	"time"
	mt "github.com/FTwOoO/vpncore/testing"
	"github.com/FTwOoO/vpncore/crypto"
	"github.com/FTwOoO/vpncore/net/conn/stream/encryption"
	"github.com/FTwOoO/vpncore/net/conn/stream/transport/tcp"
	"github.com/FTwOoO/vpncore/net/conn"
	"github.com/FTwOoO/vpncore/net/conn/message/fragment"
	"github.com/FTwOoO/vpncore/net/conn/message/object/protobuf"
	"github.com/FTwOoO/noise"
	"github.com/FTwOoO/vpncore/net/conn/message/object/msgpack"
	encryption2 "github.com/FTwOoO/vpncore/net/conn/message/encryption"

	"github.com/FTwOoO/vpncore/net/conn/message/transport/udp"
)

func createProtobufTestPackets(n int) []interface{} {

	packets := []interface{}{}

	for {
		if n < 1 {
			break
		}
		mark := (n % 2 == 0)

		msg := &mt.TestPacket{
			Mark: mark,
			Sid:  uint32(n),
			Sessions: map[string]uint64{fmt.Sprintf("hello%d", n):uint64(n), "hi":uint64(n * 2)},

		}

		packets = append(packets, msg)
		n -= 1

	}

	return packets
}

func createMsgpackTestPackets(n int) []interface{} {
	packets := []interface{}{}

	for {
		if n < 1 {
			break
		}
		msg := &msgpack.TestMsg{
			Data: []byte(fmt.Sprintf("hello%d", n)),
		}
		packets = append(packets, msg)
		n -= 1

	}

	return packets
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

func testObjectIOReadWrite(t *testing.T, listener conn.ObjectListener, connection conn.ObjectIO, msgs []interface{}) {
	defer listener.Close()
	defer connection.Close()

	var serverC, clientC conn.ObjectIO

	serverReadSignal := make(chan int, 1)
	serverWriteSignal := make(chan int, 1)
	clientWriteSignal := make(chan int, 1)
	clientReadSignal := make(chan int, 1)
	errChan := make(chan error, 2)

	clientWriteSignal <- 1

	go func() {
		var err error
		serverC, err = listener.Accept()
		if err != nil {
			errChan <- err
			return
		}

		for i, msg := range msgs {
			if i % 2 == 1 {
				<-serverWriteSignal
				fmt.Printf("[S] Write msg %v\n", msg)
				err = serverC.Write(msg)
				if err != nil {
					 errChan <- err
					return
				}

				clientReadSignal <- 1

			} else {

				<-serverReadSignal
				recvMsg, err := serverC.Read()
				if err != nil {
					 errChan <- err
					return
				}

				fmt.Printf("[S] Read msg %v\n", recvMsg)
				if !reflect.DeepEqual(recvMsg, msg) {
					 errChan <- err
					return
				}
				serverWriteSignal <- 1
			}
		}

		errChan <- nil
	}()

	clientC = connection
	go func() {
		for i, msg := range msgs {

			if i % 2 == 1 {
				<-clientReadSignal
				recvMsg, err := clientC.Read()
				if err != nil {
					 errChan <- err
					return
				}

				fmt.Printf("[C] Read msg: %v\n", recvMsg)
				if !reflect.DeepEqual(recvMsg, msg) {
					 errChan <- err
					return
				}

				clientWriteSignal <- 1

			} else {
				<-clientWriteSignal

				fmt.Printf("[C] Write msg: %v\n", msg)
				err := clientC.Write(msg)
				if err != nil {
					 errChan <- err
					return
				}

				serverReadSignal <- 1

			}
		}

		errChan <- nil
	}()

	var count int = 0
	for err := range errChan {
		count += 1


		if err != nil {
			t.Fatal(err)
		}

		if count == 2 {
			return
		}
	}

}

func TestStreamIO(t *testing.T) {
	p := conn.PROTO_TCP
	password := "123456"
	testCiphers := []crypto.StreamCipherName{crypto.AES128CFB, crypto.AES256CFB, crypto.NONE}
	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}

	for _, cipher := range testCiphers {
		for _, testDatalen := range testDatalens {
			port := mrand.Intn(100) + 30000

			fmt.Printf("Test PROTOCOL[%s] with ENCRYPTION[%s] PASS[%s]\n", p, cipher, password)

			context1 := &transport.TCPTransportStreamContext{
				Protocol:p,
				ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
				RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}
			context2 := &encryption.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{Cipher:cipher, Password:password}}

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

func TestStreamToObject(t *testing.T) {
	p := conn.PROTO_TCP
	port := mrand.Intn(100) + 30000
	cipher := crypto.AES256CFB
	password := "123456"

	context1 := &transport.TCPTransportStreamContext{
		Protocol:p,
		ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
		RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}
	context2 := &encryption.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{
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

	testObjectIOReadWrite(t, listener, connection, createProtobufTestPackets(2))

}

func createNoiseIKContextPair() []*encryption2.NoiseIKMessageContext {
	cs := noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256)
	staticI := cs.GenerateKeypair(nil)
	staticR := cs.GenerateKeypair(nil)

	context1, err := encryption2.NewNoiseIKMessageContext(
		cs,
		[]byte("vpncore"),
		staticI,
		noise.DHKey{Public:staticR.Public},
		true,
	)
	if err != nil {

	}

	context2, err := encryption2.NewNoiseIKMessageContext(
		cs,
		[]byte("vpncore"),
		noise.DHKey{},
		staticR,
		false,
	)

	return []*encryption2.NoiseIKMessageContext{context1, context2}
}

func TestStreamToObjectWithNoiseHandshake(t *testing.T) {
	p := conn.PROTO_TCP
	port := mrand.Intn(100) + 30000

	context1 := &transport.TCPTransportStreamContext{
		Protocol:p,
		ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
		RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}

	context2 := new(fragment.FragmentContext)

	contexts := createNoiseIKContextPair()
	context3_I := contexts[0]
	context3_R := contexts[1]

	context4, err := protobuf.NewProtobufMessageContext([]reflect.Type{reflect.TypeOf(&mt.TestPacket{})})
	if err != nil {
		t.Fatal(err)
	}

	contexts_client := []conn.Context{context1, context2, context3_I, context4}
	contexts_server := []conn.Context{context1, context2, context3_R, context4}

	server := new(conn.SimpleServer)
	client := new(conn.SimpleClient)

	listener, err := server.NewListener(contexts_server)
	if err != nil {
		t.Fatal(err)
	}
	connection, err := client.Dial(contexts_client)
	if err != nil {
		t.Fatal(err)
	}

	testObjectIOReadWrite(t, listener, connection, createProtobufTestPackets(4))
}

func TestMessageToObjectWithProtobuf(t *testing.T) {

	port := mrand.Intn(100) + 30000
	context1, err := udp.NewUdpMessageContext(fmt.Sprintf("127.0.0.1:%d", port))
	context2, err := protobuf.NewProtobufMessageContext([]reflect.Type{reflect.TypeOf(&mt.TestPacket{})})
	if err != nil {
		t.Fatal(err)
	}

	server := new(conn.SimpleServer)
	client := new(conn.SimpleClient)

	contexts := []conn.Context{context1, context2}

	listener, err := server.NewListener(contexts)
	if err != nil {
		t.Fatal(err)
	}

	connection, err := client.Dial(contexts)
	if err != nil {
		t.Fatal(err)
	}

	testObjectIOReadWrite(t, listener, connection, createProtobufTestPackets(2))

}

func TestMessageToObjectWithAheadAndMsgpack(t *testing.T) {

	port := mrand.Intn(100) + 30000
	context1, err := udp.NewUdpMessageContext(fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatal(err)
	}
	context2 := encryption2.NewGCM256Context([]byte("Key..."))
	context3 := new(msgpack.MsgpackContext)

	server := new(conn.SimpleServer)
	client := new(conn.SimpleClient)

	contexts := []conn.Context{context1, context2, context3}

	listener, err := server.NewListener(contexts)
	if err != nil {
		t.Fatal(err)
	}

	connection, err := client.Dial(contexts)
	if err != nil {
		t.Fatal(err)
	}

	testObjectIOReadWrite(t, listener, connection, createMsgpackTestPackets(2))

}
