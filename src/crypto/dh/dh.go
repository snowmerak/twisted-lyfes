package dh

type DH interface {
	ExportPrivate() []byte
	ExportPublic() []byte
	ImportPrivate([]byte) error
	ImportPublic([]byte) error
	Encapsulate([]byte) (cipherText []byte, secret []byte, err error)
	Decapsulate([]byte) (Secret []byte, err error)
}
