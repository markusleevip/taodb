# TaoDB
This is server and client of goleveldb
Also use Redis client to connect to the server.
## Dependencies

### Install TaoDB
-----------
	go get github.com/markusleevip/taodb

## Getting Startted

### Starting the server
-----------
	cd taodb
	./build.sh
	./taodbd -dbPath=/data/storage/taodb -addr=:7398
### Starting the client(test)
-----------
	cd taodb
	./taodb -addr=:7398
### Usage the client(example)
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
    	cn.WriteCmd("SET", []byte("hello"), []byte("Hello 世界 哈哈"))
    	cn.WriteCmdString("GET", "key")
    	if err := cn.Flush(); err != nil {
    		cn.MarkFailed()
    		panic(err)
    	}
    
    	// Consume responses
    	for i := 0; i < 4; i++ {
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
    
    	/*
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
    	*/
    
    }
	
