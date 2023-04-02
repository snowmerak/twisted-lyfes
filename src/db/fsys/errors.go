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
