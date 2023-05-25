package server

import (
	"context"
	"fmt"
	"net/http"
	"screen_stream/encoder"
	"screen_stream/screenmgr"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

type Server struct {
	ctx     context.Context
	display *screenmgr.Display
	cancel  context.CancelFunc
	options Options
	enc     *encoder.Encoder
}

var DefaultOptions Options = Options{sampleRate: 30}

func NewServer(ctx context.Context, cancel context.CancelFunc) Server {
	return Server{
		ctx:     ctx,
		cancel:  cancel,
		enc:     encoder.NewEncoder(),
		display: screenmgr.NewDisplay(0),
		options: DefaultOptions,
	}

}

func (s *Server) WithSampleRate(sampleRate int) *Server {
	s.options.sampleRate = sampleRate
	return s
}

func (s *Server) Stop() {
	fmt.Println("stopping the server")
	s.cancel()
}

func (s *Server) SpawnNewStream() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// creating the stream object for the current client
		stream := screenmgr.NewStream(s.options.sampleRate, s.display)

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		defer func(){
			ws.Close()
			stream.Stop()
		}()

		// when the ws handler receives a close message it should stop the stream
		// and send back the close handshake message
		ws.SetCloseHandler(func(code int, text string) error {
			stream.Stop()
			return ws.WriteControl(websocket.CloseMessage,[]byte{},time.Now().Add(time.Second))
		})


		// starting a function that reads the connection
		// above function in setclosehandler gets triggered 
		// when a close message is read in ReadMessage()
		// reading the connection every second to see if it was closed already
		go func(){
			ticker := time.NewTicker(time.Second)
			for{
				select{
				case <- stream.Wait():
					return
				case <-ticker.C:
					ws.ReadMessage()
				}
			}
		}()


		ch := stream.Start()
		for {
			select{
			case <-stream.Wait():
				return
			case x := <- ch:
				
				// encoding image.RGBA to jpeg []byte
				// sending the raw jpeg bytes to the clients
				res, err := s.enc.BytesToJpeg(x)
				if err != nil {
					return
				}else if err = ws.WriteMessage(websocket.BinaryMessage, []byte(res)); err != nil {
					return
				}
			}
		}
	})
}

