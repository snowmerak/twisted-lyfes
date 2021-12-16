package port

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const MIN = uint16(9999)
const MAX = uint16(65535)

func GetFree() (uint16, error) {
	conn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	sp := strings.Split(conn.Addr().String(), ":")
	if len(sp) < 2 {
		return 0, errors.New("not exist port")
	}
	port, err := strconv.ParseUint(sp[1], 10, 16)
	if err != nil {
		return 0, err
	}
	if port > uint64(65535) {
		return 0, errors.New("port is too big")
	}
	return uint16(port), nil
}

func IsAvaliable(port uint16) bool {
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
