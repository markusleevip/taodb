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
	cd taodb
	./build.sh
	./taodbd -dbPath=/data/storage/taodb -port=:7398
### Starting the client(test)
-----------
	cd taodb
	./taodb -ip=127.0.0.1 -port=:7398
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
        		fmt.Println("ctx is null")
        	}else{
        		data := make(map[string] string)
        		err:=json.Unmarshal(ctx,&data)
        		fmt.Printf("ctx.len=%d\n",len(data))
        		if err!=nil{
        			fmt.Println(err)
        		}
        		if len(data)>0{
        			for key,value := range  data {
        				fmt.Printf("pre.key=%s,%s\n",key,value)
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
	
