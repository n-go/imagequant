package imagequant

import (
	"errors"
	"unsafe"
)

/*
#cgo CFLAGS: -O3 -fno-math-errno -fopenmp -funroll-loops -fomit-frame-pointer -Wall -Wno-attributes -std=c99 -DNDEBUG -DUSE_SSE=1 -msse -fexcess-precision=fast
#cgo LDFLAGS: -fopenmp -static
#include "libimagequant.h"

const char* liqVersionString() {
	return LIQ_VERSION_STRING;
}

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

	ErrUseAfterFree = errors.New("Use after free")
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

func GetLibraryVersionString() string {
	return C.GoString(C.liqVersionString())
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

const (
	SPEED_SLOWEST = 1
	SPEED_DEFAULT = 3
	SPEED_FASTEST = 10
)

func (this *Attributes) SetSpeed(speed int) error {
	return translateError(C.liq_set_speed(this.p, C.int(speed)))
}

func (this *Attributes) GetSpeed() int {
	return int(C.liq_get_speed(this.p))
}

func (this *Attributes) SetMinOpacity(min int) error {
	return translateError(C.liq_set_min_opacity(this.p, C.int(min)))
}

func (this *Attributes) GetMinOpacity() int {
	return int(C.liq_get_min_opacity(this.p))
}

func (this *Attributes) SetMinPosterization(bits int) error {
	return translateError(C.liq_set_min_posterization(this.p, C.int(bits)))
}

func (this *Attributes) GetMinPosterization() int {
	return int(C.liq_get_min_posterization(this.p))
}

func (this *Attributes) SetLastIndexTransparent(is_last int) {
	C.liq_set_last_index_transparent(this.p, C.int(is_last))
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Attributes) Release() {
	C.liq_attr_destroy(this.p)
}

type Histogram struct {
	p *C.struct_liq_histogram
}

func (this *Attributes) CreateHistogram() *Histogram {
	ptr := C.liq_histogram_create(this.p)
	return &Histogram{p: ptr}
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Histogram) Release() {
	C.liq_histogram_destroy(this.p)
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

// Callers must not use this object once Release has been called on the parent
// Image struct.
type Result struct {
	p  *C.struct_liq_result
	im *Image
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

// Free memory. Callers must not use this object after Release has been called.
func (this *Result) Release() {
	C.liq_result_destroy(this.p)
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

// This struct has standard Go lifetime and does not need manual release.
type Palette struct {
	p C.struct_liq_palette
}

func (this *Result) GetPalette() *Palette {
	ptr := *C.liq_get_palette(this.p) // copy struct content
	return &Palette{p: ptr}
}
