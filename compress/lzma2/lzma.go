package compress

import (
	"io"

	"github.com/snowmerak/twisted-lyfes/compress"
	"github.com/ulikunitz/xz/lzma"
)

type LZMA2 struct{}

func New() compress.Compressor {
	return LZMA2{}
}

func (l LZMA2) Write(param compress.WriteParameter) (io.Writer, error) {
	w, err := lzma.NewWriter2(param.Writer)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(param.Data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return param.Writer, nil
}

func (l LZMA2) Read(param compress.ReadParameter) (io.Writer, error) {
	r, err := lzma.NewReader2(param.Reader)
	if err != nil {
		return nil, err
	}
	temp := make([]byte, 4096)
	for !r.EOS() {
		n, err := r.Read(temp)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		param.Writer.Write(temp[:n])
	}
	return param.Writer, nil
}
