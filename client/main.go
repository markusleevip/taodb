package client
import (
	"bufio"
	"fmt"
	log "github.com/markusleevip/taodb/log"
	"net"
)

// GLOBALS
var (
	opts *Options
)
func Main(){
	opts = parseArgs()

	// init logging
	log.LogTo(opts.logto, opts.loglevel)
	log.Info("port:%s",opts.port)
	log.Info("logto:%s",opts.logto)
	log.Info("loglevel:%s",opts.loglevel)
	conn,err := net.Dial("tcp",opts.ip+opts.port)
	if err!=nil{
		panic(err)
	}
	r := bufio.NewReader(conn)
	client := Client{conn,r}

	for i:=0;i<100000;i++ {
		client.Set(fmt.Sprintf("test%d",i),[]byte(fmt.Sprintf("你好世界%d",i)))
		value, _ := client.recvData()
		fmt.Println("key:test,value=", string(value[:]))
	}

	fmt.Println("读取数据==========================")
	for i:=0;i<100000;i++ {
		client.Get(fmt.Sprintf("test%d",i))
		value, _ := client.recvData()
		fmt.Printf("key:test%d,value=%s\n",i, string(value[:]))
	}
}
