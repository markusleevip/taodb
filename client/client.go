package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/markusleevip/taodb/log"
	"github.com/markusleevip/taodb/util"
	"io"
	"net"
)

type Client struct {
	net.Conn
	reader *bufio.Reader
}


func (c *Client) Get(key string) ([]byte, error) {
	kLen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", kLen, key)))
	return c.RecvData()
}

func (c *Client) Set(key string, value []byte) ([]byte, error) {
	kLen := len(key)
	vLen := len(value)
	head := fmt.Sprintf("S%d %d %s", kLen, vLen, key)
	var temp bytes.Buffer
	temp.Write([]byte(head))
	temp.Write(value)
	c.Write(temp.Bytes())
	return c.RecvData()
}

func (c *Client) Del(key string) ([]byte, error) {
	kLen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", kLen, key)))
	return c.RecvData()
}

func (c *Client) Prefix(key string) ([]byte, error) {
	kLen := len(key)
	c.Write([]byte(fmt.Sprintf("P%d %s", kLen, key)))
	return c.RecvData()
}

func (c *Client) PrefixOnlyKey(key string) ([]byte, error) {
	kLen := len(key)
	c.Write([]byte(fmt.Sprintf("K%d %s", kLen, key)))
	return c.RecvData()
}

func (c *Client) RecvData() ([]byte, error) {
	vLen, err := util.ReadLen(c.reader)
	if err != nil {
		log.Error("recvData.error:%v", err)
	}
	if vLen < 0 {
		err := make([]byte, -vLen)
		_, e := io.ReadFull(c.reader, err)
		if e != nil {
			return nil, e
		}
		return nil, errors.New(string(err))
	}

	value := make([]byte, vLen)
	_, err = io.ReadFull(c.reader, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
