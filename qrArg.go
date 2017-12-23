package qrcodee_xtend

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"
	"strconv"

	qrcode "github.com/huncent/go-qrcode"
)

type QRArg struct {
	Content   string
	bianhao   string
	size      int
	bgcolor   color.Color
	forecolor color.Color
	logo      image.Image
	level     qrcode.RecoveryLevel
	bgimg     image.Image
	bdmaxsize int
}

type Values interface {
	Get(string) string
}

func (q *QRArg) Parse(query Values) {
	q.Content = query.Get("content")
	q.size = q.parseSize(query.Get("size"))
	q.bianhao = query.Get("bianhao")
	q.bgcolor = q.parseBGColor(query.Get("bgcolor"))
	q.forecolor = q.parseForeColor(query.Get("forecolor"))
	q.logo = q.parseLogo(query.Get("logo"))
	q.bgimg = q.parseBGImg(query.Get("bgimg"))
	q.bdmaxsize = q.parseBdmaxsize(query.Get("bdmaxsize"))
	q.level = qrcode.Medium
	if q.logo == nil {
		q.level = qrcode.Highest
	}
	if q.bgimg != nil {
		if q.bgimg.Bounds().Max.X > q.bgimg.Bounds().Max.Y {
			q.size = q.bgimg.Bounds().Max.Y
		} else {
			q.size = q.bgimg.Bounds().Max.X
		}
		//		q.level = qrcode.Highest
	}
}

func (q *QRArg) parseSize(str string) int {
	size := 256
	if str != "" {
		s, err := strconv.Atoi(str)
		if err != nil {
			size = 256
		}
		size = s
	}
	return size
}

func (q *QRArg) parseBdmaxsize(str string) int {
	size := -1
	if str != "" {
		s, err := strconv.Atoi(str)
		if err != nil {
			size = -1
		}
		size = s
	}
	return size
}

func (q *QRArg) parseBGColor(str string) color.Color {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		return color.White
	}
	return color.RGBA{R: uint8(s & 0xff0000 >> 16),
		G: uint8(s & 0xff00 >> 8),
		B: uint8(s & 0xff),
		A: uint8(0xff)}
}

func (q *QRArg) parseForeColor(str string) color.Color {
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		return color.Black
	}
	return color.RGBA{R: uint8(s & 0xff0000 >> 16),
		G: uint8(s & 0xff00 >> 8),
		B: uint8(s & 0xff),
		A: uint8(uint8(0xff))}
}

func (q *QRArg) parseLogo(str string) image.Image {
	if len(str) == 0 {
		return nil
	}
	return q.downImg(str)
}

func (q *QRArg) parseBGImg(str string) image.Image {
	if len(str) == 0 {
		return nil
	}
	return q.downImg(str)
}

func (q *QRArg) downImg(str string) image.Image {
	resp, err := http.Get(str)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	logo, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil
	}
	return logo
}
