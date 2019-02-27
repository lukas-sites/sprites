// +build js,wasm

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	"image/png"
	"syscall/js"
)

func split(img image.Image, n int) *gif.GIF {
	var generatedGif gif.GIF
	bounds := img.Bounds()
	width := bounds.Max.X / n
	height := bounds.Max.Y

	for i := 0; i < n; i++ {
		out := image.NewPaletted(image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{X: width, Y: height},
		}, palette.WebSafe)

		offset := width * i
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// offset x
				out.Set(x, y, img.At(x+offset, y))
			}
		}

		generatedGif.Image = append(generatedGif.Image, out)
		generatedGif.Delay = append(generatedGif.Delay, 1)
	}

	return &generatedGif
}

func createGif(this js.Value, args []js.Value) interface{} {
	raw := args[0].String()
	numFrames := args[1].Int()
	inBuff := bytes.NewBufferString(raw)

	r := base64.NewDecoder(base64.StdEncoding, inBuff)

	img, err := png.Decode(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	generatedGif := split(img, numFrames)

	var outBuff bytes.Buffer
	w := base64.NewEncoder(base64.StdEncoding, &outBuff)
	defer w.Close()

	if err = gif.EncodeAll(w, generatedGif); err != nil {
		fmt.Println(err)
		return ""
	}

	return outBuff.String()
}

func main() {
	done := make(chan struct{})
	js.Global().Set("createGif", js.FuncOf(createGif))
	<-done
}
