package imagequant

/*
#cgo CFLAGS: -O3 -fno-math-errno -fopenmp -funroll-loops -fomit-frame-pointer -Wall -Wno-attributes -std=c99 -DNDEBUG -DUSE_SSE=1 -msse -fexcess-precision=fast
#cgo LDFLAGS: -fopenmp -static
#include "libimagequant.h"
*/
import (
	"C"
)

func GetLibraryVersion() int {
	return int(C.liq_version())
}
