package compress

import (
	"io"

	"github.com/snowmerak/generics-for-go/option"
)

type WriteParameter struct {
	Data   []byte
	Writer io.Writer
	Level  *option.Option[int]
}

type ReadParameter struct {
	Reader io.Reader
	Writer io.Writer
}

type Compressor interface {
	Write(param WriteParameter) (io.Writer, error)
	Read(param ReadParameter) (io.Writer, error)
}
