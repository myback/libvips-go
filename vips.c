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

#include "vips.h"
#include <string.h>

int vips_initialize_go() {
    return vips_init("libvips-go");
}

int vips_image_load_go(void *buf, size_t len, int imgFmt, VipsImage **out) {
	if (imgFmt == JPEG) {
		return vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == PNG) {
		return vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == WEBP) {
		return vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == GIF) {
        return vips_gifload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
    } else if (imgFmt == PDF) {
    	return vips_pdfload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == BMP) {
	    return vips_magickload_buffer(buf, len, out, NULL);
	} else if (imgFmt == TIFF) {
    	return vips_tiffload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == HEIF || imgFmt == AVIF) {
	    return vips_heifload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
	} else if (imgFmt == SVG) {
        return vips_svgload_buffer(buf, len, out, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
    } else {
        vips_error("vips_image_load", "Unsupported image format");
        return 1;
    }

    return 0;
}

int vips_pdf_load_go(void *buf, size_t len, VipsImage **out, int page, int n) {
    if (page <= 0) {
        page = 0;
    }

    if (n == 0) {
        n = 1;
    }

    return vips_pdfload_buffer(buf, len, out, "page", page, "n", n, "access", VIPS_ACCESS_SEQUENTIAL, NULL);
}

VipsImage *vips_image_new_from_bytes_go(const void *data, size_t size, int width, int height) {
    return vips_image_new_from_memory(data, size, width, height, 4, VIPS_FORMAT_UCHAR);
}

VipsBandFormat vips_band_format_go(VipsImage *in) {
    return in->BandFmt;
}

void clear_image_go(VipsImage **in) {
    if (G_IS_OBJECT(*in)) g_clear_object(in);
}

void g_free_go(void **buf) {
    g_free(*buf);
}

void swap_and_clear_go(VipsImage **in, VipsImage *out) {
    clear_image_go(in);
    *in = out;
}

gboolean vips_is_animated_go(VipsImage * in) {
    return(vips_image_get_typeof(in, "page-height") != G_TYPE_INVALID &&
           vips_image_get_typeof(in, "gif-delay") != G_TYPE_INVALID &&
           vips_image_get_typeof(in, "gif-loop") != G_TYPE_INVALID);
}

int vips_get_orientation(VipsImage *image) {
#ifdef VIPS_META_ORIENTATION
    int orientation;

    if (vips_image_get_typeof(image, VIPS_META_ORIENTATION) == G_TYPE_INT &&
        vips_image_get_int(image, VIPS_META_ORIENTATION, &orientation) == 0)
        return orientation;
#else
    const char *orientation;

	if (vips_image_get_typeof(image, EXIF_ORIENTATION) == VIPS_TYPE_REF_STRING &&
		vips_image_get_string(image, EXIF_ORIENTATION, &orientation) == 0)
        return atoi(orientation);
#endif

	return 1;
}

int vips_addalpha_go(VipsImage *in, VipsImage **out) {
    return vips_addalpha(in, out, NULL);
}

int vips_copy_go(VipsImage *in, VipsImage **out) {
    return vips_copy(in, out, NULL);
}

int vips_cast_go(VipsImage *in, VipsImage **out, VipsBandFormat format) {
    return vips_cast(in, out, format, NULL);
}

int vips_rad2float_go(VipsImage *in, VipsImage **out) {
	return vips_rad2float(in, out, NULL);
}

int vips_resize_go(VipsImage *in, VipsImage **out, double scale) {
    return vips_resize(in, out, scale, NULL);
}

int vips_rotate_go(VipsImage *in, VipsImage **out, VipsAngle angle) {
    return vips_rot(in, out, angle, NULL);
}

int vips_flip_horizontal_go(VipsImage *in, VipsImage **out) {
    return vips_flip(in, out, VIPS_DIRECTION_HORIZONTAL, NULL);
}

int vips_ensure_alpha_go(VipsImage *in, VipsImage **out) {
    if (vips_image_hasalpha(in)) {
        return vips_copy(in, out, NULL);
    }

    return vips_bandjoin_const1(in, out, 255, NULL);
}

int vips_gaussblur_go(VipsImage *in, VipsImage **out, double sigma) {
    return vips_gaussblur(in, out, sigma, NULL);
}

int vips_sharpen_go(VipsImage *in, VipsImage **out, double sigma) {
    return vips_sharpen(in, out, "sigma", sigma, NULL);
}

int vips_trim_go(VipsImage *in, VipsImage **out, double threshold, gboolean smart, double r, double g, double b,
                 gboolean equal_hor, gboolean equal_ver) {

    VipsImage *tmp;

    if (vips_image_hasalpha(in)) {
    if (vips_flatten(in, &tmp, NULL))
        return 1;
    } else {
        if (vips_copy(in, &tmp, NULL))
            return 1;
    }

    double *bg;
    int bgn;
    VipsArrayDouble *bga;

    if (smart) {
        if (vips_getpoint(tmp, &bg, &bgn, 0, 0, NULL)) {
            clear_image_go(&tmp);
            return 1;
        }

        bga = vips_array_double_new(bg, bgn);
    } else {
        bga = vips_array_double_newv(3, r, g, b);
        bg = 0;
    }

    int left, right, top, bot, width, height, diff;

    if (vips_find_trim(tmp, &left, &top, &width, &height, "background", bga, "threshold", threshold, NULL)) {
        clear_image_go(&tmp);
        vips_area_unref((VipsArea *)bga);
        g_free(bg);
        return 1;
    }

    if (equal_hor) {
        right = in->Xsize - left - width;
        diff = right - left;

        if (diff > 0) {
            width += diff;
        } else if (diff < 0) {
            left = right;
            width -= diff;
        }
    }

    if (equal_ver) {
        bot = in->Ysize - top - height;
        diff = bot - top;

        if (diff > 0) {
            height += diff;
        } else if (diff < 0) {
            top = bot;
            height -= diff;
        }
    }

    clear_image_go(&tmp);
    vips_area_unref((VipsArea *)bga);
    g_free(bg);

    if (width == 0 || height == 0) {
        return vips_copy(in, out, NULL);
    }

    return vips_extract_area(in, out, left, top, width, height, NULL);
}

int vips_apply_watermark_go(VipsImage *in, VipsImage *watermark, VipsImage **out, int left, int top, float opacity) {
    VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 5);

	if (opacity < 1) {
    if (
        vips_extract_band(watermark, &t[0], 0, "n", watermark->Bands - 1, NULL) ||
        vips_extract_band(watermark, &t[1], watermark->Bands - 1, "n", 1, NULL) ||
        vips_linear1(t[1], &t[2], opacity, 0, NULL) ||
        vips_bandjoin2(t[0], t[2], &t[3], NULL)
    ) {
        clear_image_go(&base);
	        return 1;
	    }
	} else {
        if (vips_copy(watermark, &t[3], NULL)) {
            clear_image_go(&base);
            return 1;
        }
  }

    int res =
        vips_composite2(in, t[3], &t[4], VIPS_BLEND_MODE_OVER, "compositing_space", in->Type, "x", left, "y", top, NULL) ||
        vips_cast(t[4], out, vips_image_get_format(in), NULL);

    clear_image_go(&base);

    return res;
}

