package fsys

import (
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
)

type FSys struct {
	basePath string
}

func (fs *FSys) GetMetaData(name string) (*MetaData, error) {
	data, err := os.ReadFile(filepath.Join(fs.basePath, name, MetaDataFileName))
	if err != nil {
		return nil, err
	}

	md := new(MetaData)
	if err := proto.Unmarshal(data, md); err != nil {
		return nil, err
	}

	return md, nil
}
