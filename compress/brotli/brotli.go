package compress

import (
	"bytes"
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
func (b Brotli) Write(data []byte, buf io.Writer, setting interface{}) error {
	level := brotli.DefaultCompression
	switch set := setting.(type) {
	case compress.Level:
		level = set.Level()
	}
	brt := brotli.NewWriterLevel(buf, level)
	if _, err := brt.Write(data); err != nil {
		return err
	}
	if err := brt.Flush(); err != nil {
		return err
	}
	if err := brt.Close(); err != nil {
		return err
	}
	return nil
}

func (b Brotli) Read(reader io.Reader, writer io.Writer) error {
	brt := brotli.NewReader(reader)
	buf := bytes.NewBuffer(nil)
	temp := make([]byte, 4096)
	for {
		n, err := brt.Read(temp)
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
