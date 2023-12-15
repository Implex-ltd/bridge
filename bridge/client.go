package bridge

import (
	"fmt"
	"github.com/dgrr/fastws"
	"sync"
	"time"
)

type Client struct {
	Address   string
	Port      int
	Conn      *fastws.Conn
	Connected bool

	ToPush []string
	mut    *sync.Mutex
}

func NewClient(address string, port int) (*Client, error) {
	C := Client{
		Address: address,
		Port:    port,
		ToPush:  []string{},
		mut:     &sync.Mutex{},
	}

	if err := C.Connect(); err != nil {
		return nil, err
	}

	return &C, nil
}

func (C *Client) Connect() (err error) {
	C.Conn, err = fastws.Dial(fmt.Sprintf("ws://%s:%d/ws", C.Address, C.Port))
	if err != nil {
		return err
	}

	C.Connected = true
	go C.Heartbeat()
	go C.ProcessUnpushed()

	return nil
}

func (C *Client) Heartbeat() error {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		if _, err := C.Conn.Write([]byte("ping")); err != nil {
			if err := C.Connect(); err != nil {
				C.Connected = false
				return fmt.Errorf("client disconnected, cannot reconnect")
			}
			C.Connected = true
		}
	}

	return nil
}

func (C *Client) PushData(data string) error {
	if !C.Connected {
		C.mut.Lock()
		defer C.mut.Unlock()

		C.ToPush = append(C.ToPush, data)
		return fmt.Errorf("client is not connected")
	}

	if _, err := C.Conn.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func (C *Client) ProcessUnpushed() error {
	C.mut.Lock()
	defer C.mut.Unlock()

	var notPushed []string

	for _, e := range C.ToPush {
		if err := C.PushData(e); err != nil {
			notPushed = append(notPushed, e)
		}
	}

	C.ToPush = notPushed

	if len(notPushed) > 0 {
		return fmt.Errorf("failed to push some items")
	}

	return nil
}
