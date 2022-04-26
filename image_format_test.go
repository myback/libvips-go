package libvips_go

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestFormatByMagicNumber(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     ImageFormat
	}{
		{"MagicNumberOfAVIF", ".test/blank.avif", AVIF},
		{"MagicNumberOfBMP", ".test/blank.bmp", BMP},
		{"MagicNumberOfGIF", ".test/blank.gif", GIF},
		{"MagicNumberOfHEIC", ".test/blank.heic", HEIF},
		{"MagicNumberOfHEIF", ".test/blank.heif", HEIF},
		{"MagicNumberOfICO", ".test/blank.ico", ICO},
		{"MagicNumberOfJPEG", ".test/blank.jpeg", JPEG},
		{"MagicNumberOfPDF", ".test/blank.pdf", PDF},
		{"MagicNumberOfPNG", ".test/blank.png", PNG},
		{"MagicNumberOfTIFF", ".test/blank.tiff", TIFF},
		{"MagicNumberOfWEBP", ".test/blank.webp", WEBP},
		{"MagicNumberUnknown", ".test/blank.gif.tgz", Unknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := ioutil.ReadFile(tt.filePath)
			if err != nil {
				t.Error(err)
			}
			if got := FormatByMagicNumber(buf); got != tt.want {
				t.Errorf("FormatByMagicNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageFormat_MarshalText(t *testing.T) {
	tests := []struct {
		name    string
		imgFmt  ImageFormat
		want    []byte
		wantErr bool
	}{
		{"MarshalTextAVIF", AVIF, []byte("avif"), false},
		{"MarshalTextBMP", BMP, []byte("bmp"), false},
		{"MarshalTextGIF", GIF, []byte("gif"), false},
		{"MarshalTextHEIF", HEIF, []byte("heif"), false},
		{"MarshalTextICO", ICO, []byte("ico"), false},
		{"MarshalTextJPEG", JPEG, []byte("jpg"), false},
		{"MarshalTextPDF", PDF, []byte("pdf"), false},
		{"MarshalTextPNG", PNG, []byte("png"), false},
		{"MarshalTextSVG", SVG, []byte("svg"), false},
		{"MarshalTextTIFF", TIFF, []byte("tiff"), false},
		{"MarshalTextWEBP", WEBP, []byte("webp"), false},
		{"MarshalTextInvalidImageFormat1", Unknown, nil, true},
		{"MarshalTextInvalidImageFormat2", ImageFormat(42), nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.imgFmt.MarshalText()
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

func TestImageFormat_Ext(t *testing.T) {
	tests := []struct {
		name   string
		imgFmt ImageFormat
		want   string
	}{
		{"ExtensionAVIF", AVIF, ".avif"},
		{"ExtensionBMP", BMP, ".bmp"},
		{"ExtensionGIF", GIF, ".gif"},
		{"ExtensionHEIF", HEIF, ".heif"},
		{"ExtensionICO", ICO, ".ico"},
		{"ExtensionJPEG", JPEG, ".jpg"},
		{"ExtensionPDF", PDF, ".pdf"},
		{"ExtensionPNG", PNG, ".png"},
		{"ExtensionSVG", SVG, ".svg"},
		{"ExtensionTIFF", TIFF, ".tiff"},
		{"ExtensionWEBP", WEBP, ".webp"},
		{"ExtensionUnknown1", Unknown, ""},
		{"ExtensionUnknown2", ImageFormat(42), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.imgFmt.Ext(); got != tt.want {
				t.Errorf("Ext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageFormat_String(t *testing.T) {
	tests := []struct {
		name   string
		imgFmt ImageFormat
		want   string
	}{
		{"StringAVIF", AVIF, "avif"},
		{"StringBMP", BMP, "bmp"},
		{"StringGIF", GIF, "gif"},
		{"StringHEIF", HEIF, "heif"},
		{"StringICO", ICO, "ico"},
		{"StringJPEG", JPEG, "jpg"},
		{"StringPDF", PDF, "pdf"},
		{"StringPNG", PNG, "png"},
		{"StringSVG", SVG, "svg"},
		{"StringTIFF", TIFF, "tiff"},
		{"StringWEBP", WEBP, "webp"},
		{"StringUnknown1", Unknown, "Unknown"},
		{"StringUnknown2", ImageFormat(42), "Unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.imgFmt.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageFormat_UnmarshalText(t *testing.T) {
	type args struct {
		val []byte
	}
	tests := []struct {
		name    string
		args    args
		want    ImageFormat
		wantErr bool
	}{
		{"UnmarshalTextAVIFwDot", args{[]byte(".avif")}, AVIF, false},
		{"UnmarshalTextAVIFwoDot", args{[]byte("avif")}, AVIF, false},
		{"UnmarshalTextBMPwDot", args{[]byte(".bmp")}, BMP, false},
		{"UnmarshalTextBMPwoDot", args{[]byte("bmp")}, BMP, false},
		{"UnmarshalTextGIFwDot", args{[]byte(".gif")}, GIF, false},
		{"UnmarshalTextGIFwoDot", args{[]byte("gif")}, GIF, false},
		{"UnmarshalTextHEICwDot", args{[]byte(".heic")}, HEIF, false},
		{"UnmarshalTextHEICwoDot", args{[]byte("heic")}, HEIF, false},
		{"UnmarshalTextHEIFwDot", args{[]byte(".heif")}, HEIF, false},
		{"UnmarshalTextHEIFwoDot", args{[]byte("heif")}, HEIF, false},
		{"UnmarshalTextICOwDot", args{[]byte(".ico")}, ICO, false},
		{"UnmarshalTextICOwoDot", args{[]byte("ico")}, ICO, false},
		{"UnmarshalTextJPGwDot", args{[]byte(".jpg")}, JPEG, false},
		{"UnmarshalTextJPGwoDot", args{[]byte("jpg")}, JPEG, false},
		{"UnmarshalTextJPEGwDot", args{[]byte(".jpeg")}, JPEG, false},
		{"UnmarshalTextJPEGwoDot", args{[]byte("jpeg")}, JPEG, false},
		{"UnmarshalTextPDFwDot", args{[]byte(".pdf")}, PDF, false},
		{"UnmarshalTextPDFwoDot", args{[]byte("pdf")}, PDF, false},
		{"UnmarshalTextPNGwDot", args{[]byte(".png")}, PNG, false},
		{"UnmarshalTextPNGwoDot", args{[]byte("png")}, PNG, false},
		{"UnmarshalTextSVGwDot", args{[]byte(".svg")}, SVG, false},
		{"UnmarshalTextSVGwoDot", args{[]byte("svg")}, SVG, false},
		{"UnmarshalTextTIFFwDot", args{[]byte(".tiff")}, TIFF, false},
		{"UnmarshalTextTIFFwoDot", args{[]byte("tiff")}, TIFF, false},
		{"UnmarshalTextWEBPwDot", args{[]byte(".webp")}, WEBP, false},
		{"UnmarshalTextWEBPwoDot", args{[]byte("webp")}, WEBP, false},
		{"UnmarshalTextUnknown1", args{[]byte(".avi")}, Unknown, true},
		{"UnmarshalTextUnknown2", args{[]byte("avi")}, Unknown, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imgFmt := Unknown
			if err := imgFmt.UnmarshalText(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
