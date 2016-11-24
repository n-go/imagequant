package imagequant

import (
	"errors"
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
