package aead

type AEAD interface {
	Generate(key []byte) error
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}
