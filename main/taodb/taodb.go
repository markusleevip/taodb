package main

import (
	"flag"
	"fmt"
	"github.com/markusleevip/taodb/client"
	"github.com/markusleevip/taodb/log"
	"github.com/markusleevip/taodb/resp"
	"net"
)

var flags struct {
	addr string
	logto    string
	loglevel string
}

func init() {
	flag.StringVar(&flags.addr,"addr", ":7398", "server addr")
	flag.StringVar(&flags.logto,"log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	flag.StringVar(&flags.loglevel,"log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
}

func main() {
	flag.Parse()

	// init logging
	log.LogTo(flags.logto, flags.loglevel)

	pool,_ := client.New(func() (net.Conn, error) {
		return net.Dial("tcp", flags.addr)
	})
	defer pool.Close()
	cn, _  :=pool.Get()
	defer pool.Put(cn)

	cn.WriteCmdString("PING")
	cn.WriteCmdString("ECHO", "HEllO")
	cn.WriteCmd("SET", []byte("hello"), []byte("Hello 世界"))
	cn.WriteCmdString("GET", "hello")
	cn.WriteCmdString("ITERATOR", "hello")
	if err := cn.Flush(); err != nil {
		cn.MarkFailed()
		panic(err)
	}

	// Consume responses
	for i := 0; i < 5; i++ {
		t, err := cn.PeekType()
		if err != nil {
			return
		}

		switch t {
		case resp.TypeInline:
			s, _ := cn.ReadInlineString()
			fmt.Println(s)
		case resp.TypeBulk:
			s, _ := cn.ReadBulk(nil)
			fmt.Println(string(s[:]))
		case resp.TypeInt:
			n, _ := cn.ReadInt()
			fmt.Println(n)
		case resp.TypeNil:
			_ = cn.ReadNil()
			fmt.Println(nil)
		case resp.TypeError:
			err,_:=cn.ReadError()
			fmt.Println(err)
		default:
			panic("unexpected response type")
		}
	}



}
