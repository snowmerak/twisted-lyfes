package kyber

import (
	"crypto/rand"
	"fmt"

	"github.com/cloudflare/circl/kem/kyber/kyber1024"
	"github.com/snowmerak/twisted-lyfes/crypto/dh"
)

type KeyPair struct {
	priv *kyber1024.PrivateKey
	pub  *kyber1024.PublicKey
}

func New() (dh.DH, error) {
	key := new(KeyPair)
	pub, priv, err := kyber1024.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("kyber.New: %s", err)
	}
	key.priv = priv
	key.pub = pub
	return key, nil
}

func NewNotGenerated() (dh.DH, error) {
	return &KeyPair{
		priv: new(kyber1024.PrivateKey),
		pub:  new(kyber1024.PublicKey),
	}, nil
}

func (k *KeyPair) ExportPrivate() []byte {
	buf := make([]byte, kyber1024.PrivateKeySize)
	k.priv.Pack(buf)
	return buf
}

func (k *KeyPair) ExportPublic() []byte {
	buf := make([]byte, kyber1024.PublicKeySize)
	k.pub.Pack(buf)
	return buf
}

func (k *KeyPair) ImportPrivate(buf []byte) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("kyber.ImportPrivate: %s", r)
				return
			}
			err = nil
		}()
		k.priv.Unpack(buf)
	}()
	return err
}

func (k *KeyPair) ImportPublic(buf []byte) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("kyber.ImportPublic: %s", r)
				return
			}
			err = nil
		}()
		k.pub.Unpack(buf)
	}()
	return err
}

func (k *KeyPair) Encapsulate(pub []byte) (cipherText []byte, secret []byte, err error) {
	k.ImportPublic(pub)
	cipherText = make([]byte, kyber1024.CiphertextSize)
	secret = make([]byte, kyber1024.SharedKeySize)
	k.pub.EncapsulateTo(cipherText, secret, nil)
	return
}

func (k *KeyPair) Decapsulate(cipherText []byte) (secret []byte, err error) {
	secret = make([]byte, kyber1024.SharedKeySize)
	k.priv.DecapsulateTo(secret, cipherText)
	return
}
