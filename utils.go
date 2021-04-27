/*
MIT License

Copyright (c) 2021 MyBack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package libvips_go

/*
#cgo pkg-config: vips
#cgo LDFLAGS: -s -w
#cgo CFLAGS: -O3
#include "vips.h"
*/
import "C"
import (
	"fmt"
	"image"
	"io/ioutil"
	"unsafe"
)

var DefaultEncodeConfig = encodeConfig{
	compression:     C.int(6),
	heifCompression: C.int(1),
	interlace:       C.int(1),
	lossless:        C.gboolean(1),
	palette:         C.int(1),
	quality:         C.int(95),
	strip:           C.gboolean(0),
}

type encodeConfig struct {
	compression, heifCompression, interlace, palette, quality C.int
	lossless, strip                                           C.gboolean
}

func (ec *encodeConfig) Compression(i int) {
	ec.compression = C.int(i)
}

func (ec *encodeConfig) HEIFCompression(i HEIFCompressionType) {
	ec.compression = C.int(i)
}

func (ec *encodeConfig) Interlace(b bool) {
	ec.interlace = boolToCInt(b)
}

func (ec *encodeConfig) Lossless(b bool) {
	ec.lossless = gbool(b)
}

func (ec *encodeConfig) Palette(b bool) {
	ec.palette = boolToCInt(b)
}

func (ec *encodeConfig) Quality(i int) {
	ec.quality = C.int(i)
}

func (ec *encodeConfig) StripMetadata(b bool) {
	ec.strip = gbool(b)
}

func boolToCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func vipsError() error {
	defer C.vips_error_clear()
	return fmt.Errorf(C.GoString(C.vips_error_buffer()))
}

func Load(buf []byte) (*VipsImage, error) {
	imgType := FormatByMagicNumber(buf)
	if imgType == Unknown {
		return nil, ErrUnsupportedImageFormat
	}

	img := &VipsImage{}

	err := C.int(0)

	if err = C.vips_image_load_go(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(imgType), &img.img); err != 0 {
		return nil, vipsError()
	}

	return img, nil
}

func LoadFromFile(file string) (*VipsImage, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return Load(b)
}

func LoadPDFPages(buf []byte, page, num int) (*VipsImage, error) {
	imgType := FormatByMagicNumber(buf)
	if imgType != PDF {
		return nil, ErrInvalidFileFormat
	}

	img := &VipsImage{}

	err := C.int(0)

	if err = C.vips_pdf_load_go(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), &img.img, C.int(page), C.int(num)); err != 0 {
		return nil, vipsError()
	}

	return img, nil
}

func LoadFromImage(data image.Image) *VipsImage {
	bounds := data.Bounds()

	imgSize := bounds.Dx() * bounds.Dy() * 4
	imgData := make([]byte, 0, imgSize)
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, a := data.At(i, j).RGBA()
			imgData = append(imgData, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	return &VipsImage{
		img: C.vips_image_new_from_bytes_go(unsafe.Pointer(&imgData[0]), C.size_t(imgSize), C.int(bounds.Dx()), C.int(bounds.Dy())),
	}
}

func Join(in []*VipsImage) (*VipsImage, error) {
	var tmp *C.VipsImage

	arr := make([]*C.VipsImage, len(in))
	for i, im := range in {
		arr[i] = im.img
	}

	if C.vips_arrayjoin_go(&arr[0], &tmp, C.int(len(arr))) != 0 {
		return nil, vipsError()
	}

	return &VipsImage{tmp}, nil
}

func Pixel() *VipsImage {
	imgSize := 4

	imgData := make([]byte, 0, imgSize)
	imgData = append(imgData, 255, 255, 255, 255)
	return &VipsImage{
		img: C.vips_image_new_from_bytes_go(unsafe.Pointer(&imgData[0]), C.size_t(imgSize), C.int(1), C.int(1)),
	}
}

func Cleanup() {
	C.vips_error_clear()
	C.vips_thread_shutdown()
}

func Shutdown() {
	C.vips_shutdown()
}
