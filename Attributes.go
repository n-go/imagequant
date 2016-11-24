package imagequant

import (
	"errors"
)

/*
#include "libimagequant.h"
*/
import "C"

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

func (this *Attributes) CreateHistogram() *Histogram {
	ptr := C.liq_histogram_create(this.p)
	return &Histogram{p: ptr}
}

// Free memory. Callers must not use this object after Release has been called.
func (this *Attributes) Release() {
	C.liq_attr_destroy(this.p)
}
