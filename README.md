# bridge

Make a bridge between tool using websocket protocol.

## Install
```
go get -u github.com/Implex-ltd/bridge/bridge
```

## Example
```go
package main

import (
	"log"
	"time"

	"github.com/Implex-ltd/bridge/bridge"
)

var port = 1337

func server() {
	S, err := bridge.NewServer("", port, func(msg string) {
		log.Printf("recv: %s", msg)

		if msg == "test 3" {
			log.Println("got test #3")
		}
	})

	if err != nil {
		panic(err)
	}

	log.Println("openning..")
	if err := S.Serve(); err != nil {
		panic(err)
	}

	log.Println("Server open")
}

func client() {
	C, err := bridge.NewClient("localhost", port)
	if err != nil {
		panic(err)
	}

	for _, str := range []string{
		"test 1",
		"test 2",
		"test 3",
	} {
		if err := C.PushData(str); err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Millisecond)
	}
}

func main() {
	go server()

	time.Sleep(1 * time.Second)
	client()
}

/**
    Expected output:
        - 2023/09/24 16:01:30 openning..
        - 2023/09/24 16:01:31 recv: test 1
        - 2023/09/24 16:01:31 recv: test 2
        - 2023/09/24 16:01:31 recv: test 3
        - 2023/09/24 16:01:31 got test #3
 */
   
```