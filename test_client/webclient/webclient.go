package webclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	ctx context.Context
	datachan chan []byte
	log *log.Logger
	
}


func NewClient(ctx context.Context,l *log.Logger) *Client{
	return &Client{
		ctx:ctx,
		datachan: make(chan []byte),
		log:l,
	}
}

// the main readling loop of the client
// receives a message and checks for closing errors
// if there is a closing error it closes the connection and returns
// if there is a context cancel from the app it closes the connection and returns

func (c *Client) readLoop(){

	c.log.Println("starting the reading loop")
	d := float64(1000)/float64(25)

	ticker := time.NewTicker(time.Millisecond * time.Duration(d))

	for{
		select{
		case <-c.ctx.Done():
			
			close(c.datachan)
			c.Close()
			return

		case <-ticker.C:
			_, msg, err := c.conn.ReadMessage()
			if err != nil{

				c.log.Println(err)

				if _,ok := err.(*websocket.CloseError); ok{
					close(c.datachan)
					c.Close()

					return
				}
				
				continue
			}

			c.datachan <- msg
		}
	}
}

func (c *Client) SetOnClose(f func(code int, text string) error){
	if f == nil{
		return
	} 
	c.conn.SetCloseHandler(f)
}


func (c *Client) Connect(addr string, uname string, pass string) (chan []byte, error){

	var err error
	if _, err = url.Parse(addr); err != nil{
		return nil, err
	}

	headers := http.Header{}
	headers.Add("uname",uname)
	headers.Add("pass",pass)

	c.conn, _, err = websocket.DefaultDialer.Dial(addr, headers)
	if err != nil{
		return nil, err
	}

	fmt.Println("connected:",addr)

	go c.readLoop()

	return c.datachan, nil
}

func (c *Client) Close() error{
	c.log.Println("closing the connection")
	return c.conn.Close()
}


