package compress

import (
	"io"

	"github.com/andybalholm/brotli"
	"github.com/snowmerak/twisted-lyfes/compress"
)

type Brotli struct{}

func New() compress.Compressor {
	return Brotli{}
}

// WriteBrotli has default parameter named level
// default level is 6
// you can change level by passing setting parameter implemented by compress.Level interface
func (b Brotli) Write(param compress.WriteParameter) (io.Writer, error) {
	level := param.Level.UnwrapOr(brotli.DefaultCompression)
	brt := brotli.NewWriterLevel(param.Writer, level)
	if _, err := brt.Write(param.Data); err != nil {
		return nil, err
	}
	if err := brt.Flush(); err != nil {
		return nil, err
	}
	if err := brt.Close(); err != nil {
		return nil, err
	}
	return param.Writer, nil
}

func (b Brotli) Read(param compress.ReadParameter) (io.Writer, error) {
	brt := brotli.NewReader(param.Reader)
	temp := make([]byte, 4096)
	for {
		n, err := brt.Read(temp)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		if n == 0 {
			break
		}
		param.Writer.Write(temp[:n])
	}
	return param.Writer, nil
}
