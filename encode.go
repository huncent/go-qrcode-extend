package qrcodee_xtend

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"

	"github.com/disintegration/imaging"

	"github.com/golang/freetype"
	qrcode "github.com/huncent/go-qrcode"
)

type QRDiy struct {
	Arg QRArg
}

func addLabelC(img *image.RGBA, x, y int, label string) {
	fmt.Println("addLabelC")
	data, err := ioutil.ReadFile("/opt/rsa/fonts/Arial.ttf")
	if err != nil {
		fmt.Errorf("error", err)
	}
	f, err := freetype.ParseFont(data)
	if err != nil {
		fmt.Errorf("error", err)
	}
	c := freetype.NewContext()
	c.SetSrc(image.Black)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetFont(f)
	c.SetFontSize(20)
	c.DrawString(label, freetype.Pt(x, y))
}
func (q *QRDiy) Encode() ([]byte, error) {
	var code *qrcode.QRCode
	code, err := qrcode.New(q.Arg.Content, q.Arg.level)
	if err != nil {
		return nil, err
	}
	code.BackgroundColor = q.Arg.bgcolor
	code.ForegroundColor = q.Arg.forecolor
	var img image.Image
	if q.Arg.bdmaxsize < 0 {
		img = code.Image(q.Arg.size)
	} else {
		//img = code.ImageWithBorderMaxSize(q.Arg.size, q.Arg.bdmaxsize)
		img = code.ImageWithBorderMaxSize(q.Arg.size, 0)
	}

	if q.Arg.bgimg != nil {
		q.embgimg(img, q.Arg.bgimg)
	}
	if q.Arg.logo != nil {
		logosize := q.Arg.size / 5
		logo := imaging.Resize(q.Arg.logo, logosize, logosize, imaging.Lanczos)
		q.emlogo(img, logo)
	}
	qsize := q.Arg.size
	lsize := qsize - q.Arg.bdmaxsize*2
	img = imaging.Resize(img, lsize, lsize, imaging.Lanczos)
	plusHeight := 0 //增加高度放编号
	if q.Arg.bianhao != "" {
		plusHeight = 40
	}
	rect := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{qsize, qsize + plusHeight}}
	imgBack := image.NewRGBA(rect)

	for i := 0; i < qsize; i++ {
		for j := 0; j < qsize+40; j++ {
			imgBack.Set(i, j, q.Arg.bgcolor)
		}
	}
	q.emlogo(imgBack, img)
	if q.Arg.bianhao != "" {
		addLabelC(imgBack, qsize/2-11*len(q.Arg.bianhao)/2, qsize+25, q.Arg.bianhao)
	}
	//q.emlabel(imgBack, "010120001")
	//img = imaging.Fill(img, qsize, qsize, imaging.Center, imaging.Lanczos)
	var b bytes.Buffer
	err = png.Encode(&b, imgBack)
	if err != nil {
		return nil, err
	}
	buf := b.Bytes()
	return buf, nil
}
func (q *QRDiy) emlabel(rst image.Image, label string) {
	data, err := ioutil.ReadFile("/Library/Fonts/Arial.ttf")
	if err != nil {
		fmt.Errorf("error", err)
	}
	f, err := freetype.ParseFont(data)
	if err != nil {
		fmt.Errorf("error", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 140, 22))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.Black)
	c.SetFont(f)
	c.SetFontSize(24)
	c.DrawString(label, freetype.Pt(10, 20))
	q.emlogo(rst, dst)
}
func (q *QRDiy) emlogo(rst, logo image.Image) {
	offset := rst.Bounds().Max.X/2 - logo.Bounds().Max.X/2
	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			rst.(*image.RGBA).Set(x+offset, y+offset, logo.At(x, y))
		}
	}
	return
}
func (q *QRDiy) embgimg(rst, bgimg image.Image) {
	if rst.Bounds().Max.X > q.Arg.size {
		return
	}
	qsx, qsy := 0, 0
	br, bg, bb, _ := q.Arg.bgcolor.RGBA()
	fr, fg, fb, _ := q.Arg.forecolor.RGBA()
	qex, qey := 0, 0
	oks, oke := false, false
	for z := 0; z < rst.Bounds().Max.X; z++ {
		cs := rst.(*image.RGBA).At(z, z)
		ce := rst.(*image.RGBA).At(rst.Bounds().Max.X-1-z, z)
		r, g, b, _ := cs.RGBA()
		if r == fr && g == fg && b == fb && !oks {
			qsx, qsy = z, z
			oks = true
		}
		r, g, b, _ = ce.RGBA()
		if r == fr && g == fg && b == fb && !oke {
			qex, qey = rst.Bounds().Max.X-1-z, rst.Bounds().Max.Y-1-z
			oke = true
		}
		if oks && oke {
			break
		}
	}

	for x := 0; x < rst.Bounds().Max.X; x++ {
		for y := 0; y < rst.Bounds().Max.Y; y++ {
			if x < qsx || y < qsy || x > qex || y > qey {
				rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
			} else {
				//				r, g, b, _ := rst.(*image.RGBA).At(x, y).RGBA()
				//				if r == fr && g == fg && b == fb {
				//					rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
				//				}
				r, g, b, _ := rst.(*image.RGBA).At(x, y).RGBA()
				if r == br && g == bg && b == bb {
					rst.(*image.RGBA).Set(x, y, bgimg.At(x, y))
				}
			}
		}
		if x >= qsx-2 && x <= qex+2 {
			rst.(*image.RGBA).Set(x, qsy-1, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qey+1, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qsy-2, q.Arg.bgcolor)
			rst.(*image.RGBA).Set(x, qey+2, q.Arg.bgcolor)
		}
	}

	for y := qsy - 2; y <= qey+2; y++ {
		rst.(*image.RGBA).Set(qsx-1, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qex+1, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qsx-2, y, q.Arg.bgcolor)
		rst.(*image.RGBA).Set(qex+2, y, q.Arg.bgcolor)
	}
	return
}
