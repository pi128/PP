//go:build js && wasm

package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
	"strconv"
	"syscall/js"
)

/*
Behavior:
- Ocean palette
- Big emojis / fewer cells (chunky look)
- Hover-stamp overwrite: entering a cell always paints a new random emoji
- High-rate input: pointerrawupdate + pointermove + mousemove
- Canvas stretches to fill the viewport (no gutters)
- Seam-free background repaint per tile (device-pixel snapped)
- No-draw UI zone in top-left for overlay buttons
*/

const (
	rows = 24
	cols = 36
	tile = 64

	// No-draw zone under the in-canvas UI (set to 0 to disable)
	uiW = 9 * tile
	uiH = 3 * tile
)

var (
	canvas js.Value
	ctx    js.Value

	// One emoji per cell ("" = water)
	cells [rows][cols]string

	// Fast RNG seeded once from crypto/rand
	rng *rand.Rand

	// 🌊 Ocean palette (deduped)
	emojis = []string{
		"🌊", "💧", "🫧", "🌧️",
		"🐟", "🐠", "🐡", "🐬", "🐳", "🐋", "🦈",
		"🐙", "🦑", "🪼", "🦀", "🦞", "🦐", "🦪",
		"🐢", "🐚", "🪸",
		"⚓️", "🚢", "⛵️", "🛥️", "🚤", "🛟",
		"🏝️", "🏖️", "🦭", "🏄‍♂️",
	}

	// Logical drawing size (before scaling)
	logicalW = cols * tile
	logicalH = rows * tile

	// Current transform (device pixels)
	dpr    = 1.0
	scaleX = 1.0
	scaleY = 1.0
	tx     = 0.0
	ty     = 0.0

	// Last stamped cell (avoid re-rolling while stationary)
	prevR = -1
	prevC = -1

	// Keep handlers alive (avoid GC)
	resizeFn, rawMoveFn, pointerMoveFn, mouseMoveFn, pointerLeaveFn, mouseLeaveFn js.Func
)

/* ---------- utils ---------- */

func seedRNG() {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		binary.LittleEndian.PutUint64(b[:], 0xfeedf00ddeadbeef)
	}
	rng = rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}
func randEmoji() string { return emojis[rng.Intn(len(emojis))] }
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

/* ---------- drawing ---------- */

func setCanvasState() {
	size := tile * 92 / 100 // ~92% to keep glyphs off the very edge
	ctx.Set("textAlign", "center")
	ctx.Set("textBaseline", "middle")
	ctx.Set("font", strconv.Itoa(size)+"px Apple Color Emoji, Segoe UI Emoji, Noto Color Emoji, sans-serif")
}

func fillOcean() {
	ctx.Set("fillStyle", "#0a99c9")
	ctx.Call("fillRect", 0, 0, logicalW, logicalH)
}

// repaint this tile’s water in device space, snapped to integer pixels, with a small bleed
func paintTileBG(r, c int) {
	dx := tx + float64(c*tile)*scaleX
	dy := ty + float64(r*tile)*scaleY
	dw := float64(tile) * scaleX
	dh := float64(tile) * scaleY

	dx = math.Floor(dx - 1)
	dy = math.Floor(dy - 1)
	dw = math.Ceil(dw + 2)
	dh = math.Ceil(dh + 2)

	ctx.Call("save")
	ctx.Call("setTransform", 1, 0, 0, 1, 0, 0)
	ctx.Call("clearRect", dx, dy, dw, dh)
	ctx.Set("fillStyle", "#0a99c9")
	ctx.Call("fillRect", dx, dy, dw, dh)
	ctx.Call("restore")
}

func drawCell(r, c int) {
	if r < 0 || r >= rows || c < 0 || c >= cols {
		return
	}
	paintTileBG(r, c)
	if s := cells[r][c]; s != "" {
		x := c*tile + tile/2
		y := r*tile + tile/2
		ctx.Call("fillText", s, x, y)
	}
}

func redrawAll() {
	fillOcean()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			drawCell(r, c)
		}
	}
}

func paintDot(r, c int) {
	if r < 0 || r >= rows || c < 0 || c >= cols {
		return
	}
	cells[r][c] = randEmoji() // overwrite every time we enter
	drawCell(r, c)
}

