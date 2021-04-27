package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s path/to/image/file.jpg png\n", os.Args[0])
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(os.Args[1])
	checkErr(err)

	vipsImage, err := vips.Load(b)
	checkErr(err)
	defer vipsImage.Clear()

	imageSaveOptions := vips.DefaultEncodeConfig
	imageSaveOptions.Lossless(false)

	imgFmt := vips.Unknown
	checkErr(imgFmt.UnmarshalText([]byte(os.Args[2])))

	buf, err := vipsImage.Save(imgFmt, imageSaveOptions)
	checkErr(err)

	fileName := filepath.Base(os.Args[1])
	checkErr(ioutil.WriteFile(fileName+imgFmt.Ext(), buf, 0644))
}
