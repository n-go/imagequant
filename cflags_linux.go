//+build !windows

package imagequant

/*
#cgo CFLAGS: -O3 -fopenmp -fomit-frame-pointer -Wall -Wno-attributes -std=c99 -DNDEBUG -DUSE_SSE=1 -msse
#cgo LDFLAGS: -lm -fopenmp -ldl
*/
import "C"
