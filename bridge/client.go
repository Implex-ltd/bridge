package bridge

import (
	"fmt"
	"github.com/dgrr/fastws"
)

type Client struct {
	Address string
	Port    int
	Conn    *fastws.Conn
}

func NewClient(address string, port int) (*Client, error) {
	conn, err := fastws.Dial(fmt.Sprintf("ws://%s:%d/ws", address, port))
	if err != nil {
		return nil, err
	}

	return &Client{
		Address: address,
		Port:    port,
		Conn:    conn,
	}, nil
}

func (C *Client) PushData(data string) error {
	if _, err := C.Conn.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}
