package compress

import (
	"io"

	"github.com/snowmerak/generics-for-go/option"
	"github.com/snowmerak/twisted-lyfes/compress"
	"github.com/ulikunitz/xz/lzma"
)

type LZMA2 struct{}

func New() compress.Compressor {
	return LZMA2{}
}

func (l LZMA2) Write(data []byte, buf io.Writer, level *option.Option[int]) error {
	w, err := lzma.NewWriter2(buf)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (l LZMA2) Read(reader io.Reader, writer io.Writer) error {
	r, err := lzma.NewReader2(reader)
	if err != nil {
		return err
	}
	temp := make([]byte, 4096)
	for !r.EOS() {
		n, err := r.Read(temp)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		writer.Write(temp[:n])
	}
	return nil
}
