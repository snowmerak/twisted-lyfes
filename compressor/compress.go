package compressor

import (
	"io"

	"github.com/snowmerak/generics-for-go/option"
)

type Compressor interface {
	Write(data []byte, buf io.Writer, level *option.Option[int]) error
	Read(reader io.Reader, writer io.Writer) error
}
