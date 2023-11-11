package sstcptrans

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans"
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans/internal"
	"github.com/wsrf16/swiss/sugar/base/lambda"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/chacha20poly1305kit"
	"github.com/wsrf16/swiss/sugar/crypto/aeskit/gcmkit"
	"golang.org/x/crypto/hkdf"
	"io"
	"net"
	"strings"
)

type reader struct {
	io.Reader
	cipher.AEAD
	nonce    []byte
	buf      []byte
	leftover []byte
}

type writer struct {
	io.Writer
	cipher.AEAD
	nonce []byte
	buf   []byte
}

// payloadSizeMask is the maximum size of payload in bytes.
const payloadSizeMask = 0x3FFF // 16*1024 - 1

// ErrRepeatedSalt means detected a reused salt
var ErrRepeatedSalt = errors.New("repeated salt detected")

func newReader(r io.Reader, aead cipher.AEAD) *reader {
	return &reader{
		Reader: r,
		AEAD:   aead,
		buf:    make([]byte, payloadSizeMask+aead.Overhead()),
		nonce:  make([]byte, aead.NonceSize()),
	}
}

func newWriter(w io.Writer, aead cipher.AEAD) *writer {
	return &writer{
		Writer: w,
		AEAD:   aead,
		buf:    make([]byte, 2+aead.Overhead()+payloadSizeMask+aead.Overhead()),
		nonce:  make([]byte, aead.NonceSize()),
	}
}

type SSConn struct {
	net.Conn
	Password  string
	AES       aeskit.AESTemplate
	Algorithm SSAlgorithm
	r         *reader
	w         *writer
}

type SSAlgorithm struct {
	Password *string
}

func toSSConn(conn net.Conn, config *sstrans.ShadowSocksConfig) (net.Conn, error) {
	aesTemplate := buildAESTemplate(config.Algorithm)

	ssConn := lambda.If[net.Conn](strings.ToUpper(config.Algorithm) == "PLAIN",
		conn,
		&SSConn{Conn: conn, Password: config.Password, AES: aesTemplate},
	)
	//ssConn := &SSConn{Conn: conn, Password: config.Password, AES: aesTemplate}
	return ssConn, nil
}

func (c *SSConn) getAEAD(salt []byte) (cipher.AEAD, error) {
	subkey, err := c.GetSubKey(salt)
	if err != nil {
		return nil, err
	}
	switch strings.ToUpper(c.AES.Type) {
	case "PLAIN":
		return nil, errors.New("not support")
	case "CFB":
		//stream, err := cfbkit.NewStream(subkey, c.AES.NonceSize)
		//return stream, err
		return nil, errors.New("not support")
	case "GCM":
		return gcmkit.NewAEAD(subkey)
	case "CHACHA20-IETF-POLY1305":
		return chacha20poly1305kit.NewAEAD(subkey)
	default:
		return nil, errors.New("not support")
		//return plainkit.Encrypt(plain)
	}
}

func (c *SSConn) initReader() error {
	salt := make([]byte, c.AES.NonceSize)
	if _, err := io.ReadFull(c.Conn, salt); err != nil {
		return err
	}

	// password + salt + nonce => subkey => aead
	aead, err := c.getAEAD(salt)
	if err != nil {
		return err
	}

	if internal.CheckSalt(salt) {
		return ErrRepeatedSalt
	}

	c.r = newReader(c.Conn, aead)
	return nil
}

