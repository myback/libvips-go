package main

import (
	"fmt"
	vips "github.com/myback/libvips-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s path/to/image/file.jpg path/to/image/watermark.jpg [ WmarkWxWmarkH gravity ]\n", os.Args[0])
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(os.Args[1])
	checkErr(err)

	wm, err := ioutil.ReadFile(os.Args[2])
	checkErr(err)

	vipsImage, err := vips.Load(b)
	checkErr(err)
	defer vipsImage.Clear()

	vipsWmImage, err := vips.Load(wm)
	checkErr(err)
	defer vipsWmImage.Clear()
	defer vips.Shutdown()

	gravity := vips.GravityTopLeft
	wmW := vipsWmImage.Width()
	wmH := vipsWmImage.Height()

	if len(os.Args) == 5 {
		wmDim := strings.Split(os.Args[3], "x")
		if len(wmDim) != 2 {
			checkErr(fmt.Errorf("invalid width and heigh params for crop a image.\nUsage: %s path/to/image/file.jpg path/to/image/watermark.jpg 600x400 70x30 c\n", os.Args[0]))
		}

		wmW, err = strconv.Atoi(wmDim[0])
		checkErr(err)
		wmH, err = strconv.Atoi(wmDim[1])
		checkErr(err)

		checkErr(gravity.UnmarshalText([]byte(os.Args[4])))
	}

	pt, err := gravity.PointOnPicture(vipsImage.Width(), vipsImage.Height())
	checkErr(err)

	pt.X -= wmW / 2
	pt.Y -= wmH / 2

	if pt.X < 0 {
		pt.X = 0
	}

	if pt.Y < 0 {
		pt.Y = 0
	}

	checkErr(vipsImage.AddWatermark(vipsWmImage, pt, 1))

	imgFmt := vips.FormatByMagicNumber(b)
	buf, err := vipsImage.Save(imgFmt, vips.DefaultEncodeConfig)
	checkErr(err)

	fileName := filepath.Base(os.Args[1])
	checkErr(ioutil.WriteFile("watermark-"+fileName, buf, 0644))
}
