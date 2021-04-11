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
		{"UnmarshalTextAVIFwPoint", args{[]byte(".avif")}, AVIF, false},
		{"UnmarshalTextAVIFwoPoint", args{[]byte("avif")}, AVIF, false},
		{"UnmarshalTextBMPwPoint", args{[]byte(".bmp")}, BMP, false},
		{"UnmarshalTextBMPwoPoint", args{[]byte("bmp")}, BMP, false},
		{"UnmarshalTextGIFwPoint", args{[]byte(".gif")}, GIF, false},
		{"UnmarshalTextGIFwoPoint", args{[]byte("gif")}, GIF, false},
		{"UnmarshalTextHEICwPoint", args{[]byte(".heic")}, HEIF, false},
		{"UnmarshalTextHEICwoPoint", args{[]byte("heic")}, HEIF, false},
		{"UnmarshalTextHEIFwPoint", args{[]byte(".heif")}, HEIF, false},
		{"UnmarshalTextHEIFwoPoint", args{[]byte("heif")}, HEIF, false},
		{"UnmarshalTextICOwPoint", args{[]byte(".ico")}, ICO, false},
		{"UnmarshalTextICOwoPoint", args{[]byte("ico")}, ICO, false},
		{"UnmarshalTextJPGwPoint", args{[]byte(".jpg")}, JPEG, false},
		{"UnmarshalTextJPGwoPoint", args{[]byte("jpg")}, JPEG, false},
		{"UnmarshalTextJPEGwPoint", args{[]byte(".jpeg")}, JPEG, false},
		{"UnmarshalTextJPEGwoPoint", args{[]byte("jpeg")}, JPEG, false},
		{"UnmarshalTextPDFwPoint", args{[]byte(".pdf")}, PDF, false},
		{"UnmarshalTextPDFwoPoint", args{[]byte("pdf")}, PDF, false},
		{"UnmarshalTextPNGwPoint", args{[]byte(".png")}, PNG, false},
		{"UnmarshalTextPNGwoPoint", args{[]byte("png")}, PNG, false},
		{"UnmarshalTextSVGwPoint", args{[]byte(".svg")}, SVG, false},
		{"UnmarshalTextSVGwoPoint", args{[]byte("svg")}, SVG, false},
		{"UnmarshalTextTIFFwPoint", args{[]byte(".tiff")}, TIFF, false},
		{"UnmarshalTextTIFFwoPoint", args{[]byte("tiff")}, TIFF, false},
		{"UnmarshalTextWEBPwPoint", args{[]byte(".webp")}, WEBP, false},
		{"UnmarshalTextWEBPwoPoint", args{[]byte("webp")}, WEBP, false},
		{"UnmarshalTextUnknown1", args{[]byte(".avi")}, Unknown, false},
		{"UnmarshalTextUnknown2", args{[]byte("bmpp")}, Unknown, false},

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
