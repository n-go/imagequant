package imagequant

import (
	"errors"
	"unsafe"
)

/*
#cgo CFLAGS: -O3 -fno-math-errno -fopenmp -funroll-loops -fomit-frame-pointer -Wall -Wno-attributes -std=c99 -DNDEBUG -DUSE_SSE=1 -msse -fexcess-precision=fast
#cgo LDFLAGS: -fopenmp -static
#include "libimagequant.h"
*/
import "C"

var (
	ErrQualityTooLow      = errors.New("Quality too low")
	ErrValueOutOfRange    = errors.New("Value out of range")
	ErrOutOfMemory        = errors.New("Out of memory")
	ErrAborted            = errors.New("Aborted")
	ErrBitmapNotAvailable = errors.New("Bitmap not available")
	ErrBufferTooSmall     = errors.New("Buffer too small")
	ErrInvalidPointer     = errors.New("Invalid pointer")
)

func translateError(iqe C.liq_error) error {
	switch iqe {
	case C.LIQ_OK:
		return nil
	case (C.LIQ_QUALITY_TOO_LOW):
		return ErrQualityTooLow
	case (C.LIQ_VALUE_OUT_OF_RANGE):
		return ErrValueOutOfRange
	case (C.LIQ_OUT_OF_MEMORY):
		return ErrOutOfMemory
	case (C.LIQ_ABORTED):
		return ErrAborted
	case (C.LIQ_BITMAP_NOT_AVAILABLE):
		return ErrBitmapNotAvailable
	case (C.LIQ_BUFFER_TOO_SMALL):
		return ErrBufferTooSmall
	case (C.LIQ_INVALID_POINTER):
		return ErrInvalidPointer
	default:
		return errors.New("Unknown error")
	}
}

func GetLibraryVersion() int {
	return int(C.liq_version())
}

type Attributes struct {
	p *C.struct_liq_attr
}

// Callers MUST call Release() on the returned object to free memory.
func NewAttributes() (*Attributes, error) {
	pAttr := C.liq_attr_create()
	if pAttr == nil { // nullptr
		return nil, errors.New("Unsupported platform")
	}

	return &Attributes{p: pAttr}, nil
}

func (this *Attributes) SetMaxColors(colors int) error {
	return translateError(C.liq_set_max_colors(this.p, C.int(colors)))
}

func (this *Attributes) GetMaxColors() int {
	return int(C.liq_get_max_colors(this.p))
}

func (this *Attributes) SetQuality(minimum, maximum int) error {
	return translateError(C.liq_set_quality(this.p, C.int(minimum), C.int(maximum)))
}

func (this *Attributes) GetMinQuality() int {
	return int(C.liq_get_min_quality(this.p))
}

func (this *Attributes) GetMaxQuality() int {
	return int(C.liq_get_max_quality(this.p))
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Attributes) Release() {
	C.liq_attr_destroy(this.p)
}

type Image struct {
	p *C.struct_liq_image
}

// Callers MUST call Release() on the returned object to free memory.
func NewImage(attr *Attributes, rgba32data string, width, height int, gamma float64) (*Image, error) {
	pImg := C.liq_image_create_rgba(attr.p, unsafe.Pointer(C.CString(rgba32data)), C.int(width), C.int(height), C.double(gamma))
	if pImg == nil {
		return nil, errors.New("Failed to create image (invalid argument)")
	}

	return &Image{
		p: pImg,
	}, nil
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Image) Release() {
	C.liq_image_destroy(this.p)
}

type Result struct {
	p *C.struct_liq_result
}

func (this *Image) Quantize(attr *Attributes) (*Result, error) {
	res := Result{}
	liqerr := C.liq_image_quantize(this.p, attr.p, &res.p)
	if liqerr != C.LIQ_OK {
		return nil, translateError(liqerr)
	}

	return &res, nil
}

func (this *Result) SetDitheringLevel(dither_level float32) error {
	return translateError(C.liq_set_dithering_level(this.p, C.float(dither_level)))
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Result) Release() {
	C.liq_result_destroy(this.p)
}
