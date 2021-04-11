package libvips_go

import (
	"image"
	"reflect"
	"testing"
)

func TestGravity_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		gr      Gravity
		want    []byte
		wantErr bool
	}{
		{"GravityTopLeft", GravityTopLeft, []byte("top-left"), false},
		{"GravityTop", GravityTop, []byte("top"), false},
		{"GravityTopRight", GravityTopRight, []byte("top-right"), false},
		{"GravityLeft", GravityLeft, []byte("left"), false},
		{"GravityCenter", GravityCenter, []byte("center"), false},
		{"GravityRight", GravityRight, []byte("right"), false},
		{"GravityBottomLeft", GravityBottomLeft, []byte("bottom-left"), false},
		{"GravityBottom", GravityBottom, []byte("bottom"), false},
		{"GravityBottomRight", GravityBottomRight, []byte("bottom-right"), false},
		{"GravityInvalidValue", Gravity(42), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.gr.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGravity_PointOnPicture(t *testing.T) {
	type args struct {
		w int
		h int
	}
	tests := []struct {
		name    string
		gr      Gravity
		args    args
		want    image.Point
		wantErr bool
	}{
		{"PointOnPictureGravityTopLeft", GravityTopLeft, args{100, 100}, image.Point{0, 0}, false},
		{"PointOnPictureGravityTop", GravityTop, args{100, 100}, image.Point{50, 0}, false},
		{"PointOnPictureGravityTopRight", GravityTopRight, args{100, 100}, image.Point{100, 0}, false},
		{"PointOnPictureGravityLeft", GravityLeft, args{100, 100}, image.Point{0, 50}, false},
		{"PointOnPictureGravityCenter", GravityCenter, args{100, 100}, image.Point{50, 50}, false},
		{"PointOnPictureGravityRight", GravityRight, args{100, 100}, image.Point{100, 50}, false},
		{"PointOnPictureGravityBottomLeft", GravityBottomLeft, args{100, 100}, image.Point{0, 100}, false},
		{"PointOnPictureGravityBottom", GravityBottom, args{100, 100}, image.Point{50, 100}, false},
		{"PointOnPictureGravityBottomRight", GravityBottomRight, args{100, 100}, image.Point{100, 100}, false},
		{"PointOnPictureGravityInvalidValue", Gravity(42), args{100, 100}, image.Point{0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.gr.PointOnPicture(tt.args.w, tt.args.h)
			if (err != nil) != tt.wantErr {
				t.Errorf("PointOnPicture() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PointOnPicture() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGravity_PointWatermark(t *testing.T) {
	type args struct {
		w   int
		h   int
		wmW int
		wmH int
	}
	tests := []struct {
		name string
		gr   Gravity
		args args
		want image.Point
		wantErr bool
	}{
		{"PointWatermarkGravityTopLeft", GravityTopLeft, args{100, 100, 15, 10}, image.Point{0, 0}, false},
		{"PointWatermarkGravityTop", GravityTop, args{100, 100, 15, 10}, image.Point{42, 0}, false},
		{"PointWatermarkGravityTopRight", GravityTopRight, args{100, 100, 15, 10}, image.Point{85, 0}, false},
		{"PointWatermarkGravityLeft", GravityLeft, args{100, 100, 15, 10}, image.Point{0, 45}, false},
		{"PointWatermarkGravityCenter", GravityCenter, args{100, 100, 15, 10}, image.Point{42, 45}, false},
		{"PointWatermarkGravityRight", GravityRight, args{100, 100, 15, 10}, image.Point{85, 45}, false},
		{"PointWatermarkGravityBottomLeft", GravityBottomLeft, args{100, 100, 15, 10}, image.Point{0, 90}, false},
		{"PointWatermarkGravityBottom", GravityBottom, args{100, 100, 15, 10}, image.Point{42, 90}, false},
		{"PointWatermarkGravityBottomRight", GravityBottomRight, args{100, 100, 15, 10}, image.Point{85, 90}, false},
		{"PointWatermarkGravityInvalidValue", Gravity(42), args{100, 100, 15, 10}, image.Point{0, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.gr.PointWatermark(tt.args.w, tt.args.h, tt.args.wmW, tt.args.wmH)
			if (err != nil) != tt.wantErr {
				t.Errorf("PointOnPicture() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PointWatermark() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGravity_String(t *testing.T) {
	tests := []struct {
		name string
		gr   Gravity
		want string
	}{
		{"StringGravityTopLeft", GravityTopLeft, "top-left"},
		{"StringGravityTop", GravityTop, "top"},
		{"StringGravityTopRight", GravityTopRight, "top-right"},
		{"StringGravityLeft", GravityLeft, "left"},
		{"StringGravityCenter", GravityCenter, "center"},
		{"StringGravityRight", GravityRight, "right"},
		{"StringGravityBottomLeft", GravityBottomLeft, "bottom-left"},
		{"StringGravityBottom", GravityBottom, "bottom"},
		{"StringGravityBottomRight", GravityBottomRight, "bottom-right"},
		{"StringGravityInvalidValue", Gravity(42), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gr.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGravity_UnmarshalText(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		gr      Gravity
		args    args
		want    Gravity
		wantErr bool
	}{
		{"UnmarshalTextGravityTopLeft", Gravity(42), args{val: []byte("")}, GravityTopLeft, false},
		{"UnmarshalTextGravityTopLeft", Gravity(42), args{val: []byte("tl")}, GravityTopLeft, false},
		{"UnmarshalTextGravityTopLeft", Gravity(42), args{val: []byte("lt")}, GravityTopLeft, false},
		{"UnmarshalTextGravityTop", Gravity(42), args{val: []byte("t")}, GravityTop, false},
		{"UnmarshalTextGravityTop", Gravity(42), args{val: []byte("tc")}, GravityTop, false},
		{"UnmarshalTextGravityTop", Gravity(42), args{val: []byte("ct")}, GravityTop, false},
		{"UnmarshalTextGravityTopRight", Gravity(42), args{val: []byte("tr")}, GravityTopRight, false},
		{"UnmarshalTextGravityTopRight", Gravity(42), args{val: []byte("rt")}, GravityTopRight, false},
		{"UnmarshalTextGravityLeft", Gravity(42), args{val: []byte("l")}, GravityLeft, false},
		{"UnmarshalTextGravityLeft", Gravity(42), args{val: []byte("cl")}, GravityLeft, false},
		{"UnmarshalTextGravityLeft", Gravity(42), args{val: []byte("lc")}, GravityLeft, false},
		{"UnmarshalTextGravityCenter", Gravity(42), args{val: []byte("c")}, GravityCenter, false},
		{"UnmarshalTextGravityCenter", Gravity(42), args{val: []byte("cc")}, GravityCenter, false},
		{"UnmarshalTextGravityRight", Gravity(42), args{val: []byte("r")}, GravityRight, false},
		{"UnmarshalTextGravityRight", Gravity(42), args{val: []byte("cr")}, GravityRight, false},
		{"UnmarshalTextGravityRight", Gravity(42), args{val: []byte("rc")}, GravityRight, false},
		{"UnmarshalTextGravityBottomLeft", Gravity(42), args{val: []byte("bl")}, GravityBottomLeft, false},
		{"UnmarshalTextGravityBottomLeft", Gravity(42), args{val: []byte("lb")}, GravityBottomLeft, false},
		{"UnmarshalTextGravityBottom", Gravity(42), args{val: []byte("b")}, GravityBottom, false},
		{"UnmarshalTextGravityBottom", Gravity(42), args{val: []byte("bc")}, GravityBottom, false},
		{"UnmarshalTextGravityBottom", Gravity(42), args{val: []byte("cb")}, GravityBottom, false},
		{"UnmarshalTextGravityBottomRight", Gravity(42), args{val: []byte("br")}, GravityBottomRight, false},
		{"UnmarshalTextGravityBottomRight", Gravity(42), args{val: []byte("rb")}, GravityBottomRight, false},
		{"UnmarshalTextInvalidValue", Gravity(42), args{val: []byte("lr")}, GravityTopLeft, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.gr.UnmarshalText(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
