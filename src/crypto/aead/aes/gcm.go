package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"github.com/snowmerak/twisted-lyfes/src/crypto/aead"
)

type GCM struct {
	cipher.AEAD
	nonce []byte
}

func New() aead.AEAD {
	return &GCM{}
}

func (g *GCM) Generate(key []byte) error {
	hashed := sha256.Sum256(key)
	block, err := aes.NewCipher(hashed[:])
	if err != nil {
		return err
	}

	g.AEAD, err = cipher.NewGCM(block)
	if err != nil {
		return err
	}

	g.nonce = make([]byte, g.AEAD.NonceSize())
	copy(g.nonce, key[:g.AEAD.NonceSize()])

	return nil
}

func (g *GCM) Encrypt(plaintext []byte) ([]byte, error) {
	return g.Seal(nil, g.nonce, plaintext, nil), nil
}

func (g *GCM) Decrypt(ciphertext []byte) ([]byte, error) {
	return g.Open(nil, g.nonce, ciphertext, nil)
}