/* ---------- layout: fill viewport (stretch X & Y) ---------- */

func layout() {
	win := js.Global()

	// DPR
	if v := win.Get("devicePixelRatio"); v.Truthy() {
		dpr = v.Float()
	} else {
		dpr = 1.0
	}

	// Use entire viewport (UI floats over canvas)
	vw := int(win.Get("innerWidth").Float())
	vh := int(win.Get("innerHeight").Float())

	// CSS size
	canvas.Get("style").Set("width", strconv.Itoa(vw)+"px")
	canvas.Get("style").Set("height", strconv.Itoa(vh)+"px")

	// Backing store in device px
	canvas.Set("width", int(float64(vw)*dpr))
	canvas.Set("height", int(float64(vh)*dpr))

	if ctx.IsUndefined() || ctx.IsNull() {
		ctx = canvas.Call("getContext", "2d")
	}

	// Stretch to fill both directions
	scaleX = float64(canvas.Get("width").Int()) / float64(logicalW)
	scaleY = float64(canvas.Get("height").Int()) / float64(logicalH)
	tx, ty = 0, 0

	ctx.Call("setTransform", scaleX, 0, 0, scaleY, tx, ty)
	setCanvasState()
}

/* ---------- mapping ---------- */

func mapEventToPos(ev js.Value) (lx, ly float64, r, c int) {
	b := canvas.Call("getBoundingClientRect")
	devX := (ev.Get("clientX").Float() - b.Get("left").Float()) * dpr
	devY := (ev.Get("clientY").Float() - b.Get("top").Float()) * dpr

	lx = (devX - tx) / scaleX
	ly = (devY - ty) / scaleY

	if lx < 0 {
		lx = 0
	}
	if ly < 0 {
		ly = 0
	}
	if lx > float64(logicalW) {
		lx = float64(logicalW)
	}
	if ly > float64(logicalH) {
		ly = float64(logicalH)
	}

	c = clamp(int(lx)/tile, 0, cols-1)
	r = clamp(int(ly)/tile, 0, rows-1)
	return
}

/* ---------- hover-stamp (overwrite) ---------- */

func handleMove(ev js.Value) {
	lx, ly, r, c := mapEventToPos(ev)

	// No-draw zone for in-canvas UI (top-left)
	if lx < float64(uiW) && ly < float64(uiH) {
		prevR, prevC = -1, -1
		return
	}

	// Same cell → do nothing
	if r == prevR && c == prevC {
		return
	}

	paintDot(r, c)
	prevR, prevC = r, c
}

func onMove(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		handleMove(args[0])
	}
	return nil
}
func onLeave(this js.Value, _ []js.Value) any { prevR, prevC = -1, -1; return nil }

/* ---------- JS-callable ---------- */

func clearCanvas(this js.Value, _ []js.Value) any {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cells[r][c] = ""
		}
	}
	redrawAll()
	return nil
}

/* ---------- main ---------- */

func main() {
	seedRNG()

	doc := js.Global().Get("document")
	canvas = doc.Call("getElementById", "canvas")
	canvas.Get("style").Set("touchAction", "none")

	layout()
	redrawAll() // start as empty ocean

	// export Clear (Save is handled in HTML via toDataURL)
	js.Global().Set("clearCanvas", js.FuncOf(clearCanvas))

	// keep handlers alive & bind high-rate events
	resizeFn = js.FuncOf(func(js.Value, []js.Value) any { layout(); redrawAll(); return nil })
	rawMoveFn = js.FuncOf(onMove)
	pointerMoveFn = js.FuncOf(onMove)
	mouseMoveFn = js.FuncOf(onMove)
	pointerLeaveFn = js.FuncOf(onLeave)
	mouseLeaveFn = js.FuncOf(onLeave)

	js.Global().Call("addEventListener", "resize", resizeFn)
	canvas.Call("addEventListener", "pointerrawupdate", rawMoveFn)
	canvas.Call("addEventListener", "pointermove", pointerMoveFn)
	canvas.Call("addEventListener", "mousemove", mouseMoveFn)
	canvas.Call("addEventListener", "pointerleave", pointerLeaveFn)
	canvas.Call("addEventListener", "mouseleave", mouseLeaveFn)

	select {}
}
