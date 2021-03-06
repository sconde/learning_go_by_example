package secret

import (
	"bytes"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	testMessage = []byte("Do not go gentle into the good night.")
	testKey *[32]byte
)


/* The following tests verify the positive functionality of this package:       
 * can an encrypted message be decrypted?                                        
 */

func TestGenerateKey(t *testing.T) {
	var err error
	testKey, err = GenerateKey()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestEncrypt(t *testing.T) {
	ct, err := Encrypt(testKey, testMessage)
	if err != nil {
		t.Fatalf("%v", err)
	}

	pt, err := Decrypt(testKey, ct)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if !bytes.Equal(testMessage, pt) {
		t.Fatalf("Original and decrypted msgs don't mtch")
	}
}

/* The following test verify the negative functionality of this package:
 * does tit fail when it should?
 */
func prngTester(size int, testFunc func()) {
	prng := rand.Reader
	buf := &bytes.Buffer{}

	rand.Reader = buf
	defer func() { rand.Reader = prng }()

	for i := 0; i < size; i++ {
		tmp := make([]byte, i)
		buf.Write(tmp)
		testFunc()
	}
}

func TestPRNGFailuers(t *testing.T) {
	testFunc := func() {
		_, err := GenerateKey()
		if err == nil {
			t.Fatal("expected key generation failure with bad PRNG")
		}
	}
	prngTester(32, testFunc)

	testFunc = func() {
		_, err := GenerateNonce()
		if err == nil {
			t.Fatal("expected nonce generation failure with bad PRNG")
		}
	}
	prngTester(24, testFunc)

	testFunc = func() {
		_, err := Encrypt(testKey, testMessage)
		if err == nil {
			t.Fatal("expected encryption failure with bad prng")
		}
	}
	prngTester(24, testFunc)
}

func TestDecryptFailures(t *testing.T) {
	targetLenght := 24 + secretbox.Overhead

	for i := 0; i < targetLenght; i++{
		buf := make([]byte, i)
		if _, err := Decrypt(testKey, buf); err == nil {
			t.Fatal("expected decryption failure with bad message lenght")
		}
	}

	otherKey, err:= GenerateKey()
	if err != nil {
		t.Fatal("%v", err)
	}
	ct, err := Encrypt(testKey, testMessage)
	if err != nil {
		t.Fatal("%v", err)
	}
	if _, err := Decrypt(otherKey, ct); err == nil {
		t.Fatal("decrypt should fail ith wrong key")
	}
}
