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
	"time"
	"bytes"
	mrand "math/rand"
	"sync"
	"github.com/FTwOoO/vpncore/conn/stream/crypt"
	"github.com/FTwOoO/vpncore/conn/stream/transport"
	"github.com/FTwOoO/vpncore/conn"
)

func TestNewListener(t *testing.T) {
	proto := conn.PROTO_TCP
	password := "123456"
	port := mrand.Intn(100) + 20000
	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}
	testCiphers := []crypto.Cipher{crypto.AES128CFB, crypto.AES256CFB, /*enc.SALSA20,*/
		crypto.NONE}

	for _, testDatalen := range testDatalens {
		for _, cipher := range testCiphers {
			fmt.Printf("Test PROTOCOL[%s] with ENCRYPTION[%s] PASS[%s] DATALEN[%d]\n", proto, cipher, password, testDatalen)
			testStreamOnly(t, proto, cipher, port, password, testDatalen)

		}
	}
}

func testStreamOnly(t *testing.T, proto conn.TransportProtocol, cipher crypto.Cipher, port int, password string, testDatalen int) {

	context1 := &transport.TransportStreamContext{
		Protocol:proto,
		ListenAddr:fmt.Sprintf("0.0.0.0:%d", port),
		RemoveAddr:fmt.Sprintf("127.0.0.1:%d", port)}
	context2 := &crypt.CryptStreamContext{EncrytionConfig:&crypto.EncrytionConfig{Cipher:cipher, Password:password}}

	lisener, err := context1.Listen()
	if err != nil {
		t.Fatal(err)
	}
	lisener = &conn.WrapStreamListener{Base:lisener, Contexts:[]conn.StreamContext{context2}}
	defer lisener.Close()

	connection, err := context1.Dial()
	connection = conn.WrapStream([]conn.StreamContext{context2}, connection)
	defer connection.Close()

	testData := mt.RandomBytes(testDatalen)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		receiveData := make([]byte, testDatalen)

		c, err := lisener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		nread := 0
		for {

			n, err := c.Read(receiveData[nread:])
			if err != nil {
				t.Fatal(err)
			}

			mt.PrintBytes(receiveData[nread:], 0x10, "Read bytes")

			nread += n
			if nread >= testDatalen {
				if !bytes.Equal(receiveData, testData) {
					t.Fatal("Bytes does not equal!")
				}
				return
			}

		}

	}()

	<-time.After(3 * time.Second)

	go func() {
		defer wg.Done()
		nwrite := 0
		for {
			n, err := connection.Write(testData[nwrite:])
			if err != nil {
				t.Fatal(err)
			}

			mt.PrintBytes(testData[nwrite:], 0x10, "Write bytes")

			nwrite += n
			if nwrite >= testDatalen {
				break
			}
		}
	}()

	wg.Wait()
}

func testAllStack(t *testing.T) {
	//context3 := new(fragment.FragmentContext)
	//context4 := protobuf.NewProtobufMessageContext()
}