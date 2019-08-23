package main

import (
	"flag"
	"github.com/markusleevip/taodb"
	"github.com/markusleevip/taodb/log"
	"net"
)

var flags struct {
	addr, DBPath,logto, loglevel string
}

func init() {
	flag.StringVar(&flags.addr, "addr", ":7398", "The TCP address to bind to")
	flag.StringVar(&flags.DBPath, "dbPath", "/data/storage/taodb", "db save path")
	flag.StringVar(&flags.logto,"log", "stdout", "Write log messages to this file. 'stdout' and 'none' have special meanings")
	flag.StringVar(&flags.loglevel,"log-level", "DEBUG", "The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")
}

func main() {

	flag.Parse()

	log.LogTo(flags.logto, flags.loglevel)
	if err := run(); err != nil {
		log.Error("start error.",err)
		panic(err)
	}
}

func run() error {
	srv := taodb.NewServer(nil, flags.DBPath)
	srv.Handle("ping", taodb.Ping())
	srv.Handle("echo", taodb.Echo())
	srv.Handle("info", taodb.Info(srv))
	srv.Handle("get", taodb.Get(srv))
	srv.Handle("set", taodb.Set(srv))
	srv.Handle("del",taodb.Del(srv))
	srv.Handle("iterator",taodb.Iterator(srv))

	lis, err := net.Listen("tcp", flags.addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	log.Info("waiting for connections on %s", lis.Addr().String())

	return srv.Serve(lis)
}
