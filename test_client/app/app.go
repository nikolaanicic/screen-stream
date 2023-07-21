package app

import (
	"context"
	"image"
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
	img *image.RGBA
	cimg *canvas.Image
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

	a.img = image.NewRGBA(image.Rect(0,0,2560,1080))
	a.cimg = canvas.NewImageFromImage(a.img)

	a.cimg.Image = a.img

	a.card.SetContent(a.cimg)
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
				a.img.Pix = []uint8(x)

				<-ticker.C
				a.cimg.Refresh()
				a.card.Refresh()
			}
		}
		
}


func (a *App) Start(addr string) error{
	c, err := a.c.Connect(addr,"Nikola","anicic")
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

