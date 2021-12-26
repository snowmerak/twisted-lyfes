package kyber_test

import (
	"bytes"
	"testing"

	"github.com/snowmerak/twisted-lyfes/crypto/dh/kyber"
)

func TestKeyExchange(t *testing.T) {
	server, err := kyber.New()
	if err != nil {
		t.Fatal(err)
	}

	client, err := kyber.NewNotGenerated()
	if err != nil {
		t.Fatal(err)
	}

	ct, clientSecret, err := client.Encapsulate(server.ExportPublic())
	if err != nil {
		t.Fatal(err)
	}

	serverSecret, err := server.Decapsulate(ct)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(serverSecret, clientSecret) {
		t.Fatal("client and server secrets do not match")
	}

}
