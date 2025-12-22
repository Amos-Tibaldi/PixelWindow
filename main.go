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
)

var g_app_done bool

func main() {
	const imgsizebytes = 640 * 480 * 4
	var redbuffer [imgsizebytes]byte
	var greenbuffer [imgsizebytes]byte
	var bluebuffer [imgsizebytes]byte
	var thewindows [3]PixelWindowGo.PixelWindow = [...]PixelWindowGo.PixelWindow{
		{H: 0, ThePointer: 0, Title: "RED", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
		{H: 0, ThePointer: 0, Title: "GREEN", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
		{H: 0, ThePointer: 0, Title: "BLUE", Xpixsize: 640, Ypixsize: 480, VSync: true, Width: 640, Height: 480},
	}

	for i := 0; i < 3; i++ {
		PixelWindowGo.CreatePixelWindow(&thewindows[i])
	}

	go func() {
		var level byte = 100

		for !g_app_done {
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

			thewindows[0].LDAPIXELWindowDisplayBuffer(&redbuffer[0])
			thewindows[1].LDAPIXELWindowDisplayBuffer(&greenbuffer[0])
			thewindows[2].LDAPIXELWindowDisplayBuffer(&bluebuffer[0])
		}
	}()

	PixelWindowGo.TheMessagePump()
}
