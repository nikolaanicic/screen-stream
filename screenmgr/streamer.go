package screenmgr

import (
	"fmt"
	"image"
	"time"
)

type DisplayStream struct {
	sampleRate 	int
	display *Display
	imgchan chan *image.RGBA
	closechan chan struct{}
}

func NewStream(sampleRate int,d *Display) *DisplayStream{

	return &DisplayStream{
		sampleRate: sampleRate,
		display: d,
		imgchan: make(chan *image.RGBA, 1),
		closechan: make(chan struct{}),
		
	}
}


func (s *DisplayStream) Stop(){
	fmt.Println("closing the stream")
	s.close()
}

func (s *DisplayStream) close(){
	close(s.closechan)
}

func(s *DisplayStream) Wait() <-chan struct{}{
	return s.closechan
}


func (s *DisplayStream) Start() chan *image.RGBA{
	fmt.Println("starting the stream")

	samplePeriod := float64(1000) / float64(s.sampleRate)
	sampleTicker := time.NewTicker(time.Millisecond * time.Duration(samplePeriod))

	var img *image.RGBA
	go func(){
		for{
			select{
			case <-s.closechan:
				sampleTicker.Stop()
				close(s.imgchan)
				return
				
			case <-sampleTicker.C:
				img, _ = s.display.Capture()
				s.imgchan <- img
		}
	}
	}()

	return s.imgchan
}


