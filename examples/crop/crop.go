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
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s path/to/image/file.jpg WxH gravity\n", os.Args[0])
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(os.Args[1])
	checkErr(err)

	vipsImage, err := vips.Load(b)
	checkErr(err)
	defer vipsImage.Clear()
	defer vips.Shutdown()

	gravity := vips.GravityTopLeft
	checkErr(gravity.UnmarshalText([]byte(os.Args[3])))

	pt, err := gravity.PointOnPicture(vipsImage.Width(), vipsImage.Height())
	checkErr(err)

	dim := strings.Split(os.Args[2], "x")
	if len(dim) != 2 {
		checkErr(fmt.Errorf("invalid width and heigh params for crop a image.\nUsage: %s path/to/image/file.jpg 200x300 c\n", os.Args[0]))
	}

	w, err := strconv.Atoi(dim[0])
	checkErr(err)
	h, err := strconv.Atoi(dim[1])
	checkErr(err)

	pt.X -= w / 2
	pt.Y -= h / 2

	checkErr(vipsImage.Crop(w, h, pt))

	// Resize a big photo before crop. For a small photo (less than 100x100), the first stage is crop then resize
	//checkErr(vipsImage.Fill(w, h, pt))

	imgFmt := vips.FormatByMagicNumber(b)
	buf, err := vipsImage.Save(imgFmt, vips.DefaultEncodeConfig)
	checkErr(err)

	fileName := filepath.Base(os.Args[1])
	checkErr(ioutil.WriteFile("resized-"+fileName, buf, 0644))
}
