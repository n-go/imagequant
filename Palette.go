package imagequant

/*
#include "libimagequant.h"
*/
import "C"

// This struct has standard Go lifetime and does not need manual release.
type Palette struct {
	p C.struct_liq_palette
}

func (this *Palette) Count() uint {
	return uint(this.p.count)
}

func (this *Palette) At(idx int) (Color, error) {
	if idx < 0 || idx >= int(this.Count()) {
		return Color{}, ErrValueOutOfRange
	}

	return Color{
		r: uint8(this.p.entries[idx].r),
		g: uint8(this.p.entries[idx].g),
		b: uint8(this.p.entries[idx].b),
		a: uint8(this.p.entries[idx].a),
	}, nil

}
