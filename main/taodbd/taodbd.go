package main

import (
	"flag"
	"log"
	"github.com/markusleevip/taodb"
	"net"
)


var flags struct {
	addr ,DBPath string

}

func init() {
	flag.StringVar(&flags.addr, "addr", ":7398", "The TCP address to bind to")
	flag.StringVar(&flags.DBPath,"dbPath", "/data/storage/leveldb", "db save path")
}

func main() {

	flag.Parse()

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	srv := taodb.NewServer(nil,flags.DBPath)
	srv.Handle("ping",taodb.Ping())
	srv.Handle("echo",taodb.Echo())
	srv.Handle("info",taodb.Info(srv))
	srv.Handle("get",taodb.Get(srv))
	srv.Handle("set",taodb.Set(srv))

	lis, err := net.Listen("tcp", flags.addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	log.Printf("waiting for connections on %s", lis.Addr().String())

	return srv.Serve(lis)
}
