/*
Based on the source code https://github.com/imgproxy/imgproxy


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
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"runtime"
	"strconv"
	"unsafe"
)

const Version = string(C.VIPS_VERSION)

var (
	ErrUnsupportedImageFormat = fmt.Errorf("unsupported image file format")
	ErrInvalidFileFormat      = fmt.Errorf("invalid file format")
)

func init() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := C.vips_initialize_go(); err != 0 {
		Shutdown()
		fmt.Println("failed to initialize vips library:", err)
		os.Exit(1)
	}

	// Disable libvips cache. Since processing pipeline is fine tuned, we won't get much profit from it.
	// Enabled cache can cause SIGSEGV on Musl-based systems like Alpine.

	var cacheMaxMem uint64
	var cacheMax int
	var vectorEnabled, leakCheck bool
	var err error

	cacheMaxMemVal := os.Getenv("VIPS_CACHE_MAX_MEM")
	if cacheMaxMemVal != "" {
		if cacheMaxMem, err = strconv.ParseUint(cacheMaxMemVal, 10, 64); err != nil {
			fmt.Printf("failed to initialize vips library: invalid VIPS_CACHE_MAX_MEM value \"%s\": %s\n", cacheMaxMemVal, err)
			os.Exit(1)
		}
	}

	cacheMaxVal := os.Getenv("VIPS_CACHE_MAX")
	if cacheMaxVal != "" {
		if cacheMax, err = strconv.Atoi(cacheMaxVal); err != nil {
			fmt.Printf("failed to initialize vips library: invalid VIPS_CACHE_MAX value \"%s\": %s\n", cacheMaxVal, err)
			os.Exit(1)
		}
	}

	if len(os.Getenv("VIPS_CACHE_TRACE")) > 0 {
		C.vips_cache_set_trace(C.gboolean(1))
	}

	C.vips_cache_set_max_mem(C.size_t(cacheMaxMem))
	C.vips_cache_set_max(C.int(cacheMax))

	C.vips_concurrency_set(1)

	vectorEnabledVal := os.Getenv("VIPS_VECTOR_ENABLED")
	if vectorEnabledVal != "" {
		vectorEnabled, err = strconv.ParseBool(vectorEnabledVal)
		if err != nil {
			fmt.Printf("failed to initialize vips library: invalid VIPS_VECTOR_ENABLED value \"%s\": %s\n", vectorEnabledVal, err)
			os.Exit(1)
		}
	}
	// Vector calculations cause SIGSEGV sometimes when working with JPEG.
	// It's better to disable it since profit it quite small
	C.vips_vector_set_enabled(gbool(vectorEnabled))

	leakCheckVal := os.Getenv("VIPS_LEAK_CHECK")
	if leakCheckVal != "" {
		leakCheck, err = strconv.ParseBool(vectorEnabledVal)
		if err != nil {
			fmt.Printf("failed to initialize vips library: invalid VIPS_LEAK_CHECK value \"%s\": %s\n", leakCheckVal, err)
			os.Exit(1)
		}
	}

	C.vips_leak_set(gbool(leakCheck))
}

type VipsImage struct{ img *C.VipsImage }

func (img *VipsImage) Width() int {
	return int(img.img.Xsize)
}

func (img *VipsImage) Height() int {
	return int(img.img.Ysize)
}

func (img *VipsImage) IsAnimated() bool {
	return C.vips_is_animated_go(img.img) > 0
}

func (img *VipsImage) HasAlpha() bool {
	return C.vips_image_hasalpha(img.img) > 0
}

func (img *VipsImage) Clear() {
	if img.img != nil {
		C.clear_image_go(&img.img)
	}
}

func (img *VipsImage) CopyMemory() error {
	var tmp *C.VipsImage
	if tmp = C.vips_image_copy_memory(img.img); tmp == nil {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Save(imgType ImageFormat, opts encodeConfig) ([]byte, error) {
	if imgType == ICO {
		b, err := img.SaveAsIco()
		return b, err
	}

	err := C.int(0)

	var ptr unsafe.Pointer
	defer C.g_free_go(&ptr)

	imgSize := C.size_t(0)
	switch imgType {
	case JPEG:
		err = C.vips_jpegsave_go(img.img, &ptr, &imgSize, opts.quality, opts.strip, opts.interlace)
	case PNG:
		err = C.vips_pngsave_go(img.img, &ptr, &imgSize, opts.compression, opts.strip, opts.interlace, opts.palette)
	case WEBP:
		err = C.vips_webpsave_go(img.img, &ptr, &imgSize, opts.quality, opts.strip, opts.lossless)
	case GIF:
		err = C.vips_gifsave_go(img.img, &ptr, &imgSize)
	case TIFF:
		err = C.vips_tiffsave_go(img.img, &ptr, &imgSize, opts.quality)
	case AVIF:
		err = C.vips_avifsave_go(img.img, &ptr, &imgSize, opts.quality)
	case HEIF:
		err = C.vips_heifsave_go(img.img, &ptr, &imgSize, opts.quality, opts.heifCompression, opts.lossless)
	case BMP:
		err = C.vips_bmpsave_go(img.img, &ptr, &imgSize)
	case PDF:
		err = C.vips_pdfsave_go(img.img, &ptr, &imgSize)
	default:
		return nil, ErrUnsupportedImageFormat
	}
	if err != 0 {
		return nil, vipsError()
	}

	return C.GoBytes(ptr, C.int(imgSize)), nil
}

func (img *VipsImage) SaveAsIco() ([]byte, error) {
	if img.Width() > 256 || img.Height() > 256 {
		return nil, fmt.Errorf("image dimensions is too big. Max dimension size for ICO is 256")
	}

	imgSize := C.size_t(0)

	var ptr unsafe.Pointer
	defer C.g_free_go(&ptr)

	if C.vips_pngsave_go(img.img, &ptr, &imgSize, 0, 0, 256, 1) != 0 {
		return nil, vipsError()
	}

	buf := new(bytes.Buffer)
	buf.Grow(22 + int(imgSize))

	// ICONDIR header
	if _, err := buf.Write([]byte{0, 0, 1, 0, 1, 0}); err != nil {
		return nil, err
	}

	// ICONDIRENTRY
	if _, err := buf.Write([]byte{
		byte(img.Width() % 256),
		byte(img.Height() % 256),
	}); err != nil {
		return nil, err
	}
	// Number of colors. Not supported in our case
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	// Reserved
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	// Color planes. Always 1 in our case
	if _, err := buf.Write([]byte{1, 0}); err != nil {
		return nil, err
	}
	// Bits per pixel
	if img.HasAlpha() {
		if _, err := buf.Write([]byte{32, 0}); err != nil {
			return nil, err
		}
	} else {
		if _, err := buf.Write([]byte{24, 0}); err != nil {
			return nil, err
		}
	}
	// Image data size
	if err := binary.Write(buf, binary.LittleEndian, uint32(imgSize)); err != nil {
		return nil, err
	}
	// Image data offset. Always 22 in our case
	if _, err := buf.Write([]byte{22, 0, 0, 0}); err != nil {
		return nil, err
	}

	if _, err := buf.Write(C.GoBytes(ptr, C.int(imgSize))); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (img *VipsImage) Resize(scale float64, hasAlpha bool) error {
	var tmp *C.VipsImage

	if hasAlpha {
		if C.vips_resize_with_premultiply_go(img.img, &tmp, C.double(scale)) != 0 {
			return vipsError()
		}
	} else {
		if C.vips_resize_go(img.img, &tmp, C.double(scale)) != 0 {
			return vipsError()
		}
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Rotate(angle int) error {
	var tmp *C.VipsImage

	vipsAngle := (angle / 90) % 4

	if C.vips_rotate_go(img.img, &tmp, C.VipsAngle(vipsAngle)) != 0 {
		return vipsError()
	}

	C.vips_autorot_remove_angle(tmp)
	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Flip() error {
	var tmp *C.VipsImage
	if C.vips_flip_horizontal_go(img.img, &tmp) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) EnsureAlpha() error {
	var tmp *C.VipsImage
	if C.vips_ensure_alpha_go(img.img, &tmp) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Blur(sigma float32) error {
	var tmp *C.VipsImage
	if C.vips_gaussblur_go(img.img, &tmp, C.double(sigma)) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Sharpen(sigma float32) error {
	var tmp *C.VipsImage
	if C.vips_sharpen_go(img.img, &tmp, C.double(sigma)) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Trim(threshold float64, smart bool, color color.RGBA, equalHor bool, equalVer bool) error {
	var tmp *C.VipsImage

	if err := img.CopyMemory(); err != nil {
		return err
	}

	if C.vips_trim_go(img.img, &tmp, C.double(threshold),
		gbool(smart), C.double(color.R), C.double(color.G), C.double(color.B),
		gbool(equalHor), gbool(equalVer)) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Extract(out *VipsImage, pt image.Point, w, h int) error {
	if C.vips_extract_area_go(img.img, &out.img, C.int(pt.X), C.int(pt.Y), C.int(w), C.int(h)) != 0 {
		return vipsError()
	}

	return nil
}

func (img *VipsImage) Fill(dstW, dstH int, pt image.Point) error {
	if dstW < 0 || dstH < 0 {
		return fmt.Errorf("dimensions must be a positive values")
	}

	srcW, srcH := img.Width(), img.Height()

	if dstW == 0 || dstH == 0 || (dstW >= srcW && dstH >= srcH) {
		return nil
	}

	scale := 1.0

	if srcW >= 100 && srcH >= 100 {
		if err := img.fillCrop(dstW, dstH, pt); err != nil {
			return err
		}

		if srcW > srcH {
			scale = float64(dstW) / float64(img.Width())
		} else if srcW < srcH {
			scale = float64(dstH) / float64(img.Height())
		}

		return img.Resize(scale, img.HasAlpha())
	} else {
		if srcW > srcH {
			scale = float64(dstW) / float64(srcW)
		} else if srcW < srcH {
			scale = float64(dstH) / float64(srcH)
		}

		if err := img.Resize(scale, img.HasAlpha()); err != nil {
			return err
		}

		return img.fillCrop(dstW, dstH, pt)
	}
}

func (img *VipsImage) fillCrop(dstW, dstH int, pt image.Point) error {
	srcW, srcH := img.Width(), img.Height()

	srcAspectRatio := float64(srcW) / float64(srcH)
	dstAspectRatio := float64(dstW) / float64(dstH)

	var w, h int
	if srcAspectRatio < dstAspectRatio {
		cropH := float64(srcW) * float64(dstH) / float64(dstW)
		w = srcW
		h = int(math.Max(1, cropH) + 0.5)
	} else {
		cropW := float64(srcH) * float64(dstW) / float64(dstH)
		w = int(math.Max(1, cropW) + 0.5)
		h = srcH
	}

	return img.Crop(w, h, pt)
}

func (img *VipsImage) AddWatermark(wm *VipsImage, pt image.Point, opacity float64) error {
	var tmp *C.VipsImage
	if err := C.vips_apply_watermark_go(img.img, wm.img, &tmp, C.int(pt.X), C.int(pt.Y), C.float(opacity)); err != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

// Strip - Remove EXIF data
func (img *VipsImage) Strip() error {
	var tmp *C.VipsImage
	if C.vips_strip_go(img.img, &tmp) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

// VipsJpeg: out of order read at line XXXX
// https://github.com/libvips/libvips/issues/639#issuecomment-294513915
func (img *VipsImage) SmartCrop(w, h int) error {
	var tmp *C.VipsImage
	if C.vips_smartcrop_go(img.img, &tmp, C.int(w), C.int(h)) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) Crop(dstW, dstH int, pt image.Point) error {
	var tmp *C.VipsImage
	if C.vips_extract_area_go(img.img, &tmp, C.int(pt.X), C.int(pt.Y), C.int(dstW), C.int(dstH)) != 0 {
		return vipsError()
	}

	C.swap_and_clear_go(&img.img, tmp)

	return nil
}

func (img *VipsImage) SetInt(key string, val int) {
	C.vips_image_set_int(img.img, C.CString(key), C.int(val))
}
