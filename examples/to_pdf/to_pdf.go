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

	var arrVips []*vips.VipsImage
	for i := 1; i <= 3; i++ {
		i, err := vips.LoadFromFile(".test/wiki_a4.png")
		checkErr(err)

		arrVips = append(arrVips, i)
	}

	defer func() {
		for _, img := range arrVips {
			img.Clear()
		}
	}()

	jVips, err := vips.Join(arrVips)
	checkErr(err)
	defer jVips.Clear()

	jVips.SetInt("page-height", jVips.Height()/len(arrVips))

	buf, err := jVips.Save(vips.PDF, vips.DefaultEncodeConfig)
	checkErr(err)

	checkErr(ioutil.WriteFile("a4-2.pdf", buf, 0644))
}
