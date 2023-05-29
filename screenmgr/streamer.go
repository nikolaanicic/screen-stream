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
	cancel context.CancelFunc

}

func NewStream(sampleRate int,d *Display) *DisplayStream{
	ctx, cancel := context.WithCancel(context.Background())

	return &DisplayStream{
		sampleRate: sampleRate,
		ctx:ctx,
		display: d,
		imgchan: make(chan *image.RGBA, 1),
		cancel: cancel,
	}
}


func (s *DisplayStream) Stop(){
	fmt.Println("closing the stream")
	s.cancel()
}

func(s *DisplayStream) Wait() <-chan struct{}{
	return s.ctx.Done()
}


func (s *DisplayStream) Start() chan *image.RGBA{
	fmt.Println("starting the stream")

	samplePeriod := float64(1000) / float64(s.sampleRate)
	sampleTicker := time.NewTicker(time.Millisecond * time.Duration(samplePeriod))

	var img *image.RGBA
	go func(){
		for{
			select{
			case <-s.ctx.Done():
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


