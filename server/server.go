package server

import (
	"log"
	"net/http"
	"screen_stream/screenmgr"
	"screen_stream/util"
	cfg "screen_stream/util/config"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {return true},
	}
)

type Server struct {
	display *screenmgr.Display
	cancelChan chan struct{}
	options Options
	log *log.Logger
	config cfg.Config
}

var DefaultOptions Options = Options{sampleRate: 30}

func New(config cfg.Config,log *log.Logger) Server {
	return Server{
		cancelChan:make(chan struct{}),
		display: screenmgr.NewDisplay(0),
		options: DefaultOptions,
		log:log,
		config: config,
	}

}

func (s *Server) WithSampleRate(sampleRate int) *Server {
	s.options.sampleRate = sampleRate
	return s
}

func (s *Server) Stop() {
	s.log.Println("stopping the server")
	close(s.cancelChan)
	
}

// func (s *Server) CheckUsername(uname string) error {
// 	user, err := userlib.Current()
	
// 	if err != nil{
// 		return err		
// 	} else if user.Username != uname{
// 		return fmt.Errorf("invalid username: %s", uname)
// 	}
	
// 	return nil
// }


func (s *Server) SpawnNewStream() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		

		uname := r.Header.Get("uname")
		pass := r.Header.Get("pass")

		if err := util.CompareHash(s.config.Password, pass);err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid username or password"))
			return
		} else if err := util.CompareHash(s.config.Username, uname); err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid username or password"))
			return
		}

		stream := screenmgr.NewStream(s.options.sampleRate, s.display)

		ws, err := upgrader.Upgrade(w, r, nil)
		
		if err != nil {
			w.Write([]byte(err.Error()))
			s.log.Println(err)

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
					ticker.Stop()
					return
				case <- s.cancelChan:
					stream.Stop()
					ticker.Stop()
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
			case <- s.cancelChan:
				stream.Stop()
				return
			case x := <- ch:
				// just sending the received pixels of the image.RGBA object
				if err != nil {
					s.log.Println(err)
					return
				}else if err = ws.WriteMessage(websocket.BinaryMessage, x.Pix); err != nil {
					s.log.Println(err)
					return
				}
			}
		}
	})
}

