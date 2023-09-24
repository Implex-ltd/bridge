package bridge

import (
	"fmt"
	"log"

	"github.com/dgrr/fastws"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Address   string
	Port      int
	Processor func(msg string)
}

func NewServer(address string, port int, processor func(msg string)) (*Server, error) {
	return &Server{
		Processor: processor,
		Address:   address,
		Port:      port,
	}, nil
}

func (S *Server) Serve() error {
	if err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%v", S.Address, S.Port), fastws.Upgrade(func(conn *fastws.Conn) {
		var msg []byte
		var err error

		for {
			_, msg, err = conn.ReadMessage(msg[:0])
			if err != nil {
				if err != fastws.EOF {
					log.Println(err)
				}
				break
			}

			body := string(msg)

			S.Processor(body)
		}
	})); err != nil {
		return err
	}

	return nil
}
