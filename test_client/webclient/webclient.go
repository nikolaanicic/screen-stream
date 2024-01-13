package webclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"screen_stream/test_client/events"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	ctx context.Context
	datachan chan []byte
	log *log.Logger
	eventCatcher *events.EventCatcher
	
}


func NewClient(ctx context.Context,l *log.Logger) *Client{
	return &Client{
		ctx:ctx,
		datachan: make(chan []byte),
		log:l,
		eventCatcher: events.New(),
	}
}

func (c *Client) startReadLoop(conn *websocket.Conn) chan []byte{

	d := float64(1000)/float64(30)
	datachan := make(chan []byte)

	ticker := time.NewTicker(time.Millisecond * time.Duration(d))

	go func(){
		for{
			select{
			case <-c.ctx.Done():
				close(datachan)
				c.Close()

				return
	
			case <-ticker.C:
				_, msg, err := conn.ReadMessage()

				if err != nil{
					c.log.Println(err)
	
					if _,ok := err.(*websocket.CloseError); ok{
						close(datachan)
						c.Close()
						return
					}
					
					continue
				}
	
				datachan <- msg
			}
		}
	}()

	return datachan
}

func (c *Client) startEventWriteLoop(conn *websocket.Conn){
	eventCh := c.eventCatcher.Start()

	for {
		select {
		case <-c.ctx.Done():
			c.eventCatcher.Stop()
			return
		case e := <-eventCh:
			conn.WriteJSON(e)
		}
	}
}

func (c *Client) SetOnClose(f func(code int, text string) error){
	if f == nil{
		return
	} 
	c.conn.SetCloseHandler(f)
}


func (c *Client) ConnectEvents(addr string, uname string, pass string) error {
	var err error
	if _, err = url.Parse(addr); err != nil{
		return err
	}

	headers := http.Header{}
	headers.Add("uname",uname)
	headers.Add("pass",pass)

	conn, _, err := websocket.DefaultDialer.Dial(addr, headers)
	if err != nil{
		return err
	}

	fmt.Println("connected:",addr)

	go c.startEventWriteLoop(conn)

	return nil
}



func (c *Client) ConnectScreen(addr string, uname string, pass string) (chan []byte, error){

	var err error
	if _, err = url.Parse(addr); err != nil{
		return nil, err
	}

	headers := http.Header{}
	headers.Add("uname",uname)
	headers.Add("pass",pass)

	conn, _, err := websocket.DefaultDialer.Dial(addr, headers)
	if err != nil{
		return nil, err
	}

	fmt.Println("connected:",addr)

	datachan := c.startReadLoop(conn)

	return datachan, nil
}

func (c *Client) Close() error{
	c.log.Println("closing the connection")
	return c.conn.Close()
}


