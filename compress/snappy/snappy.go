package compress

import (
	"bytes"
	"io"

	"github.com/golang/snappy"
	"github.com/snowmerak/twisted-lyfes/compress"
)

type Snappy struct{}

func New() compress.Compressor {
	return Snappy{}
}

func (s Snappy) Write(data []byte, writer io.Writer, _ interface{}) error {
	w := snappy.NewBufferedWriter(writer)
	if _, err := w.Write(data); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func (s Snappy) Read(reader io.Reader, writer io.Writer) error {
	r := snappy.NewReader(reader)
	temp := make([]byte, 4096)
	buf := bytes.NewBuffer(nil)
	for {
		n, err := r.Read(temp)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		if n == 0 {
			break
		}
		buf.Write(temp[:n])
	}
	return nil
}
