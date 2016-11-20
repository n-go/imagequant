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

type ImageQuantErr int

var (
	ErrOK                 = ImageQuantErr(C.LIQ_OK)
	ErrQualityTooLow      = ImageQuantErr(C.LIQ_QUALITY_TOO_LOW)
	ErrValueOutOfRange    = ImageQuantErr(C.LIQ_VALUE_OUT_OF_RANGE)
	ErrOutOfMemory        = ImageQuantErr(C.LIQ_OUT_OF_MEMORY)
	ErrAborted            = ImageQuantErr(C.LIQ_ABORTED)
	ErrBitmapNotAvailable = ImageQuantErr(C.LIQ_BITMAP_NOT_AVAILABLE)
	ErrBufferTooSmall     = ImageQuantErr(C.LIQ_BUFFER_TOO_SMALL)
	ErrInvalidPointer     = ImageQuantErr(C.LIQ_INVALID_POINTER)
)

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

func (this *Attributes) SetMaxColors(colors int) ImageQuantErr {
	ret := C.liq_set_max_colors(this.p, C.int(colors))
	return ImageQuantErr(ret)
}

func (this *Attributes) GetMaxColors() int {
	return int(C.liq_get_max_colors(this.p))
}

func (this *Attributes) SetQuality(minimum, maximum int) ImageQuantErr {
	ret := C.liq_set_quality(this.p, C.int(minimum), C.int(maximum))
	return ImageQuantErr(ret)
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
func NewImage(attr *Attributes, rgba32data string, width, height int, gamma float64) *Image {
	pImg := C.liq_image_create_rgba(attr.p, unsafe.Pointer(C.CString(rgba32data)), C.int(width), C.int(height), C.double(gamma))
	return &Image{
		p: pImg,
	}
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Image) Release() {
	C.liq_image_destroy(this.p)
}
