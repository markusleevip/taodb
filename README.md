# TaoDB
This is server and client of goleveldb
## Dependencies

### Install goleveldb
-----------
	go get github.com/syndtr/goleveldb/leveldb
### Install TaoDB
-----------
	go get github.com/markusleevip/taodb

## Getting Startted

### Starting the server
-----------
	cd taodb/main/server
	go build
	./server -dbPath=/data/storage/taodb -port=:7398
### Starting the client(test)
-----------
	cd taodb/main/client
	go build
	./client -ip=127.0.0.1 -port=:7398
### Usage the client(example)
	import (
		"fmt"
		"github.com/markusleevip/taodb/client"
	)
	func main() {
		client :=client.New("127.0.0.1:7398")
		client.Set("test",[]byte("Hello World!"))
		value,_:=client.Get("test")
		fmt.Println(string(value[:]))

		// iterator
		for i:=0;i<100;i++ {
        	client.Set(fmt.Sprintf("hello%d", i), []byte(fmt.Sprintf("Hello World!%d",i)))
        	}
        	ctx,_:=client.Prefix("hello")
        	if len(ctx)==0{
        		log.Info("ctx is null")
        	}else{
        		data := make(map[string] string)
        		err:=json.Unmarshal(ctx,&data)
        		log.Info("ctx.len=%d",len(data))
        		if err!=nil{
        			fmt.Println(err)
        		}
        		if len(data)>0{
        			for key,value := range  data {
        				log.Info("pre.key=%s,%s\n",key,value)
        			}
        		}
        	}
	}
	
