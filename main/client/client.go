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
	client := client.New(opts.ip + opts.port)

	for i := 0; i < 100; i++ {
		client.Set(fmt.Sprintf("hello%d", i), []byte(fmt.Sprintf("Hello World!%d", i)))
	}

	for i := 0; i < 100; i++ {
		value, _ := client.Get(fmt.Sprintf("hello%d", i))
		log.Info("get key:hello%d,value=%s\n", i, string(value[:]))
	}

	ctx, _ := client.Prefix("hello")
	if len(ctx) == 0 {
		log.Info("ctx is null")
	} else {
		fmt.Println("cit is not null")
		data := make(map[string]string)
		err := json.Unmarshal(ctx, &data)
		if err != nil {
			log.Error("json error:", err)
		}
		if len(data) > 0 {
			for key, value := range data {
				log.Info("pre.key=%s,%s\n", key, value)
			}
		}
	}
	ctx, _ = client.PrefixOnlyKey("hello")
	if len(ctx) == 0 {
		log.Info("ctx is null")
	} else {
		data := make([]string, 0)
		err := json.Unmarshal(ctx, &data)
		if err != nil {
			log.Error("json error:", err)
			return
		}
		if len(data) > 0 {
			log.Info("data.len=%d", len(data))
			for i, key := range data {
				log.Info("pre.i=%d,key=%s", i, key)
				value, _ := client.Get(key)
				log.Info("getValue=%s", value[:])
			}
		}

	}

}
