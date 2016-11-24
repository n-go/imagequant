package imagequant

/*
#include "libimagequant.h"
*/
import "C"

// Callers must not use this object once Release has been called on the parent
// Image struct.
type Result struct {
	p  *C.struct_liq_result
	im *Image
}

func (this *Result) SetDitheringLevel(dither_level float32) error {
	return translateError(C.liq_set_dithering_level(this.p, C.float(dither_level)))
}

func (this *Result) GetQuantizationError() float64 {
	return float64(C.liq_get_quantization_error(this.p))
}

func (this *Result) GetRemappingError() float64 {
	return float64(C.liq_get_remapping_error(this.p))
}

func (this *Result) GetQuantizationQuality() float64 {
	return float64(C.liq_get_quantization_quality(this.p))
}

func (this *Result) GetRemappingQuality() float64 {
	return float64(C.liq_get_remapping_quality(this.p))
}

func (this *Result) SetOutputGamma(gamma float64) error {
	return translateError(C.liq_set_output_gamma(this.p, C.double(gamma)))
}

func (this *Result) GetImageWidth() int {
	// C.liq_image_get_width
	return this.im.w
}

func (this *Result) GetImageHeight() int {
	// C.liq_image_get_height
	return this.im.h
}

func (this *Result) GetOutputGamma() float64 {
	return float64(C.liq_get_output_gamma(this.p))
}

func (this *Result) WriteRemappedImage() ([]byte, error) {
	if this.im.released {
		return nil, ErrUseAfterFree
	}

	buff_size := this.im.w * this.im.h
	buff := make([]byte, buff_size)

	// n.b. C.CBytes() added in go1.7.3

	iqe := C.liq_write_remapped_image(this.p, this.im.p, C.CBytes(buff), C.size_t(buff_size))
	if iqe != C.LIQ_OK {
		return nil, translateError(iqe)
	}

	return buff, nil
}

func (this *Result) GetPalette() *Palette {
	ptr := *C.liq_get_palette(this.p) // copy struct content
	return &Palette{p: ptr}
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Result) Release() {
	C.liq_result_destroy(this.p)
}
