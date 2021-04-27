package main

import (
	"fmt"
	"io/ioutil"
	"os"

	vips "github.com/myback/libvips-go"
)

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func main() {
	defer vips.Shutdown()

	vipsImage := vips.Pixel()
	defer vipsImage.Clear()

	enc := vips.DefaultEncodeConfig
	enc.Lossless(false)
	enc.HEIFCompression(vips.HEIF_COMPRESSION_HEVC)
	enc.Interlace(false)

	for _, t := range []vips.ImageFormat{vips.AVIF, vips.BMP, vips.GIF, vips.HEIF, vips.ICO, vips.JPEG, vips.PDF,
		vips.PNG, vips.TIFF, vips.WEBP} {

		buf, err := vipsImage.Save(t, enc)
		checkErr(err)

		checkErr(ioutil.WriteFile(".test/blank"+t.Ext(), buf, 0644))
	}
}
