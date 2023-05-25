package app

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"log"
	"screen_stream/test_client/webclient"
	"time"

	"fyne.io/fyne/v2"
	fa "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)


var(
	DefaultWidth float32 = 840
	DefaultHeight float32 = 640
)

type App struct {
	ctx context.Context
	cancel context.CancelFunc
	app fyne.App
	window fyne.Window
	card *widget.Card
	log *log.Logger
	c *webclient.Client
}


func NewApp(title string, l *log.Logger) *App{


	a := &App{}
	a.ctx,a.cancel = context.WithCancel(context.Background())

	a.app = fa.New()
	a.window = a.app.NewWindow(title)
	a.window.CenterOnScreen()

	a.SetSize(DefaultWidth,DefaultHeight)

	a.window.SetOnClosed(a.onClose)


	a.log = l
	a.c = webclient.NewClient(a.ctx, a.log)
	
	a.card = widget.NewCard("stream","",nil)
	a.window.SetContent(a.card)
	
	return a
}



func (a *App) SetSize(width, height float32){
	a.window.Resize(fyne.NewSize(width,height))
}

// render loop accepts image bytes from the client channel
// puts an image on screen every 1000/30 milliseconds
func (a *App) renderLoop(c chan []byte){
	samplePeriod := float64(1000) / float64(30)
	ticker := time.NewTicker(time.Millisecond * time.Duration(samplePeriod))

		for {
			select {
			case <-a.ctx.Done():
				return
			case x:= <-c:
				img, err := jpeg.Decode(bytes.NewReader(x))
				if err != nil {
					fmt.Println(err)
					continue
				}

				<-ticker.C
				a.card.SetContent(canvas.NewImageFromImage(img))
				a.card.Refresh()
			}
		}
		
}


func (a *App) Start(addr string) error{
	c, err := a.c.Connect(addr)
	if err != nil{
		return err
	}

	go a.renderLoop(c)
	a.window.ShowAndRun()

	return nil
}

func (a *App) onClose(){
	a.cancel()
}

