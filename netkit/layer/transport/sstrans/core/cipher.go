package core

import (
	"crypto/md5"
	"errors"
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans/shadowaead"
	"net"
	"sort"
	"strings"
)

type Cipher interface {
	StreamConnCipher
	PacketConnCipher
}

type StreamConnCipher interface {
	StreamConn(net.Conn) net.Conn
}

type PacketConnCipher interface {
	PacketConn(net.PacketConn) net.PacketConn
}

// ErrCipherNotSupported occurs when a cipher is not supported (likely because of security concerns).
var ErrCipherNotSupported = errors.New("cipher not supported")

const (
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
)

type CipherBuilder struct {
	KeySize int
	New     func(key []byte) (shadowaead.Cipher, error)
}

// List of AEAD ciphers: key size in bytes and constructor
var aeadList = map[string]CipherBuilder{
	aeadAes128Gcm:        CipherBuilder{16, shadowaead.AESGCM},
	aeadAes256Gcm:        CipherBuilder{32, shadowaead.AESGCM},
	aeadChacha20Poly1305: CipherBuilder{32, shadowaead.Chacha20Poly1305},
}

// ListCipher returns a list of available cipher names sorted alphabetically.
func ListCipher() []string {
	var l []string
	for k := range aeadList {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

// PickCipher returns a Cipher of the given name. Derive key from password if given key is empty.
func PickCipher(name string, key []byte, password string) (Cipher, error) {
	name = strings.ToUpper(name)

	switch name {
	case "PLAIN":
		return &plain{}, nil
	case "CHACHA20-IETF-POLY1305":
		name = aeadChacha20Poly1305
	case "AES-128-GCM":
		name = aeadAes128Gcm
	case "AES-256-GCM":
		name = aeadAes256Gcm
	}

	choice, ok := aeadList[name]
	if !ok {
		return nil, ErrCipherNotSupported
	}

	if len(key) == 0 {
		key = kdf(password, choice.KeySize)
	}
	if len(key) != choice.KeySize {
		return nil, shadowaead.KeySizeError(choice.KeySize)
	}
	cipher, err := choice.New(key)
	return &cipherProvider{cipher}, err

}

type cipherProvider struct{ shadowaead.Cipher }

func (aead *cipherProvider) StreamConn(src net.Conn) net.Conn {
	return shadowaead.NewConn(src, aead)
}
func (aead *cipherProvider) PacketConn(src net.PacketConn) net.PacketConn {
	return shadowaead.NewPacketConn(src, aead)
}

// plain cipher does not encrypt
type plain struct{}

func (plain) StreamConn(c net.Conn) net.Conn             { return c }
func (plain) PacketConn(c net.PacketConn) net.PacketConn { return c }

// key-derivation function from original Shadowsocks
func kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}
