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

#include <stdlib.h>

#include <vips/vips.h>
#include <vips/vips7compat.h>
#include <vips/vector.h>

enum ImageFormat {
    UNKNOWN = 0,
    AVIF,
    BMP,
  	GIF,
  	HEIF,
  	ICO,
  	JPEG,
  	PDF,
  	PNG,
  	SVG,
  	TIFF,
  	WEBP
};

int vips_initialize_go();
int vips_image_load_go(void *buf, size_t len, int imgtype, VipsImage **out);
int vips_pdf_load_go(void *buf, size_t len, VipsImage **out, int page, int n);
VipsImage *vips_image_new_from_bytes_go(const void *data, size_t size, int width, int height);

VipsBandFormat vips_band_format_go(VipsImage *in);

void clear_image_go(VipsImage **in);
void g_free_go(void **buf);
void swap_and_clear_go(VipsImage **in, VipsImage *out);

gboolean vips_is_animated_go(VipsImage * in);

int vips_get_orientation(VipsImage *image);
int vips_addalpha_go(VipsImage *in, VipsImage **out);
int vips_copy_go(VipsImage *in, VipsImage **out);
int vips_cast_go(VipsImage *in, VipsImage **out, VipsBandFormat format);
int vips_rad2float_go(VipsImage *in, VipsImage **out);
int vips_resize_go(VipsImage *in, VipsImage **out, double scale);
int vips_rotate_go(VipsImage *in, VipsImage **out, VipsAngle angle);
int vips_flip_horizontal_go(VipsImage *in, VipsImage **out);
int vips_ensure_alpha_go(VipsImage *in, VipsImage **out);
int vips_gaussblur_go(VipsImage *in, VipsImage **out, double sigma);
int vips_sharpen_go(VipsImage *in, VipsImage **out, double sigma);
int vips_trim_go(VipsImage *in, VipsImage **out, double threshold, gboolean smart, double r, double g, double b,
                 gboolean equal_hor, gboolean equal_ver);
int vips_apply_watermark_go(VipsImage *in, VipsImage *watermark, VipsImage **out, int left, int top, float opacity);
int vips_strip_go(VipsImage *in, VipsImage **out);
int vips_smartcrop_go(VipsImage *in, VipsImage **out, int width, int height);
int vips_extract_area_go(VipsImage *in, VipsImage **out, int left, int top, int width, int height);

int vips_jpegsave_go(VipsImage *in, void **buf, size_t *len, int quality, int strip, int interlace);
int vips_pngsave_go(VipsImage *in, void **buf, size_t *len, int compression, int strip, int interlace, int palette);
int vips_webpsave_go(VipsImage *in, void **buf, size_t *len, int quality, int strip, int lossless);
int vips_gifsave_go(VipsImage *in, void **buf, size_t *len);
int vips_tiffsave_go(VipsImage *in, void **buf, size_t *len, int quality);
int vips_avifsave_go(VipsImage *in, void **buf, size_t *len, int quality);
int vips_heifsave_go(VipsImage *in, void **buf, size_t *len, int quality, int compression, int lossless);
int vips_bmpsave_go(VipsImage *in, void **buf, size_t *len);
int vips_pdfsave_go(VipsImage *in, void **buf, size_t *len);

int vips_resize_with_premultiply_go(VipsImage *in, VipsImage **out, double scale);

int vips_arrayjoin_go(VipsImage **in, VipsImage **out, int n);
