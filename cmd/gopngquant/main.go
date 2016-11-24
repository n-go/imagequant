package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"

	"code.ivysaur.me/imagequant"
)

func GoImageToRgba32(im image.Image) []byte {
	w := im.Bounds().Max.X
	h := im.Bounds().Max.Y
	ret := make([]byte, w*h*4)

	for y := 0; y < h; y += 1 {
		for x := 0; x < w; x += 1 {
			r16, g16, b16, a16 := im.At(x, y).RGBA() // Each value ranges within [0, 0xffff]

			ret[y*h+x+0] = uint8(r16 >> 8)
			ret[y*h+x+1] = uint8(g16 >> 8)
			ret[y*h+x+2] = uint8(b16 >> 8)
			ret[y*h+x+3] = uint8(a16 >> 8)
		}
	}

	return ret
}

func Crush(sourcefile, destfile string) error {

	fh, err := os.OpenFile(sourcefile, os.O_RDONLY, 0444)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer fh.Close()

	img, err := png.Decode(fh)
	if err != nil {
		return fmt.Errorf("png.Decode: %s", err.Error())
	}

	attr, err := imagequant.NewAttributes()
	if err != nil {
		return fmt.Errorf("NewAttributes: %s", err.Error())
	}
	defer attr.Release()

	rgba32data := GoImageToRgba32(img)

	iqm, err := imagequant.NewImage(attr, string(rgba32data), img.Bounds().Max.X, img.Bounds().Max.Y, 0)
	if err != nil {
		return fmt.Errorf("NewImage: %s", err.Error())
	}
	defer iqm.Release()

	res, err := iqm.Quantize(attr)
	if err != nil {
		return fmt.Errorf("Quantize: %s", err.Error())
	}
	defer res.Release()

	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		return fmt.Errorf("WriteRemappedImage: %s", err.Error())
	}

	rect := image.Rectangle{Max: image.Point{X: res.GetImageWidth(), Y: res.GetImageHeight()}}
	pal := res.GetPalette()
	im2 := image.NewPaletted(rect, pal)
	for y := 0; y < rect.Max.Y; y += 1 {
		for x := 0; x < rect.Max.X; x += 1 {
			im2.Set(x, y, pal[rgb8data[y*rect.Max.Y+x]])
		}
	}

	fh2, err := os.OpenFile(destfile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("os.OpenFile: %s", err.Error())
	}
	defer fh2.Close()

	penc := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = penc.Encode(fh2, im2)
	if err != nil {
		return fmt.Errorf("png.Encode: %s", err.Error())
	}

	return nil
}

func main() {
	ShouldDisplayVersion := flag.Bool("Version", false, "")
	Infile := flag.String("In", "", "Input filename")
	Outfile := flag.String("Out", "", "Output filename")

	flag.Parse()

	if *ShouldDisplayVersion {
		fmt.Printf("libimagequant '%s' (%d)\n", imagequant.GetLibraryVersionString(), imagequant.GetLibraryVersion())
		os.Exit(1)
	}

	err := Crush(*Infile, *Outfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
