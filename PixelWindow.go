package PixelWindowGo

import (
	"sync"
	"syscall"
	"unsafe"
)

const WS_VISIBLE = 0x10000000
const WS_OVERLAPPEDWINDOW = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
const SW_SHOWDEFAULT = 10
const CW_USEDEFAULT = ^0x7fffffff
const CS_VREDRAW = 0x00000001
const CS_HREDRAW = 0x00000002
const IDI_APPLICATION = 32512
const IDC_ARROW = 32512
const COLOR_WINDOW = 5
const WM_DESTROY = 2

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}
func PostQuitMessage(exitCode int) {
	procPostQuitMessage.Call(
		uintptr(exitCode))
}

func WndProc(hWnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_DESTROY:
		PostQuitMessage(0)
	default:
		return DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
func DefWindowProc(hwnd HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProc.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam)

	return ret
}

func GetModuleHandle(modulename string) HINSTANCE {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return HINSTANCE(ret)
}

var modkernel32 = syscall.NewLazyDLL("kernel32.dll")

var procGetModuleHandle = modkernel32.NewProc("GetModuleHandleW")

var moduser32 = syscall.NewLazyDLL("user32.dll")

var procDefWindowProc = moduser32.NewProc("DefWindowProcW")
var procPostQuitMessage = moduser32.NewProc("PostQuitMessage")
var procUpdateWindow = moduser32.NewProc("UpdateWindow")
var procShowWindow = moduser32.NewProc("ShowWindow")
var procCreateWindowEx = moduser32.NewProc("CreateWindowExW")
var procRegisterClassEx = moduser32.NewProc("RegisterClassExW")
var procLoadCursor = moduser32.NewProc("LoadCursorW")
var procLoadIcon = moduser32.NewProc("LoadIconW")
var procGetMessage = moduser32.NewProc("GetMessageW")
var procTranslateMessage = moduser32.NewProc("TranslateMessage")
var procDispatchMessage = moduser32.NewProc("DispatchMessageW")

type HANDLE uintptr
type HWND HANDLE

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}
type HMENU HANDLE
type DWORD uint32
type HINSTANCE HANDLE
type HICON HANDLE
type HCURSOR HANDLE
type HBRUSH HANDLE
type ATOM uint16

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms633577.aspx
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   *uint16
	ClassName  *uint16
	IconSm     HICON
}

func LoadIcon(instance HINSTANCE, iconName *uint16) HICON {
	ret, _, _ := procLoadIcon.Call(
		uintptr(instance),
		uintptr(unsafe.Pointer(iconName)))

	return HICON(ret)

}

func UpdateWindow(hwnd HWND) bool {
	ret, _, _ := procUpdateWindow.Call(
		uintptr(hwnd))
	return ret != 0
}

func LoadCursor(instance HINSTANCE, cursorName *uint16) HCURSOR {
	ret, _, _ := procLoadCursor.Call(
		uintptr(instance),
		uintptr(unsafe.Pointer(cursorName)))

	return HCURSOR(ret)

}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))

	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret != 0

}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))

	return ret

}

func RegisterClassEx(wndClassEx *WNDCLASSEX) ATOM {
	ret, _, _ := procRegisterClassEx.Call(uintptr(unsafe.Pointer(wndClassEx)))
	return ATOM(ret)
}
func ShowWindow(hwnd HWND, cmdshow int) bool {
	ret, _, _ := procShowWindow.Call(
		uintptr(hwnd),
		uintptr(cmdshow))

	return ret != 0

}

func CreateWindowEx(exStyle uint, className, windowName *uint16,
	style uint, x, y, width, height int, parent HWND, menu HMENU,
	instance HINSTANCE, param unsafe.Pointer) HWND {
	ret, _, _ := procCreateWindowEx.Call(
		uintptr(exStyle),
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		uintptr(style),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(parent),
		uintptr(menu),
		uintptr(instance),
		uintptr(param))

	return HWND(ret)
}

func CreatePixelWindow(pwg *sync.WaitGroup, mytitle string, xpix int, ypix int, isonsync bool, ppw *PixelWindow) {
	if pwg == nil {
		return
	}
	hInstance := GetModuleHandle("")

	lpszClassName := syscall.StringToUTF16Ptr("CN" + mytitle)

	var wcex WNDCLASSEX
	wcex.Size = uint32(unsafe.Sizeof(wcex))
	wcex.Style = CS_HREDRAW | CS_VREDRAW
	wcex.WndProc = syscall.NewCallback(WndProc)
	wcex.ClsExtra = 0
	wcex.WndExtra = 0
	wcex.Instance = hInstance
	wcex.Icon = LoadIcon(hInstance, MakeIntResource(IDI_APPLICATION))
	wcex.Cursor = LoadCursor(0, MakeIntResource(IDC_ARROW))
	wcex.Background = COLOR_WINDOW + 11
	wcex.MenuName = nil
	wcex.ClassName = lpszClassName
	wcex.IconSm = LoadIcon(hInstance, MakeIntResource(IDI_APPLICATION))
	RegisterClassEx(&wcex)

	hWnd := CreateWindowEx(
		0, lpszClassName, syscall.StringToUTF16Ptr(mytitle),
		WS_OVERLAPPEDWINDOW|WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 400, 0, 0, hInstance, nil)

	ShowWindow(hWnd, SW_SHOWDEFAULT)
	UpdateWindow(hWnd)

	theMessagePump()
	pwg.Done()
}

type LDAPIXELWINDOWHANDLE int64

type PixelWindow struct {
	H        HWND
	Apointer int64
}

func (pw *PixelWindow) LDAPIXELWindowDisplayBuffer(
	//win LDAPIXELWINDOWHANDLE,
	bfr *byte) {
	//pw.Apointer = int64(win)
	pw.DisplayBuffer(bfr)
}

const PIXELWINDOW_BUFFER_SIZE = 4
const imgsizebytes = 640 * 480 * 4

var ThePixelBuffer PixelBuffer

type PixelBuffer struct {
	Pixels [PIXELWINDOW_BUFFER_SIZE * imgsizebytes]*byte
	Head   int
	Tail   int
	Size   int
}

func (pw *PixelWindow) DisplayBuffer(b *byte) {
	ThePixelBuffer.Pixels[ThePixelBuffer.Head] = b

	ThePixelBuffer.Head++
	ThePixelBuffer.Head %= PIXELWINDOW_BUFFER_SIZE
	ThePixelBuffer.Size++
}

func theMessagePump() int {
	var msg MSG
	for {
		if GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
	return int(msg.WParam)
}
