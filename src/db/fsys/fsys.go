package fsys

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FSys struct {
	basePath      string
	hashFunc      func([]byte) []byte
	stringEncoder func([]byte) string

	chunkSize int
}

func New(basePath string, hashFunc func([]byte) []byte, stringEncoder func([]byte) string) *FSys {
	return &FSys{
		basePath:      basePath,
		hashFunc:      hashFunc,
		stringEncoder: stringEncoder,
		chunkSize:     1024 * 1024,
	}
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

func (fs *FSys) PutMetaData(metaData *MetaData) error {
	data, err := proto.Marshal(metaData)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(fs.basePath, metaData.Name, MetaDataFileName), data, 0o775); err != nil {
		return err
	}

	return nil
}

func (fs *FSys) GetPartition(name, part string) (*Partition, error) {
	data, err := os.ReadFile(filepath.Join(fs.basePath, name, part+".part"))
	if err != nil {
		return nil, err
	}

	partition := new(Partition)
	if err := proto.Unmarshal(data, partition); err != nil {
		return nil, err
	}

	return partition, nil
}

func (fs *FSys) PutPartition(name string, part *Partition) error {
	data, err := proto.Marshal(part)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(fs.basePath, name, part.Key+".part"), data, 0o775); err != nil {
		return err
	}

	return nil
}

func (fs *FSys) makePartKey(partition []byte) string {
	return fs.stringEncoder(fs.hashFunc(partition))
}

func (fs *FSys) CheckParts(metaData *MetaData) ([]string, []string, error) {
	basePath := filepath.Join(fs.basePath, metaData.Name)

	dirs, err := os.ReadDir(basePath)
	if err != nil {
		return nil, nil, err
	}

	partMap := map[string]struct{}{}
	for _, p := range metaData.Partitions {
		partMap[p] = struct{}{}
	}

	exists := []string(nil)

	part := new(Partition)
	for _, dir := range dirs {
		if dir.IsDir() {
			continue
		}

		name := dir.Name()

		if _, ok := partMap[name]; !ok {
			continue
		}

		data, err := os.ReadFile(filepath.Join(basePath, name))
		if err != nil {
			return nil, nil, err
		}

		if err := proto.Unmarshal(data, part); err != nil {
			continue
		}

		hashed := fs.makePartKey(data)
		if part.Key != hashed {
			continue
		}

		delete(partMap, name)
		exists = append(exists, name)
	}

	notFounds := []string(nil)
	for k := range partMap {
		notFounds = append(notFounds, k)
	}

	return exists, notFounds, nil
}

func (fs *FSys) CombinePartsTo(writer io.Writer, metaData *MetaData) error {
	basePath := filepath.Join(fs.basePath, metaData.Name)

	partition := new(Partition)
	for _, part := range metaData.Partitions {
		partPath := filepath.Join(basePath, part+".part")

		data, err := os.ReadFile(partPath)
		if err != nil {
			return err
		}

		if err := proto.Unmarshal(data, partition); err != nil {
			return err
		}

		if part != partition.Key || partition.Key != fs.makePartKey(partition.Data) {
			return ErrPartitionKeyDisMatch
		}

		n, err := writer.Write(partition.Data)
		if err != nil {
			return err
		}
		if len(data) > n {
			return ErrBufferFull
		}
	}

	return nil
}

func (fs *FSys) Delete(name string) error {
	if err := os.RemoveAll(filepath.Join(fs.basePath, name)); err != nil {
		return err
	}
	return nil
}

func (fs *FSys) Bring(path string) error {
	name := filepath.Base(path)

	metaData := new(MetaData)
	metaData.Name = name
	metaData.Timestamp = time.Now().Unix()

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	buf := make([]byte, fs.chunkSize)

	part := new(Partition)
	for {
		n, err := f.Read(buf)
		if n == 0 && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		part.Data = make([]byte, n)
		copy(part.Data, buf[:n])

		part.Key = fs.makePartKey(part.Data)

		if err := fs.PutPartition(name, part); err != nil {
			return err
		}

		metaData.Partitions = append(metaData.Partitions, part.Key)
	}

	if err := fs.PutMetaData(metaData); err != nil {
		return err
	}

	return nil
}
