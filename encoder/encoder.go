package encoder

import (
	"bytes"
	"image"
	"image/jpeg"
)

type Encoder struct{}

func (e *Encoder) BytesToJpeg(m image.Image) ([]byte,error) {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf,m,nil);err != nil{
		return nil, err
	}

	return buf.Bytes(), nil
}


func (e *Encoder) JpegToImage(i []byte)(image.Image, error){
	return jpeg.Decode(bytes.NewReader(i))
}


func NewEncoder() *Encoder{
	return &Encoder{}
}



