package main

import (
	"PixelWindowGo/PixelWindowGo"
	"sync"
)

var g_app_done bool
var (
	redbuffer   *byte
	greenbuffer *byte
	bluebuffer  *byte
)

func main() {
	const imgsizebytes = 640 * 480 * 4
	var redbuffer [imgsizebytes]byte
	var greenbuffer [imgsizebytes]byte
	var bluebuffer [imgsizebytes]byte
	var names [3]string = [...]string{"RED", "GREEN", "BLUE"}
	var thewindows [3]PixelWindowGo.PixelWindow = [...]PixelWindowGo.PixelWindow{{0, 0}, {0, 0}, {0, 0}}
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go PixelWindowGo.CreatePixelWindow(&wg, names[i], 640, 480, true, &thewindows[i])
	}

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
	wg.Wait()
}
