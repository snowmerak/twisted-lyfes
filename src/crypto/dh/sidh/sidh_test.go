package sidh_test

import (
	"bytes"
	"github.com/snowmerak/twisted-lyfes/src/crypto/dh/sidh"
	"testing"
)

func TestKeyExchangeA(t *testing.T) {
	a, err := sidh.NewKem(sidh.Fp751)
	if err != nil {
		t.Fatal(err)
	}
	b, err := sidh.NewKeyPair(sidh.Fp751)
	if err != nil {
		t.Fatal(err)
	}
	ctA, ssA, err := a.Encapsulate(b.ExportPublic())
	if err != nil {
		t.Fatal(err)
	}
	ssB, err := b.Decapsulate(ctA)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ssA, ssB) {
		t.Fatal("secret shares do not match")
	}
	t.Log("secret shares match")
}

func TestKeyExchangeB(t *testing.T) {
	a, err := sidh.NewKeyPair(sidh.Fp503)
	if err != nil {
		t.Fatal(err)
	}
	b, err := sidh.NewKem(sidh.Fp503)
	if err != nil {
		t.Fatal(err)
	}
	ctB, ssB, err := b.Encapsulate(a.ExportPublic())
	if err != nil {
		t.Fatal(err)
	}
	ssA, err := a.Decapsulate(ctB)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ssA, ssB) {
		t.Fatal("secret shares do not match")
	}
	t.Log("secret shares match")
}
