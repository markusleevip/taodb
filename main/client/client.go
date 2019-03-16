package main

import (
	"encoding/json"
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
	client :=client.New(opts.ip+opts.port)

	for i:=0;i<100;i++ {
		client.Set(fmt.Sprintf("hello%d", i), []byte(fmt.Sprintf("Hello World!%d",i)))
	}

	for i:=0;i<100;i++ {
		value, _ := client.Get(fmt.Sprintf("hello%d",i))
		fmt.Printf("get key:hello%d,value=%s\n",i, string(value[:]))
	}

	ctx,_:=client.Prefix("hello")
	if len(ctx)==0{
		log.Info("ctx is null")
	}else{
		data := make(map[string] string)
		err:=json.Unmarshal(ctx,&data)
		log.Info("%d",len(data))
		if err!=nil{
			log.Error("json error:",err)
		}
		if len(data)>0{
			for key,value := range  data {
				log.Info("pre.key=%s,%s\n",key,value)
			}
		}
	}


}
