package fsys

import "errors"

var ErrPartitionKeyDisMatch = errors.New("partition key dismatch")

func IsErrPartitionKeyDisMatch(err error) bool {
	return errors.Is(err, ErrPartitionKeyDisMatch)
}

var ErrBufferFull = errors.New("buffer full")

func IsErrBufferFull(err error) bool {
	return errors.Is(err, ErrBufferFull)
}

var ErrPartitionNotFound = errors.New("partition not found")

func IsErrPartitionNotFound(err error) bool {
	return errors.Is(err, ErrPartitionNotFound)
}
