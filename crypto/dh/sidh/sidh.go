package sidh

import (
	"crypto/rand"
	"errors"

	"github.com/cloudflare/circl/dh/sidh"
	"github.com/snowmerak/twisted-lyfes/crypto/dh"
)

const (
	Fp751 = sidh.Fp751
	Fp434 = sidh.Fp434
	Fp503 = sidh.Fp503
)

type KeyPair struct {
	priv *sidh.PrivateKey
	pub  *sidh.PublicKey
	kem  *sidh.KEM
	fp   uint8
}

func NewKeyPair(fp uint8) (dh.DH, error) {
	priv := sidh.NewPrivateKey(fp, sidh.KeyVariantSike)
	if err := priv.Generate(rand.Reader); err != nil {
		return nil, err
	}
	pub := sidh.NewPublicKey(fp, sidh.KeyVariantSike)
	priv.GeneratePublicKey(pub)

	var kem *sidh.KEM
	switch fp {
	case Fp751:
		kem = sidh.NewSike751(rand.Reader)
	case Fp434:
		kem = sidh.NewSike434(rand.Reader)
	case Fp503:
		kem = sidh.NewSike503(rand.Reader)
	default:
		return nil, errors.New("sidh: invalid field size")
	}

	return &KeyPair{
		priv: priv,
		pub:  pub,
		kem:  kem,
		fp:   fp,
	}, nil
}

func NewNotGeneratedKeyPair(fp uint8) (dh.DH, error) {
	kp := &KeyPair{
		priv: sidh.NewPrivateKey(fp, sidh.KeyVariantSike),
		pub:  sidh.NewPublicKey(fp, sidh.KeyVariantSike),
		fp:   fp,
	}

	switch fp {
	case Fp751:
		kp.kem = sidh.NewSike751(rand.Reader)
	case Fp434:
		kp.kem = sidh.NewSike434(rand.Reader)
	case Fp503:
		kp.kem = sidh.NewSike503(rand.Reader)
	default:
		return nil, errors.New("sidh: invalid field size")
	}

	return kp, nil
}

func NewKem(fp uint8) (dh.DH, error) {
	kp := &KeyPair{}

	kp.fp = fp

	switch fp {
	case Fp751:
		kp.kem = sidh.NewSike751(rand.Reader)
	case Fp434:
		kp.kem = sidh.NewSike434(rand.Reader)
	case Fp503:
		kp.kem = sidh.NewSike503(rand.Reader)
	default:
		return nil, errors.New("sidh: invalid field size")
	}

	return kp, nil
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
	other := sidh.NewPublicKey(k.fp, sidh.KeyVariantSike)
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