func (c *SSConn) initWriter() error {
	salt := make([]byte, c.AES.NonceSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	// password + salt + nonce => subkey => aead
	aead, err := c.getAEAD(salt)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(salt)
	if err != nil {
		return err
	}
	internal.AddSalt(salt)
	c.w = newWriter(c.Conn, aead)
	return nil
}

func (c *SSConn) Read(p []byte) (int, error) {
	if c.r == nil {
		if err := c.initReader(); err != nil {
			return 0, err
		}
	}
	return c.r.Read(p)
}
func (c *SSConn) Write(p []byte) (int, error) {
	if c.w == nil {
		if err := c.initWriter(); err != nil {
			return 0, err
		}
	}
	return c.w.Write(p)
	//n, err := c.Conn.Write(p)
	//if err != nil {
	//    return n, err
	//}
	//encrypt, err := c.AES.Encrypt(p)
	//if err != nil {
	//    return n, err
	//}
	//p = encrypt
	//return n, err
}

func (r *reader) Read(b []byte) (int, error) {
	// copy decrypted bytes (if any) from previous record first
	if len(r.leftover) > 0 {
		n := copy(b, r.leftover)
		r.leftover = r.leftover[n:]
		return n, nil
	}

	n, err := r.read()
	m := copy(b, r.buf[:n])
	if m < n { // insufficient len(b), keep leftover for next read
		r.leftover = r.buf[m:n]
	}
	return m, err
}

// read and decrypt a record into the internal buffer. Return decrypted payload length and any error encountered.
func (r *reader) read() (int, error) {
	// decrypt payload size
	buf := r.buf[:2+r.Overhead()]
	_, err := io.ReadFull(r.Reader, buf)
	if err != nil {
		return 0, err
	}

	_, err = r.Open(buf[:0], r.nonce, buf, nil)
	increment(r.nonce)
	if err != nil {
		return 0, err
	}

	size := (int(buf[0])<<8 + int(buf[1])) & payloadSizeMask

	// decrypt payload
	buf = r.buf[:size+r.Overhead()]
	_, err = io.ReadFull(r.Reader, buf)
	if err != nil {
		return 0, err
	}

	_, err = r.Open(buf[:0], r.nonce, buf, nil)
	increment(r.nonce)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// increment little-endian encoded unsigned integer b. Wrap around on overflow.
func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}

// Write encrypts b and writes to the embedded io.Writer.
func (w *writer) Write(b []byte) (int, error) {
	n, err := w.ReadFrom(bytes.NewBuffer(b))
	return int(n), err
}

// ReadFrom reads from the given io.Reader until EOF or error, encrypts and
// writes to the embedded io.Writer. Returns number of bytes read from r and
// any error encountered.
func (w *writer) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		buf := w.buf
		payloadBuf := buf[2+w.Overhead() : 2+w.Overhead()+payloadSizeMask]
		nr, er := r.Read(payloadBuf)

		if nr > 0 {
			n += int64(nr)
			buf = buf[:2+w.Overhead()+nr+w.Overhead()]
			payloadBuf = payloadBuf[:nr]
			buf[0], buf[1] = byte(nr>>8), byte(nr) // big-endian payload size
			w.Seal(buf[:0], w.nonce, buf[:2], nil)
			increment(w.nonce)

			w.Seal(payloadBuf[:0], w.nonce, payloadBuf, nil)
			increment(w.nonce)

			_, ew := w.Writer.Write(buf)
			if ew != nil {
				err = ew
				break
			}
		}

		if er != nil {
			if er != io.EOF { // ignore EOF as per io.ReaderFrom contract
				err = er
			}
			break
		}
	}

	return n, err
}

//func getNonceSize(name string) (int, error) {
//    name = strings.ToUpper(name)
//    switch name {
//    case "PLAIN":
//        return 16, nil
//    case "CHACHA20-IETF-POLY1305":
//        return 32, nil
//    case "AES-128-GCM":
//        return 16, nil
//    case "AES-196-GCM":
//        return 28, nil
//    case "AES-256-GCM":
//        return 32, nil
//    default:
//        return 0, errors.New("not support " + name)
//    }
//    return 0, nil
//}

func buildAESTemplate(name string) aeskit.AESTemplate {
	var key []byte = nil
	name = strings.ToUpper(name)
	var template aeskit.AESTemplate
	switch name {
	case "PLAIN":
		template = aeskit.AESTemplate{Key: key, NonceSize: 16, Type: "PLAIN"}
	case "CHACHA20-IETF-POLY1305":
		//aeskit.AESTemplate{Type: "CHACHA20", }
		// 32
		//name = aeadChacha20Poly1305
		template = aeskit.AESTemplate{Key: key, NonceSize: 32, Type: "CHACHA20-IETF-POLY1305"}
	case "AES-128-GCM":
		template = aeskit.AESTemplate{Key: key, NonceSize: 16, Type: "GCM"}
	case "AES-192-GCM":
		template = aeskit.AESTemplate{Key: key, NonceSize: 28, Type: "GCM"}
	case "AES-256-GCM":
		template = aeskit.AESTemplate{Key: key, NonceSize: 32, Type: "GCM"}
	}
	return template
}

type SSTemplate struct {
	AES aeskit.AESTemplate
}

// password + salt => subkey
func (ss *SSConn) GetSubKey(salt []byte) ([]byte, error) {
	key := getKey(ss.Password, ss.AES.NonceSize)
	subkey := make([]byte, ss.AES.NonceSize)
	hkdfSHA1(key, salt, []byte("ss-subkey"), subkey)
	return subkey, nil
}

func hkdfSHA1(secret, salt, info, outkey []byte) {
	r := hkdf.New(sha1.New, secret, salt, info)
	if _, err := io.ReadFull(r, outkey); err != nil {
		panic(err) // should never happen
	}
}

func getKey(password string, keyLen int) []byte {
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
