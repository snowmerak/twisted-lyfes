package compress

import (
	"io"

	"github.com/golang/snappy"
	"github.com/snowmerak/twisted-lyfes/compress"
)

type Snappy struct{}

func New() compress.Compressor {
	return Snappy{}
}

func (s Snappy) Write(param compress.WriteParameter) (io.Writer, error) {
	w := snappy.NewBufferedWriter(param.Writer)
	if _, err := w.Write(param.Data); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return param.Writer, nil
}

func (s Snappy) Read(param compress.ReadParameter) (io.Writer, error) {
	r := snappy.NewReader(param.Reader)
	temp := make([]byte, 4096)
	for {
		n, err := r.Read(temp)
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
