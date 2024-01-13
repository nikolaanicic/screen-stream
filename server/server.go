package server

import (
	"log"
	"net/http"
	"screen_stream/eventsmgr"
	"screen_stream/screenmgr"
	"screen_stream/util"
	cfg "screen_stream/util/config"
	"time"

	"github.com/gorilla/websocket"
	hook "github.com/robotn/gohook"
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
	eventManager *eventsmgr.EventManager
}

var DefaultOptions Options = Options{sampleRate: 30}

func New(config cfg.Config,log *log.Logger, disp *screenmgr.Display) Server {
	return Server{
		cancelChan:make(chan struct{}),
		display: disp,
		options: DefaultOptions,
		log:log,
		config: config,
		eventManager: &eventsmgr.EventManager{},
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

func (s *Server) checkCredentials(r *http.Request) error{

	uname := r.Header.Get("uname")
	pass := r.Header.Get("pass")

	if err := util.CompareHash(s.config.Password, pass);err != nil {
		return err
	} else if err := util.CompareHash(s.config.Username, uname); err != nil{
		return err
	}

	return nil
}


func (s *Server) SpawnNewEventsHandler() func(http.ResponseWriter, *http.Request){
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := s.checkCredentials(r); err != nil{
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid username or password"))
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil{
			w.Write([]byte(err.Error()))
			s.log.Println(err)

			return
		}

		defer ws.Close()

		ws.SetCloseHandler(func(code int, text string) error {
			return ws.WriteControl(websocket.CloseMessage,[]byte{},time.Now().Add(time.Second))
		})

		var e hook.Event
		for{
			select{
			case <- s.cancelChan:
				return
			default:
				err := ws.ReadJSON(&e)
				if err != nil{
					s.log.Println(err)
					return
				}
				
				s.eventManager.HandleEvent(&e)
			}
		}

	})
}

func (s *Server) SpawnNewScreenStream() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if err := s.checkCredentials(r); err != nil{
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

		ws.SetCloseHandler(func(code int, text string) error {
			stream.Stop()
			return ws.WriteControl(websocket.CloseMessage,[]byte{},time.Now().Add(time.Second))
		})
	
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

