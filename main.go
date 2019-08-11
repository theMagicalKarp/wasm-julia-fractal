package main

import (
	"fmt"
	"time"
	"syscall/js"
	"math"
)

func min(a,b byte) byte {
	if a < b {
		return a
	}
	return b
}

func fillBuffer(x0, y0 float32, width, height int, buffer []byte, max_iterations, r,g,b byte) {
	widthf := float32(width)
	heightf := float32(height)
	var x, x_n, y, y_n, i, j float32

	var iteration byte

	for i = 0.0; i < heightf; i++ {
		for j = 0.0; j < widthf; j++ {
			x = -1.5 + j * 3.0 / widthf
			y = -1.0 + i * 2.0 / heightf
			iteration = 0

			for (x * x + y * y < 4.0) && (iteration < max_iterations) {
				x_n = x * x - y * y + x0
				y_n = 2.0 * x * y + y0
				x = x_n
				y = y_n
				iteration++
			}

			z := int(i) * width * 4 + int(j) * 4
			buffer[z + 0] = min(iteration*r, 254)
			buffer[z + 1] = min(iteration*g, 254)
			buffer[z + 2] = min(iteration*b, 254)
			buffer[z + 3] = 254
		}
	}
}

func main() {
	fmt.Println("Starting")
	doc := js.Global().Get("document")
	clientWidth := doc.Get("body").Get("clientWidth").Int()
	clientHeight := doc.Get("body").Get("clientHeight").Int()

	canvasEl := doc.Call("getElementById", "canvas")
	canvasEl.Set("width", clientWidth)
	canvasEl.Set("height", clientHeight)
	ctx := canvasEl.Call("getContext", "2d")
	ctx.Set("imageSmoothingEnabled", false)

	scale := 8
	ctx.Call("scale", scale, scale)

	width := (clientWidth/scale) + 1
	height := (clientHeight/scale) + 1

	scaledCanvas := doc.Call("createElement", "canvas");
	scaledCanvas.Set("width", width)
	scaledCanvas.Set("height", height)
	scaledCtx := scaledCanvas.Call("getContext", "2d")
	scaledCtx.Set("imageSmoothingEnabled", false)

	uint8Array := js.Global().Get("Uint8Array")
	jsBuffer := uint8Array.New(width*height*4)
	goBuffer := make([]byte, width*height*4)

	var max_iterations, r, g, b byte
	max_iterations = 50
	r = 20
	g = 5
	b = 8
	x0 := float32(-0.4)
	y0 := float32(-0.6)

	fillBuffer(x0, y0, width, height, goBuffer, max_iterations, r, g, b)

	imageData := scaledCtx.Call("createImageData", width, height)
	imageDataData := imageData.Get("data")

	mousex := 0.0
	mousey := 0.0
	thetaX := 0.0
	thetaY := 0.0

	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		e.Call("preventDefault")
		touches := e.Get("touches")

		if touches.Truthy() && touches.Length() >= 1 {
			touch := touches.Index(0)
			mousex = touch.Get("clientX").Float() / float64(clientWidth) * 2.0 - 1.0
			mousey = touch.Get("clientY").Float() / float64(clientHeight) * 2.0 - 1.0
		} else {
			mousex = e.Get("clientX").Float() / float64(clientWidth) * 2.0 - 1.0
			mousey = e.Get("clientY").Float() / float64(clientHeight) * 2.0 - 1.0
		}
		return nil
	})

	defer mouseMoveEvt.Release()

	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
	doc.Call("addEventListener", "touchmove", mouseMoveEvt)

	for {
		thetaX = thetaX - mousex
		thetaY = thetaY - mousey

		x0 =float32(math.Sin(thetaX*0.0174533)) * 0.50 - 0.4
		y0 = float32(math.Sin(thetaY*0.0174533)) * 0.50 - 0.6

		now := time.Now()
		js.CopyBytesToJS(jsBuffer, goBuffer)
		imageDataData.Call("set", jsBuffer)

		scaledCtx.Call("putImageData", imageData, 0, 0)
		ctx.Call("drawImage", scaledCanvas, 0, 0)

		fillBuffer(x0, y0, width, height, goBuffer, max_iterations, r, g, b)

		wait := 16 * time.Millisecond - time.Since(now)
		if wait <= 0 {
			wait = 1 * time.Millisecond
		}

		time.Sleep(wait)
	}
}
