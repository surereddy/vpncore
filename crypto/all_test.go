/*
go test -v ./... -bench '^Bench*' -run ^$

BenchmarkAES128-4   	  100000	     11934 ns/op
BenchmarkAES192-4   	  100000	     13296 ns/op
BenchmarkAES256-4   	  100000	     14143 ns/op
BenchmarkTEA-4      	   50000	     34950 ns/op
BenchmarkSimpleXOR-4	 2000000	       774 ns/op
BenchmarkBlowfish-4 	   20000	     59895 ns/op
BenchmarkNone-4     	20000000	       109 ns/op
BenchmarkCast5-4    	   20000	     69818 ns/op
BenchmarkTripleDES-4	    2000	   1008371 ns/op
BenchmarkTwofish-4  	   20000	     89424 ns/op
BenchmarkXTEA-4     	   20000	     80140 ns/op
BenchmarkSalsa20-4  	  300000	      5118 ns/op

 */


package crypto

import (
	"bytes"
	mrand "math/rand"
	crand "crypto/rand"
	"io"
	"testing"
	"fmt"
	"github.com/FTwOoO/vpncore/crypto/cipher"
)

func EnryptionOne(t *testing.T, encrytion Cipher, testKey string, testDataLen int) {
	fmt.Printf("Test %s for EnryptionOne with key[%s] and test data length %d\n", encrytion, testKey, testDataLen)

	bc, err := NewCipher(&EncrytionConfig{Cipher:encrytion, Password:testKey})
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, testDataLen)
	io.ReadFull(crand.Reader, data)
	dec := make([]byte, testDataLen)
	enc := make([]byte, testDataLen)
	bc.Encrypt(enc, data)
	bc.Decrypt(dec, enc)
	if !bytes.Equal(data, dec) {
		t.Fail()
	}
}

func EnryptionStreamingIO(t *testing.T, encrytion Cipher, testKey string, testDataLen int) {
	// create test data
	len1 := mrand.Intn(testDataLen) + testDataLen
	len2 := mrand.Intn(testDataLen) + testDataLen
	len3 := mrand.Intn(testDataLen) + testDataLen
	data1 := make([]byte, len1)
	data2 := make([]byte, len2)
	data3 := make([]byte, len3)
	io.ReadFull(crand.Reader, data1)
	io.ReadFull(crand.Reader, data2)
	io.ReadFull(crand.Reader, data3)

	alldata := make([]byte, len1 + len2 + len3)
	copy(alldata[:len1], data1)
	copy(alldata[len1:len1 + len2], data2)
	copy(alldata[len1 + len2:], data3)
	fmt.Printf("Test %s for EnryptionIO with data length %d-%d-%d\n", encrytion, len1, len2, len3)

	result1 := make([]byte, len(alldata))
	result2 := make([]byte, len(alldata))

	buf1 := bytes.NewBuffer([]byte{})
	cf1, err := NewCipher(&EncrytionConfig{Cipher:encrytion, Password:testKey})
	if err != nil {
		t.Fatal(err)
	}
	r1, err := NewCryptionReadWriter(buf1, cf1)
	if err != nil {
		t.Fatal(err)
	}
	r1.Write(data1)
	r1.Write(data2)
	r1.Write(data3)
	buf1.Read(result1)

	buf2 := bytes.NewBuffer([]byte{})
	cf2, err := NewCipher(&EncrytionConfig{Cipher:encrytion, Password:testKey})
	if err != nil {
		t.Fatal(err)
	}
	r2, err := NewCryptionReadWriter(buf2, cf2)
	if err != nil {
		t.Fatal(err)
	}
	r2.Write(alldata)
	buf2.Read(result2)

	if !bytes.Equal(result1, result2) {
		t.Fatal("Error encryption 1!")
	}

	if encrytion != NONE && bytes.Equal(result1, alldata) {
		t.Fatal("Error encryption 2!")

	}

	r1.Write(alldata)
	r1.Read(result1)

	r2.Write(data1)
	r2.Write(data2)
	r2.Write(data3)
	r2.Read(result2)

	if !bytes.Equal(result1, result2) {
		t.Fatal("Error encryption 3!")
	}

}

func TestAll(t *testing.T) {
	password := "I'm test key"

	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}
	testCiphers := []Cipher{AES128CFB, AES256CFB, SALSA20, NONE}

	for _, testDatalen := range testDatalens {
		for _, cf := range testCiphers {

			EnryptionOne(t, cf, password, testDatalen)

			if cf != SALSA20 {
				EnryptionStreamingIO(t, cf, password, testDatalen)
			}
		}
	}

}

func BenchmarkSalsa20(b *testing.B) {
	var testDataLen = 2047

	pass := make([]byte, 32)
	io.ReadFull(crand.Reader, pass)
	bc, err := cipher.NewSalsa20Stream(pass)
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, testDataLen)
	io.ReadFull(crand.Reader, data)
	dec := make([]byte, testDataLen)
	enc := make([]byte, testDataLen)

	for i := 0; i < b.N; i++ {
		bc.Encrypt(enc, data)
		bc.Decrypt(dec, enc)
	}
}
