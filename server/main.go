package server

import (
	"bufio"
	"github.com/markusleevip/taodb"
	"github.com/markusleevip/taodb/leveldb"
	"github.com/markusleevip/taodb/log"
	"io"
	"net"
)

// GLOBALS
var (
	opts *Options
	db  taodb.DB
)

func Main() {
	opts = parseArgs()
	// init logging
	log.LogTo(opts.logto, opts.loglevel)
	db = leveldb.NewDB(opts.DBPath)
	if state,err:=db.State(""); err==nil{
		log.Info(state)
	}
	listen,err := net.Listen("tcp",opts.port)
	if err!=nil{
		log.Error("listen error by port:%s,error:%v",opts.port,err)
		panic(err)
	}
	for{
		conn,err:= listen.Accept()
		if err!=nil{
			log.Error("listen error by port:%s,error:%v",opts.port,err)
		}
		go process(conn)

	}

}

func process(conn net.Conn){
	defer conn.Close()
	reader :=bufio.NewReader(conn)
	for {
		op,err := reader.ReadByte()
		if err !=nil {
			if err!=io.EOF{
				log.Info("client %s is error:%v",conn.RemoteAddr().String(),err)
			}else if err == io.EOF{
				log.Info("client %s is close",conn.RemoteAddr().String())
			}

			return
		}
		server :=Server{}
		if op =='S'{
			server.set(conn,reader)
		}else if op=='G'{
			server.get(conn,reader)
		}else if op=='D'{
			server.del(conn,reader)
		}else if op=='P'{
			server.prefix(conn,reader)
		}else{
			log.Info("close connection due to invalid operation:%v",op)
		}
	}

}

