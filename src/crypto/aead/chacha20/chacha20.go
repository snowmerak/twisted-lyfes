package chacha20

import (
	"crypto/cipher"
	"github.com/snowmerak/twisted-lyfes/src/crypto/aead"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
)

type Chacha20 struct {
	cipher.AEAD
	nonce []byte
}

func New(key []byte) aead.AEAD {
	return &Chacha20{}
}

func (c *Chacha20) Generate(key []byte) error {
	hashed := blake2b.Sum256(key)
	aead, err := chacha20poly1305.New(hashed[:])
	if err != nil {
		return err
	}

	c.AEAD = aead
	c.nonce = make([]byte, c.AEAD.NonceSize())
	copy(c.nonce, key[:c.AEAD.NonceSize()])

	return nil
}

func (c *Chacha20) Encrypt(plaintext []byte) ([]byte, error) {
	return c.Seal(nil, c.nonce, plaintext, nil), nil
}

func (c *Chacha20) Decrypt(ciphertext []byte) ([]byte, error) {
	return c.Open(nil, c.nonce, ciphertext, nil)
}
