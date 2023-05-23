package server

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"screen_stream/encoder"
	"screen_stream/screenmgr"

	"github.com/gorilla/websocket"
)

var(
	upgrader = websocket.Upgrader{}
)

type Server struct {
	http.Server
	ctx context.Context
	display *screenmgr.Display
	cancel context.CancelFunc
	options Options
	s *screenmgr.DisplayStream
	enc *encoder.Encoder
}

var DefaultOptions Options = Options{sampleRate: 30}

func NewServer(ctx context.Context, cancel context.CancelFunc) Server{
	return Server{
		ctx:ctx,
		cancel:cancel,
		enc: encoder.NewEncoder(),
		display: screenmgr.NewDisplay(0),
		options: DefaultOptions,
	}
}

func (s *Server) WithSampleRate(sampleRate int) *Server{
	s.options.sampleRate = sampleRate
	return s
}

func (s *Server) Stop(){
	fmt.Println("stopping the server")
	s.cancel()
}


func (s *Server) initStream() chan *image.RGBA{
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx, s.cancel = ctx,cancel

	s.s = screenmgr.NewStream(s.options.sampleRate,s.ctx,s.display)
	return s.s.Start()
}


func (s *Server) GetDisplayStreamHandler() func(http.ResponseWriter, *http.Request){
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool {return true}

		fmt.Println("connected")

		ws, err := upgrader.Upgrade(w,r,nil)
		if err != nil{
			w.Write([]byte(err.Error()))
			return
		}

		defer ws.Close()
		fmt.Println("upgraded")


		ws.SetCloseHandler(func(code int, text string) error {
			fmt.Println("conn closed")
			s.cancel()
			return nil
		})

		ch := s.initStream()

		fmt.Println("stream open")

		for {

			select{
			case <-s.ctx.Done():
				return
			case x := <- ch:
				res, err := s.enc.BytesToBase64(x)
				if err != nil{
					fmt.Println(err)
					s.Stop()
					return
				}

				err = ws.WriteMessage(1, []byte(res))

				if err != nil{
					fmt.Println(err)
					s.Stop()
					return
				} 
			}

		}
	})
}