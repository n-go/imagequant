package imagequant

import (
	"errors"
	"unsafe"
)

/*
#include "libimagequant.h"
*/
import "C"

type Image struct {
	p        *C.struct_liq_image
	w, h     int
	released bool
}

// Callers MUST call Release() on the returned object to free memory.
func NewImage(attr *Attributes, rgba32data string, width, height int, gamma float64) (*Image, error) {
	pImg := C.liq_image_create_rgba(attr.p, unsafe.Pointer(C.CString(rgba32data)), C.int(width), C.int(height), C.double(gamma))
	if pImg == nil {
		return nil, errors.New("Failed to create image (invalid argument)")
	}

	return &Image{
		p:        pImg,
		w:        width,
		h:        height,
		released: false,
	}, nil
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Image) Release() {
	C.liq_image_destroy(this.p)
	this.released = true
}

func (this *Image) Quantize(attr *Attributes) (*Result, error) {
	res := Result{}
	liqerr := C.liq_image_quantize(this.p, attr.p, &res.p)
	if liqerr != C.LIQ_OK {
		return nil, translateError(liqerr)
	}

	return &res, nil
}
