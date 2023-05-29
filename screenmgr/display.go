package screenmgr

import (
	"image"
	"image/jpeg"
	"io"

	"github.com/kbinani/screenshot"
)

type Display struct {
	width      	int
	height     	int
	displayNum 	int
}


func getWidthHeight(index int) (width,height int){
	bounds := screenshot.GetDisplayBounds(index).Size()
	width, height = bounds.X, bounds.Y

	return
}

func NewDisplay(index int) *Display{

	w,h := getWidthHeight(index)

	display := &Display{
		displayNum: index,
		width: w,
		height: h,
	}

	return display
}

func (d *Display) GetSize() (width, height int){
	return d.width, d.height
}

func (d *Display) Capture() (*image.RGBA, error){
	return screenshot.CaptureDisplay(d.displayNum)
}


func (d *Display) Sample(output io.Writer) error {
	img, err := d.Capture()

	if err != nil{
		return err
	}else if err := jpeg.Encode(output,img,nil); err != nil{
		return err
	}
	return nil
}
