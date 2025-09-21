//go:build js && wasm

package main

import (
	crand "crypto/rand"
	"syscall/js"
)

const (
	rows          = 64
	cols          = 92
	bytesPerPixel = 4
)

var (
	grid    [rows][cols]uint8
	palette = [256][4]uint8{
		0: {0, 0, 0, 255},     // black
		1: {255, 0, 0, 255},   // red
		2: {0, 255, 0, 255},   // green
		3: {0, 0, 255, 255},   // blue
		4: {50, 200, 80, 255}, // not sure :)
		5: {255, 255, 0, 255}, // yellow
	}
	nColors = 6
)

var (
	ctx js.Value
	img js.Value
	buf []byte
)

func putPixel(r, c int, color uint8) {
	if r < 0 || r >= rows || c < 0 || c >= cols {
		return // safety
	}
	if color >= uint8(nColors) {
		color = 0 // guard against out-of-range
	}
	grid[r][c] = color

	off := (r*cols + c) * bytesPerPixel
	rgba := palette[color]
	buf[off+0] = rgba[0]
	buf[off+1] = rgba[1]
	buf[off+2] = rgba[2]
	buf[off+3] = rgba[3]
}

func redraw() {
	dst := img.Get("data")
	js.CopyBytesToJS(dst, buf)
	ctx.Call("putImageData", img, 0, 0)
}

// JS: setCell(r, c, color)
func setCell(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return nil
	}
	r, c, color := args[0].Int(), args[1].Int(), args[2].Int()
	putPixel(r, c, uint8(color))
	redraw()
	return nil
}

// JS: randomize()
func randomize(this js.Value, args []js.Value) any {
	tmp := make([]byte, rows*cols)
	if _, err := crand.Read(tmp); err != nil {
		// deterministic fallback
		for i := range tmp {
			tmp[i] = byte(i)
		}
	}
	i := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := uint8(int(tmp[i]) % nColors) // 0..nColors-1
			putPixel(r, c, color)
			i++
		}
	}
	redraw()
	return nil
}

func main() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "canvas")

	canvas.Set("width", cols)
	canvas.Set("height", rows)

	ctx = canvas.Call("getContext", "2d") //
	img = js.Global().Get("ImageData").New(cols, rows)
	buf = make([]byte, rows*cols*bytesPerPixel)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			putPixel(r, c, 1)
		}
	}
	redraw()

	js.Global().Set("setCell", js.FuncOf(setCell))
	js.Global().Set("randomize", js.FuncOf(randomize))

	select {} // keep WASM running
}
