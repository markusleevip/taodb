package util

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ReadLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	l, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return l, nil
}

func SendData(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	vLen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(vLen), value...))
	return e
}
