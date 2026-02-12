// -----------------------------------------------------------------------
//
// This file is part of the PixelWindow Project
//
//	by Amos Tibaldi - tibaldi at users.sourceforge.net
//
// https://sourceforge.net/projects/pixelwindow/
//
// https://github.com/Amos-Tibaldi/PixelWindow
//
// COPYRIGHT: http://www.gnu.org/licenses/gpl.html
//
//	       COPYRIGHT-gpl-3.0.txt
//
//	The PixelWindow Project
//	   PixelWindow gives high performance pixel access to DirectX windows
//	in go and in c++.
//
//	Copyright (C) 2022 Amos Tibaldi
//
//	This program is free software: you can redistribute it and/or modify
//	it under the terms of the GNU General Public License as published by
//	the Free Software Foundation, either version 3 of the License, or
//	(at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even the implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// -----------------------------------------------------------------------
package main

import (
	"PixelWindowGo/PixelWindowGo"
	"image"
	"image/color"
	"log"
	"os"

	"unsafe"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var bIOwnTheMessagePump bool = false
var luckyhwnd PixelWindowGo.HWND = 0

const imgsizebytes = 640 * 480 * 4
const textwidthx = 640

var redbuffer [imgsizebytes]byte
var greenbuffer [imgsizebytes]byte
var bluebuffer [imgsizebytes]byte

var thewindows [3]PixelWindowGo.PixelWindow = [...]PixelWindowGo.PixelWindow{
	{H: 0, ThePointer: 0, Title: "RED", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
	{H: 0, ThePointer: 0, Title: "GREEN", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
	{H: 0, ThePointer: 0, Title: "BLUE", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
}

func main() {

	for i := 0; i < 3; i++ {
		PixelWindowGo.CreatePixelWindow(&thewindows[i])
	}

	PixelWindowGo.TheMessagePump(true)
	go func() {
		for true {
			fillBuffersAndUpdatePixels()
		}
	}()

	//go func() {
	{
		PixelWindowGo.TheMessagePump(false)
	}
	//}()
}

var level byte = 100

func fillBuffersAndUpdatePixels() {
	for i := 0; i < 640; i++ {
		for j := 0; j < 480; j++ {
			var offset int64 = int64(4 * (j*640 + i))

			redbuffer[offset+0] = 0
			redbuffer[offset+1] = 0
			redbuffer[offset+2] = level
			redbuffer[offset+3] = 0

			greenbuffer[offset+0] = 0
			greenbuffer[offset+1] = level
			greenbuffer[offset+2] = 0
			greenbuffer[offset+3] = 0

			bluebuffer[offset+0] = level
			bluebuffer[offset+1] = 0
			bluebuffer[offset+2] = 0
			bluebuffer[offset+2] = 0
		}
	}
	level++
	level %= 255
	for i := 0; i < 100; i++ {
		var offset int64 = int64(4 * (i*640 + i))
		redbuffer[offset+0] = 255
		greenbuffer[offset+0] = 255
		bluebuffer[offset+0] = 255
		redbuffer[offset+1] = 255
		greenbuffer[offset+1] = 255
		bluebuffer[offset+1] = 255
		redbuffer[offset+2] = 255
		greenbuffer[offset+2] = 255
		bluebuffer[offset+2] = 255
		redbuffer[offset+3] = 255
		greenbuffer[offset+3] = 255
		bluebuffer[offset+3] = 255
		offset = int64(4 * ((480-i-1)*640 + i))
		redbuffer[offset+0] = 255
		greenbuffer[offset+0] = 255
		bluebuffer[offset+0] = 255
		redbuffer[offset+1] = 255
		greenbuffer[offset+1] = 255
		bluebuffer[offset+1] = 255
		redbuffer[offset+2] = 255
		greenbuffer[offset+2] = 255
		bluebuffer[offset+2] = 255
		redbuffer[offset+3] = 255
		greenbuffer[offset+3] = 255
		bluebuffer[offset+3] = 255

	}

	drawtext(250, 36, 6, 28, &redbuffer[0], "jellyfish")
	drawtext(250, 36, 6, 28, &greenbuffer[0], "bluewindow")
	drawtext(250, 36, 6, 28, &bluebuffer[0], "yahoo")

	thewindows[0].LDAPIXELWindowDisplayBuffer(&redbuffer[0])
	thewindows[1].LDAPIXELWindowDisplayBuffer(&greenbuffer[0])
	thewindows[2].LDAPIXELWindowDisplayBuffer(&bluebuffer[0])

}

func drawtext(width int, height int, startingDotX int, startingDotY int, bfr *byte, whattowrite string) error {

	f, err := opentype.Parse(goitalic.TTF)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	dst := image.NewGray(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	d.DrawString(whattowrite)
	d.Src = image.NewUniform(color.Gray{0x7F})
	d.DrawString(" fish")

	const asciiArt = ".++8"
	buf := make([]byte, 0, height*(width+1))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := asciiArt[dst.GrayAt(x, y).Y>>6]
			if c != '.' {
				// No-op.
				writetextpixel(x, y, 255, 0, 0, bfr)
			} else if x == startingDotX-1 {
				writetextpixel(x, y, 0, 255, 0, bfr)
			} else if y == startingDotY-1 {
				writetextpixel(x, y, 0, 0, 255, bfr)
			}

		}

	}
	os.Stdout.Write(buf)

	return nil
}

func writetextpixel(x int, y int, r byte, g byte, b byte, thebuf *byte) {
	var offseta int64 = int64(4 * (y*textwidthx + x))
	*(*byte)(unsafe.Pointer(
		uintptr(unsafe.Pointer(thebuf)) + uintptr(0+offseta))) = g
	*(*byte)(unsafe.Pointer(
		uintptr(unsafe.Pointer(thebuf)) + uintptr(1+offseta))) = r
	*(*byte)(unsafe.Pointer(
		uintptr(unsafe.Pointer(thebuf)) + uintptr(2+offseta))) = b
	*(*byte)(unsafe.Pointer(
		uintptr(unsafe.Pointer(thebuf)) + uintptr(3+offseta))) = 255
}
