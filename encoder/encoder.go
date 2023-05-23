package encoder

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
)

type Encoder struct{}

func (e *Encoder) BytesToBase64(m image.Image) (res string,err error) {
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf,m,nil)
	
	if err != nil{
		return
	}

	res = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}

func NewEncoder() *Encoder{
	return &Encoder{}
}