int vips_strip_go(VipsImage *in, VipsImage **out) {
    if (vips_copy(in, out, NULL)) return 1;

    gchar **fields = vips_image_get_fields(in);

    for (int i = 0; fields[i] != NULL; i++) {
        gchar *name = fields[i];

        if (strcmp(name, VIPS_META_ICC_NAME) == 0) continue;

        vips_image_remove(*out, name);
    }

    g_strfreev(fields);

    return 0;
}

int vips_smartcrop_go(VipsImage *in, VipsImage **out, int width, int height) {
    return vips_smartcrop(in, out, width, height, NULL);
}

int vips_extract_area_go(VipsImage *in, VipsImage **out, int left, int top, int width, int height) {
    return vips_extract_area(in, out, left, top, width, height, NULL);
}

int vips_jpegsave_go(VipsImage *in, void **buf, size_t *len, int quality, int strip, int interlace) {
    return vips_jpegsave_buffer(in, buf, len,
        "Q", quality,
        "strip", strip,
        "optimize_coding", TRUE,
        "interlace", interlace,
        NULL);
}

int vips_pngsave_go(VipsImage *in, void **buf, size_t *len, int compression, int strip, int interlace, int palette) {
    return vips_pngsave_buffer(in, buf, len,
        "compression", compression,
        "strip", strip,
        "filter", VIPS_FOREIGN_PNG_FILTER_NONE,
        "interlace", interlace,
        "palette", palette,
        NULL);
}

int vips_webpsave_go(VipsImage *in, void **buf, size_t *len, int quality, int strip, gboolean lossless) {
    return vips_webpsave_buffer(in, buf, len, "strip", strip, "Q", quality, "lossless", lossless, NULL);
}

int vips_gifsave_go(VipsImage *in, void **buf, size_t *len) {
    return vips_magicksave_buffer(in, buf, len, "format", "gif", NULL);
}

int vips_tiffsave_go(VipsImage *in, void **buf, size_t *len, int quality) {
    return vips_tiffsave_buffer(in, buf, len, "Q", quality, NULL);
}

int vips_avifsave_go(VipsImage *in, void **buf, size_t *len, int quality) {
    return vips_heifsave_buffer(in, buf, len, "Q", quality, "compression", VIPS_FOREIGN_HEIF_COMPRESSION_AV1, NULL);
}

int vips_heifsave_go(VipsImage *in, void **buf, size_t *len, int quality, int compression, gboolean lossless) {
    return vips_heifsave_buffer(in, buf, len, "Q", quality, "compression", compression, "lossless", lossless, NULL);
}

int vips_bmpsave_go(VipsImage *in, void **buf, size_t *len) {
    return vips_magicksave_buffer(in, buf, len, "format", "bmp", NULL);
}

int vips_resize_with_premultiply_go(VipsImage *in, VipsImage **out, double scale) {
	VipsBandFormat format;
    VipsImage *tmp1, *tmp2;

    format = vips_band_format_go(in);

    if (vips_premultiply(in, &tmp1, NULL))
        return 1;

    if (vips_resize(tmp1, &tmp2, scale, NULL)) {
        clear_image_go(&tmp1);
        return 1;
    }
    swap_and_clear_go(&tmp1, tmp2);

    if (vips_unpremultiply(tmp1, &tmp2, NULL)) {
        clear_image_go(&tmp1);
        return 1;
    }
    swap_and_clear_go(&tmp1, tmp2);

    if (vips_cast(tmp1, out, format, NULL)) {
        clear_image_go(&tmp1);
        return 1;
    }

    clear_image_go(&tmp1);

    return 0;
}
