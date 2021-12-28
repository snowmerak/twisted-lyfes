package compress

import "io"

type Level interface {
	Level() int
}

type Compressor interface {
	Write(data []byte, buf io.Writer, setting interface{}) error
	Read(reader io.Reader, writer io.Writer) error
}
