package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/markusleevip/taodb/log"
	"github.com/markusleevip/taodb/util"
	"io"
	"net"
)

type Client struct{
	net.Conn
	reader *bufio.Reader
}

func (c *Client) Get(key string){
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s",klen,key)))
}

func (c *Client) Set(key string, value []byte){
	kLen := len(key)
	vLen := len(value)
	head:=fmt.Sprintf("S%d %d %s",kLen,vLen,key)
	var temp bytes.Buffer
	temp.Write([]byte(head))
	temp.Write(value)
	c.Write(temp.Bytes())
}

func (c *Client) recvData() ([]byte,error){
	vlen ,err := util.ReadLen(c.reader)
	if err!=nil{
		log.Error("recvDataerror:%v",err)
	}
	if vlen <0{
		err:=make([]byte,-vlen)
		_,e := io.ReadFull(c.reader,err)
		if e!=nil{
			return nil,e
		}
		return nil, errors.New(string(err))
	}

	value:= make([]byte,vlen)
	_,err = io.ReadFull(c.reader,value)
	if err!=nil{
		return nil,err
	}
	return value,nil

}