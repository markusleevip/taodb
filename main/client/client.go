package main

import (
	"fmt"
	"github.com/markusleevip/taodb/client"
	"github.com/markusleevip/taodb/log"
)
// GLOBALS
var (
	opts *Options
)
func main() {
	opts = parseArgs()

	// init logging
	log.LogTo(opts.logto, opts.loglevel)
	log.Info("port:%s",opts.port)
	log.Info("logto:%s",opts.logto)
	log.Info("loglevel:%s",opts.loglevel)
	//client :=client.New(opts.ip+opts.port)
	client :=client.New("127.0.0.1:7398")
	//defer client.Close()


	for i:=0;i<100;i++ {
		value, _ := client.Set(fmt.Sprintf("hello%d",i),[]byte(fmt.Sprintf("Hello World%d",i)))
		log.Info("set key:hello%d,value=%s", i,string(value[:]))
	}

	for i:=0;i<100;i++ {
		value, _ := client.Get(fmt.Sprintf("hello%d",i))
		//value, _ := client.RecvData()
		log.Info("get key:hello%d,value=%s",i, string(value[:]))
	}

	client.Set("test",[]byte("Hello World!"))
	value,_:=client.Get("test")
	fmt.Println(string(value[:]))

}
