#TaoDB
Asynchronous server and client for goleveldb
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
### Starting the client(example)
	import (
		"fmt"
		"github.com/markusleevip/taodb/client"
	)
	func main() {
		client :=client.New("127.0.0.1:7398")
		client.Set("test",[]byte("Hello World!"))
		value,_:=client.Get("test")
		fmt.Println(string(value[:]))
	}
	
