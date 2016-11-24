package imagequant

/*
#include "libimagequant.h"
*/
import "C"

type Histogram struct {
	p *C.struct_liq_histogram
}

func (this *Histogram) AddImage(attr *Attributes, img *Image) error {
	return translateError(C.liq_histogram_add_image(this.p, attr.p, img.p))
}

func (this *Histogram) Quantize(attr *Attributes) (*Result, error) {
	res := Result{}
	liqerr := C.liq_histogram_quantize(this.p, attr.p, &res.p)
	if liqerr != C.LIQ_OK {
		return nil, translateError(liqerr)
	}

	return &res, nil
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Histogram) Release() {
	C.liq_histogram_destroy(this.p)
}
