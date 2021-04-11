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

import (
	"fmt"
	"image"
	"strings"
)

type Gravity int

const (
	GravityTopLeft Gravity = iota
	GravityTop
	GravityTopRight
	GravityLeft
	GravityCenter
	GravityRight
	GravityBottomLeft
	GravityBottom
	GravityBottomRight
)

func (gr Gravity) String() string {
	b, err := gr.MarshalText()
	if err != nil {
		return "Unknown"
	}

	return string(b)
}

func (gr Gravity) MarshalText() ([]byte, error) {
	switch gr {
	case GravityTopLeft:
		return []byte("top-left"), nil
	case GravityTop:
		return []byte("top"), nil
	case GravityTopRight:
		return []byte("top-right"), nil
	case GravityLeft:
		return []byte("left"), nil
	case GravityCenter:
		return []byte("center"), nil
	case GravityRight:
		return []byte("right"), nil
	case GravityBottomLeft:
		return []byte("bottom-left"), nil
	case GravityBottom:
		return []byte("bottom"), nil
	case GravityBottomRight:
		return []byte("bottom-right"), nil
	}

	return nil, fmt.Errorf("not a valid gravity %d", gr)
}

func (gr *Gravity) UnmarshalText(val []byte) error {
	txt := string(val)

	switch strings.ToLower(txt) {
	case "tl", "lt", "":
		*gr = GravityTopLeft
	case "t", "tc", "ct":
		*gr = GravityTop
	case "tr", "rt":
		*gr = GravityTopRight
	case "l", "cl", "lc":
		*gr = GravityLeft
	case "c", "cc":
		*gr = GravityCenter
	case "r", "cr", "rc":
		*gr = GravityRight
	case "bl", "lb":
		*gr = GravityBottomLeft
	case "b", "bc", "cb":
		*gr = GravityBottom
	case "br", "rb":
		*gr = GravityBottomRight
	default:
		return fmt.Errorf("not a valid gravity %q", txt)
	}

	return nil
}

func (gr Gravity) PointOnPicture(w, h int) (image.Point, error) {
	var pt image.Point

	switch gr {
	case GravityTopLeft:
	case GravityTop:
		pt.X = w / 2
	case GravityTopRight:
		pt.X = w
	case GravityLeft:
		pt.Y = h / 2
	case GravityCenter:
		pt.X = w / 2
		pt.Y = h / 2
	case GravityRight:
		pt.X = w
		pt.Y = h / 2
	case GravityBottomLeft:
		pt.Y = h
	case GravityBottom:
		pt.X = w / 2
		pt.Y = h
	case GravityBottomRight:
		pt.X = w
		pt.Y = h
	default:
		return pt, fmt.Errorf("not a valid gravity %d", gr)
	}

	return pt, nil
}

func (gr Gravity) PointWatermark(w, h, wmW, wmH int) (image.Point, error) {
	var pt image.Point

	switch gr {
	case GravityTopLeft:
	case GravityTop:
		pt.X = (w - wmW) / 2
	case GravityTopRight:
		pt.X = w - wmW
	case GravityLeft:
		pt.Y = (h - wmH) / 2
	case GravityCenter:
		pt.X = (w - wmW) / 2
		pt.Y = (h - wmH) / 2
	case GravityRight:
		pt.X = w - wmW
		pt.Y = (h - wmH) / 2
	case GravityBottomLeft:
		pt.Y = h - wmH
	case GravityBottom:
		pt.X = (w - wmW) / 2
		pt.Y = h - wmH
	case GravityBottomRight:
		pt.X = w - wmW
		pt.Y = h - wmH
	default:
		return pt, fmt.Errorf("not a valid gravity %d", gr)
	}

	return pt, nil
}
