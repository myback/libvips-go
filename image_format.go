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
	"strings"
)

type ImageFormat int
type HEIFCompressionType int

const (
	Unknown = ImageFormat(C.UNKNOWN)
	AVIF    = ImageFormat(C.AVIF)
	BMP     = ImageFormat(C.BMP)
	GIF     = ImageFormat(C.GIF)
	HEIF    = ImageFormat(C.HEIF)
	ICO     = ImageFormat(C.ICO)
	JPEG    = ImageFormat(C.JPEG)
	PDF     = ImageFormat(C.PDF)
	PNG     = ImageFormat(C.PNG)
	SVG     = ImageFormat(C.SVG)
	TIFF    = ImageFormat(C.TIFF)
	WEBP    = ImageFormat(C.WEBP)

	HEIF_COMPRESSION_HEVC = HEIFCompressionType(C.VIPS_FOREIGN_HEIF_COMPRESSION_HEVC)
	HEIF_COMPRESSION_AVC  = HEIFCompressionType(C.VIPS_FOREIGN_HEIF_COMPRESSION_AVC)
	HEIF_COMPRESSION_JPEG = HEIFCompressionType(C.VIPS_FOREIGN_HEIF_COMPRESSION_JPEG)
)

func (imgFmt ImageFormat) Ext() string {
	b, err := imgFmt.MarshalText()
	if err != nil {
		return ""
	}

	return "." + string(b)
}

func (imgFmt ImageFormat) String() string {
	b, err := imgFmt.MarshalText()
	if err != nil {
		return "Unknown"
	}

	return string(b)
}

func (imgFmt ImageFormat) MarshalText() ([]byte, error) {
	switch imgFmt {
	case AVIF:
		return []byte("avif"), nil
	case BMP:
		return []byte("bmp"), nil
	case GIF:
		return []byte("gif"), nil
	case HEIF:
		return []byte("heif"), nil
	case ICO:
		return []byte("ico"), nil
	case JPEG:
		return []byte("jpg"), nil
	case PDF:
		return []byte("pdf"), nil
	case PNG:
		return []byte("png"), nil
	case SVG:
		return []byte("svg"), nil
	case TIFF:
		return []byte("tiff"), nil
	case WEBP:
		return []byte("webp"), nil
	}

	return nil, fmt.Errorf("unknown image format")
}

func (imgFmt *ImageFormat) UnmarshalText(val []byte) error {
	var ext string
	if val[0] == '.' {
		ext = string(val[1:])
	} else {
		ext = string(val)
	}

	switch strings.ToLower(ext) {
	case "avif":
		*imgFmt = AVIF
	case "bmp":
		*imgFmt = BMP
	case "gif":
		*imgFmt = GIF
	case "heic", "heif":
		*imgFmt = HEIF
	case "ico":
		*imgFmt = ICO
	case "jpg", "jpeg":
		*imgFmt = JPEG
	case "pdf":
		*imgFmt = PDF
	case "png":
		*imgFmt = PNG
	case "svg":
		*imgFmt = SVG
	case "tiff":
		*imgFmt = TIFF
	case "webp":
		*imgFmt = WEBP
	default:
		*imgFmt = Unknown
	}

	return nil
}

func FormatByMagicNumber(buf []byte) ImageFormat {
	if buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF {
		return JPEG
	}

	if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		return PNG
	}

	if buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46 {
		return GIF
	}

	if buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50 {
		return WEBP
	}

	if buf[0] == 0x25 && buf[1] == 0x50 && buf[2] == 0x44 && buf[3] == 0x46 && buf[4] == 0x2D {
		return PDF
	}

	if buf[0] == 0x42 && buf[1] == 0x4D {
		return BMP
	}

	if buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x01 && buf[3] == 0x00 {
		return ICO
	}

	if (buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x0) ||
		(buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x0 && buf[3] == 0x2A) {
		return TIFF
	}

	//like MPEG-4 video files
	if buf[3] == 0x18 && buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70 {
		return HEIF
	}

	if buf[3] == 0x1C && buf[4] == 0x66 && buf[5] == 0x74 && buf[6] == 0x79 && buf[7] == 0x70 {
		return AVIF
	}

	return Unknown
}
