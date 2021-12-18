package sidh

import (
	"crypto/rand"
	"errors"

	"github.com/cloudflare/circl/dh/sidh"
	"github.com/twisted-lyfes/utility/dh"
)

type KeyPair struct {
	priv *sidh.PrivateKey
	pub  *sidh.PublicKey
	kem  *sidh.KEM
}

func NewKeyPair() (dh.DH, error) {
	priv := sidh.NewPrivateKey(sidh.Fp751, sidh.KeyVariantSike)
	if err := priv.Generate(rand.Reader); err != nil {
		return nil, err
	}
	pub := sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSike)
	priv.GeneratePublicKey(pub)
	kem := sidh.NewSike751(rand.Reader)

	return &KeyPair{
		priv: priv,
		pub:  pub,
		kem:  kem,
	}, nil
}

func NewNotGeneratedKeyPair() (dh.DH, error) {
	return &KeyPair{
		priv: sidh.NewPrivateKey(sidh.Fp751, sidh.KeyVariantSike),
		pub:  sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSike),
		kem:  sidh.NewSike751(rand.Reader),
	}, nil
}

func NewKem() (dh.DH, error) {
	return &KeyPair{
		priv: nil,
		pub:  nil,
		kem:  sidh.NewSike751(rand.Reader),
	}, nil
}

func (k *KeyPair) ExportPrivate() []byte {
	bs := make([]byte, k.priv.Size())
	k.priv.Export(bs)
	return bs
}

func (k *KeyPair) ExportPublic() []byte {
	bs := make([]byte, k.pub.Size())
	k.pub.Export(bs)
	return bs
}

func (k *KeyPair) ImportPrivate(bs []byte) error {
	if len(bs) != k.priv.Size() {
		return errors.New("invalid private key size")
	}
	if err := k.priv.Import(bs); err != nil {
		return err
	}
	return nil
}

func (k *KeyPair) ImportPublic(bs []byte) error {
	if len(bs) != k.pub.Size() {
		return errors.New("invalid public key size")
	}
	if err := k.pub.Import(bs); err != nil {
		return err
	}
	return nil
}

func (k *KeyPair) Encapsulate(publicKey []byte) (ct []byte, ss []byte, err error) {
	ct = make([]byte, k.kem.CiphertextSize())
	ss = make([]byte, k.kem.SharedSecretSize())
	other := sidh.NewPublicKey(sidh.Fp751, sidh.KeyVariantSike)
	if err := other.Import(publicKey); err != nil {
		return nil, nil, err
	}
	if err := k.kem.Encapsulate(ct, ss, other); err != nil {
		return nil, nil, err
	}
	return ct, ss, nil
}

func (k *KeyPair) Decapsulate(cipherText []byte) (ss []byte, err error) {
	ss = make([]byte, k.kem.SharedSecretSize())
	if err := k.kem.Decapsulate(ss, k.priv, k.pub, cipherText); err != nil {
		return nil, err
	}
	return ss, nil
}
