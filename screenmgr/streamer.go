package screenmgr

import (
	"context"
	"fmt"
	"image"
	"time"
)

type DisplayStream struct {
	sampleRate 	int
	ctx context.Context
	display *Display
	imgchan chan *image.RGBA

}

func NewStream(sampleRate int, ctx context.Context,d *Display) *DisplayStream{
	return &DisplayStream{
		sampleRate: sampleRate,
		ctx:ctx,
		display: d,
		imgchan: make(chan *image.RGBA, 1),
	}
}


func (s *DisplayStream) stop(){
	close(s.imgchan)
}


func (s *DisplayStream) Start() chan *image.RGBA{

	samplePeriod := float64(1000) / float64(s.sampleRate)
	sampleTicker := time.NewTicker(time.Millisecond * time.Duration(samplePeriod))


	go func(){
		for{
			select{
			case <-s.ctx.Done():
				fmt.Println("closing the stream")
				sampleTicker.Stop()
				s.stop()
				
				return
			case <-sampleTicker.C:
				img, _ := s.display.Capture()
				s.imgchan <- img
		}
	}
	}()

	return s.imgchan
}


