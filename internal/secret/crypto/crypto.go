package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var (
	ErrEmptyKey       = errors.New("empty encryption key")
	ErrInvalidKey     = errors.New("invalid decryption key")
	ErrInvalidMessage = errors.New("invalid message")
)

var pgpConfig = &packet.Config{
	DefaultHash:            crypto.SHA256,
	DefaultCipher:          packet.CipherAES256,
	DefaultCompressionAlgo: packet.CompressionNone,
}

var pgpHeader = map[string]string{
	"Comment": "https://yopass.se",
}

func Decrypt(r io.Reader, key string) (content, filename string, err error) {
	tried := false
	prompt := func([]openpgp.Key, bool) ([]byte, error) {
		if tried {
			return nil, ErrInvalidKey
		}
		tried = true
		return []byte(key), nil
	}
	a, err := armor.Decode(r)
	if err != nil {
		return "", "", ErrInvalidMessage
	}
	m, err := openpgp.ReadMessage(a.Body, nil, prompt, pgpConfig)
	if err != nil {
		return "", "", fmt.Errorf("could not decrypt: %w", err)
	}
	p, err := io.ReadAll(m.UnverifiedBody)
	if err != nil {
		return "", "", fmt.Errorf("could not read plaintext: %w", err)
	}
	if m.LiteralData.IsBinary {
		filename = m.LiteralData.FileName
	}
	return string(p), filename, nil
}

func Encrypt(r io.Reader, key string) (string, error) {
	if key == "" {
		return "", ErrEmptyKey
	}

	var hints *openpgp.FileHints
	if f, ok := r.(*os.File); ok && r != os.Stdin {
		stat, err := f.Stat()
		if err != nil {
			return "", fmt.Errorf("could not get file info: %w", err)
		}
		hints = &openpgp.FileHints{
			IsBinary: true,
			FileName: stat.Name(),
			ModTime:  stat.ModTime(),
		}
	}

	buf := new(bytes.Buffer)
	a, err := armor.Encode(buf, "PGP MESSAGE", pgpHeader)
	if err != nil {
		return "", fmt.Errorf("could not create armor encoder: %w", err)
	}
	w, err := openpgp.SymmetricallyEncrypt(a, []byte(key), hints, pgpConfig)
	if err != nil {
		return "", fmt.Errorf("could not encrypt: %w", err)
	}
	if _, err := io.Copy(w, r); err != nil {
		return "", fmt.Errorf("could not copy data: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("could not close writer: %w", err)
	}
	if err := a.Close(); err != nil {
		return "", fmt.Errorf("could not close armor: %w", err)
	}

	return buf.String(), nil
}

func GenerateKey() (string, error) {
	const length = 22
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
