package sidh_test

import (
	"bytes"
	"testing"

	"github.com/snowmerak/twisted-lyfes/crypto/dh/sidh"
)

func TestKeyExchangeA(t *testing.T) {
	a, err := sidh.NewKem()
	if err != nil {
		t.Fatal(err)
	}
	b, err := sidh.NewKeyPair()
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
	a, err := sidh.NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	b, err := sidh.NewKem()
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
