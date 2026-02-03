//-----------------------------------------------------------------------
//
// This file is part of the PixelWindow Project
//
//  by Amos Tibaldi - tibaldi at users.sourceforge.net
//
// https://sourceforge.net/projects/pixelwindow/
//
// https://github.com/Amos-Tibaldi/PixelWindow
//
//
// COPYRIGHT: http://www.gnu.org/licenses/gpl.html
//            COPYRIGHT-gpl-3.0.txt
//
//     The PixelWindow Project
//        PixelWindow gives high performance pixel access to DirectX windows
//     in go and in c++.
//
//     Copyright (C) 2022 Amos Tibaldi
//
//     This program is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     This program is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
//-----------------------------------------------------------------------

package PixelWindowGo

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

const SW_SHOW = 5
const D3D_SDK_VERSION = 32
const GWLP_USERDATA = -21
const WS_VISIBLE = 0x10000000
const WS_OVERLAPPEDWINDOW = 0x00000000 | 0x00C00000 | 0x00080000 | 0x00040000 | 0x00020000 | 0x00010000
const SW_SHOWDEFAULT = 10
const CW_USEDEFAULT = ^0x7fffffff
const CS_VREDRAW = 0x00000001
const CS_HREDRAW = 0x00000002
const CS_OWNDC = 0x00000020
const BLACK_BRUSH = 4
const IDI_APPLICATION = 32512
const IDC_ARROW = 32512
const COLOR_WINDOW = 5
const WM_DESTROY = 2
const WS_SYSMENU = 0x00080000
const WS_MINIMIZEBOX = 0x00020000
const MULTISAMPLE_NONE = 0

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

var procSetWindowLongPtr = moduser32.NewProc("SetWindowLongW")
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
var procMoveWindow = moduser32.NewProc("MoveWindow")
var procGetWindowRect = moduser32.NewProc("GetWindowRect")
var procGetClientRect = moduser32.NewProc("GetClientRect")

var modgdi32 = syscall.NewLazyDLL("gdi32.dll")

var procGetStockObject = modgdi32.NewProc("GetStockObject")

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

func SetWindowLongPtr(hwnd HWND, index int, value uintptr) uintptr {
	ret, _, _ := procSetWindowLongPtr.Call(
		uintptr(hwnd),
		uintptr(index),
		value)

	return ret
}

func LoadIcon(instance HINSTANCE, iconName *uint16) HICON {
	ret, _, _ := procLoadIcon.Call(
		uintptr(instance),
		uintptr(unsafe.Pointer(iconName)))

	return HICON(ret)

}

func UpdateWindow(hwnd HWND) bool {
	retuw, retuwptr, erroruw := procUpdateWindow.Call(uintptr(hwnd))
	if retuwptr != 0 {

	}
	if erroruw != nil {
		//fmt.Println("prntfuw %d", erroruw)
	}
	return retuw != 0
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

type HGDIOBJ HANDLE

func GetStockObject(fnObject int) HGDIOBJ {
	ret, _, _ := procGetStockObject.Call(
		uintptr(fnObject))

	return HGDIOBJ(ret)
}

func CreatePixelWindow(ppw *PixelWindow) {
	ppw.MYBUF.Bufcondvar = sync.NewCond(&ppw.MYBUF.Bufmutex)

	hInstance := GetModuleHandle("")

	lpszClassName := syscall.StringToUTF16Ptr("CN" + ppw.Title)

	ppw.IsRectUsed = false
	ppw.Rect.Bottom = 0
	ppw.Rect.Left = 0
	ppw.Rect.Top = 0
	ppw.Rect.Right = 0

	var wcex WNDCLASSEX
	wcex.Size = uint32(unsafe.Sizeof(wcex))
	wcex.Style = CS_OWNDC
	wcex.WndProc = syscall.NewCallback(WndProc)
	wcex.ClsExtra = 0
	wcex.WndExtra = 0
	wcex.Instance = hInstance
	wcex.Icon = LoadIcon(hInstance, MakeIntResource(IDI_APPLICATION))
	wcex.Cursor = LoadCursor(0, MakeIntResource(IDC_ARROW))
	wcex.Background = HBRUSH(GetStockObject(BLACK_BRUSH))
	wcex.MenuName = nil
	wcex.ClassName = lpszClassName
	wcex.IconSm = LoadIcon(hInstance, MakeIntResource(IDI_APPLICATION))
	RegisterClassEx(&wcex)

	hWnd := CreateWindowEx(
		0, lpszClassName, syscall.StringToUTF16Ptr(ppw.Title),
		WS_OVERLAPPEDWINDOW|WS_VISIBLE|WS_SYSMENU|WS_MINIMIZEBOX,
		0, 0, ppw.Xpixsize, ppw.Ypixsize, 0, 0, hInstance, nil)
	ppw.H = hWnd
	SetWindowLongPtr(hWnd, GWLP_USERDATA, uintptr(unsafe.Pointer(ppw)))

	g_D3D, theerr := Create(D3D_SDK_VERSION)
	fmt.Sprintln("graphics %d", theerr)
	var pp PRESENT_PARAMETERS
	pp.BackBufferCount = 1
	pp.BackBufferWidth = uint32(ppw.Xpixsize)
	pp.BackBufferHeight = uint32(ppw.Ypixsize)
	pp.MultiSampleType = MULTISAMPLE_NONE
	pp.MultiSampleQuality = 0
	pp.SwapEffect = SWAPEFFECT_DISCARD
	pp.HDeviceWindow = ppw.H
	pp.Windowed = 1
	pp.Flags = PRESENTFLAG_LOCKABLE_BACKBUFFER
	pp.FullScreen_RefreshRateInHz = D3DPRESENT_RATE_DEFAULT
	const D3DPRESENT_INTERVAL_DEFAULT = 0x00000000
	const D3DPRESENT_INTERVAL_IMMEDIATE = 0x80000000
	if ppw.VSync {
		pp.PresentationInterval = D3DPRESENT_INTERVAL_DEFAULT
	} else {
		pp.PresentationInterval = D3DPRESENT_INTERVAL_IMMEDIATE
	}
	const D3DFMT_X8R8G8B8 = 22
	pp.BackBufferFormat = D3DFMT_X8R8G8B8 //Display format
	pp.EnableAutoDepthStencil = 0         //No depth/stencil buffer
	const D3DADAPTER_DEFAULT = 0
	const D3DDEVTYPE_HAL = 1
	const D3DCREATE_HARDWARE_VERTEXPROCESSING = 0x00000040
	//var errorg3d error = nil
	ppw.P_device, _, _ = g_D3D.CreateDevice(D3DADAPTER_DEFAULT,
		D3DDEVTYPE_HAL,
		hWnd,
		D3DCREATE_HARDWARE_VERTEXPROCESSING, // D3DCREATE_SOFTWARE_VERTEXPROCESSING,
		pp)
	//fmt.Printf("errorg3d %s", errorg3d)
	//var err error
	const D3DPOOL_SYSTEMMEM = 2
	ppw.PFrontBuffer, _ = ppw.P_device.CreateOffscreenPlainSurface(
		uint(ppw.Xpixsize),
		uint(ppw.Ypixsize),
		D3DFMT_X8R8G8B8,
		D3DPOOL_SYSTEMMEM,
		0, //uintptr(puntfrontbuffer),
	)
	//fmt.Println(err)
	const D3DBACKBUFFER_TYPE_MONO = 0
	ppw.PBackBuffer, _ = ppw.P_device.GetBackBuffer(0, 0, D3DBACKBUFFER_TYPE_MONO)
	//fmt.Println(err)
	ppw.ResizeWindow(ppw.Width, ppw.Height)

	ShowWindow(ppw.H, SW_SHOW)
	UpdateWindow(ppw.H)

	go pixwinthread(ppw)
}

type BACKBUFFER_TYPE uint32

// GetBackBuffer retrieves a back buffer from the device's swap chain.
// Call Release on the returned surface when finished using it.
func (obj *Device) GetBackBuffer(
	swapChain uint,
	backBuffer uint,
	typ BACKBUFFER_TYPE,
) (*Surface, Error) {
	var surface *Surface
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.GetBackBuffer,
		5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(swapChain),
		uintptr(backBuffer),
		uintptr(typ),
		uintptr(unsafe.Pointer(&surface)),
		0,
	)
	return surface, toErr(ret)
}

func (pwp *PixelWindow) ResizeWindow(width int, height int) {
	for ; !pwp.IsRectUsed; pwp.IsRectUsed = true {
		fmt.Sprintf("in Resizewindow %d \n", pwp.H)
		pwp.CalculateExactRect(int32(width), int32(height), &pwp.Rect)
		MoveWindow(pwp.H,
			int(pwp.Rect.Left), int(pwp.Rect.Top),
			int(pwp.Rect.Right)-int(pwp.Rect.Left),
			int(pwp.Rect.Bottom)-int(pwp.Rect.Top),
			true)
	}
	pwp.IsRectUsed = false
}

func pixwinthread(pwp *PixelWindow) {
	for true {
		//fmt.Println("PWT start pixwinthread before lock size", pwp.MYBUF.MyPixelBuffer.Size)
		pwp.MYBUF.Bufmutex.Lock()
		//fmt.Println("PWT start pixwinthread after lock")
		for pwp.MYBUF.MyPixelBuffer.Size == 0 {
			//fmt.Println("PWT thrinnerH sizeis ", pwp.H, pwp.MYBUF.MyPixelBuffer.Size)
			pwp.MYBUF.Bufcondvar.Wait()
			//fmt.Println("PWT exit pixwinthread condwait size", pwp.H, pwp.MYBUF.MyPixelBuffer.Size)
		}
		pwp.CopyFrameToFrontBuffer()
		pwp.MYBUF.MyPixelBuffer.Tail++
		pwp.MYBUF.MyPixelBuffer.Tail %= PIXELWINDOW_BUFFER_SIZE
		pwp.MYBUF.MyPixelBuffer.Size--
		//fmt.Println("PWT fine copy size= ", pwp.MYBUF.MyPixelBuffer.Size)
		if pwp.MYBUF.MyPixelBuffer.Size > 0 {
			pwp.MYBUF.Bufcondvar.Signal()
		}
		pwp.PutFrontBufferOntoScreen()
		pwp.MYBUF.Bufmutex.Unlock()
		//fmt.Println("PWT end pwt HHH  Size ", pwp.H, pwp.MYBUF.MyPixelBuffer.Size)
	}
}

type COLOR uint32

// Clear clears one or more surfaces such as a render target, multiple render
// targets, a stencil buffer, and a depth buffer.
func (obj *Device) Clear(
	rects []RECT,
	flags uint32,
	color COLOR,
	z float32,
	stencil uint32,
) Error {
	var rectPtr *RECT
	if len(rects) > 0 {
		rectPtr = &rects[0]
	}
	ret, _, _ := syscall.Syscall9(
		obj.vtbl.Clear,
		7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(len(rects)),
		uintptr(unsafe.Pointer(rectPtr)),
		uintptr(flags),
		uintptr(color),
		uintptr(z),
		uintptr(stencil),
		0,
		0,
	)
	return toErr(ret)
}

// BeginScene begins a scene.
// Applications must call BeginScene before performing any rendering and must
// call EndScene when rendering is complete and before calling BeginScene again.
func (obj *Device) BeginScene() Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.BeginScene,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return toErr(ret)
}

// EndScene ends a scene that was begun by calling BeginScene.
func (obj *Device) EndScene() (err Error) {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.EndScene,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return toErr(ret)
}

// Present presents the contents of the next buffer in the sequence of back
// buffers owned by the device.
func (obj *Device) Present(
	sourceRect *RECT,
	destRect *RECT,
	destWindowOverride HWND,
	dirtyRegion *RGNDATA,
) Error {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.Present,
		5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(sourceRect)),
		uintptr(unsafe.Pointer(destRect)),
		uintptr(destWindowOverride),
		uintptr(unsafe.Pointer(dirtyRegion)),
		0,
	)
	return toErr(ret)
}

const NULL = 0
const D3DCLEAR_TARGET = 0x00000001 /* Clear target surface */
// RGNDATA contains region data.
type RGNDATA struct {
	Rdh    RGNDATAHEADER
	Buffer [1]byte
}

// RGNDATAHEADER describes region data.
type RGNDATAHEADER struct {
	DwSize   uint32
	IType    uint32
	NCount   uint32
	NRgnSize uint32
	RcBound  RECT
}

func (ppw *PixelWindow) PutFrontBufferOntoScreen() {
	ppw.P_device.Clear( //Number of rectangles to clear, we're clearing everything so set it to 0
		nil,             //Pointer to the rectangles to clear, NULL to clear whole display
		D3DCLEAR_TARGET, //What to clear.  We don't have a Z Buffer or Stencil Buffer
		0x00000000,      //Colour to clear to (AARRGGBB)
		1.0,             //Value to clear ZBuffer to, doesn't matter since we don't have one
		0)               //Stencil clear value, again, we don't have one, this value doesn't matter
	ppw.P_device.BeginScene()
	ppw.P_device.UpdateSurface(ppw.PFrontBuffer, nil, ppw.PBackBuffer, nil)
	ppw.P_device.EndScene()
	ppw.P_device.Present(nil, //Source rectangle to display, NULL for all of it
		nil,  //Destination rectangle, NULL to fill whole display
		NULL, //Target window, if NULL uses device window set in CreateDevice
		nil)  //Dirty Region, set it to NULL
}

// UpdateSurface copies rectangular subsets of pixels from one surface to another.
func (obj *Device) UpdateSurface(
	sourceSurface *Surface,
	sourceRect *RECT,
	destSurface *Surface,
	destPoint *POINT,
) Error {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.UpdateSurface,
		5,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(sourceSurface)),
		uintptr(unsafe.Pointer(sourceRect)),
		uintptr(unsafe.Pointer(destSurface)),
		uintptr(unsafe.Pointer(destPoint)),
		0,
	)
	return toErr(ret)
}

// LOCKED_RECT describes a locked rectangular region.
type LOCKED_RECT struct {
	Pitch int32
	PBits uintptr
}

// UnlockRect unlocks a rectangle on a surface.
func (obj *Surface) UnlockRect() Error {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.UnlockRect,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return toErr(ret)
}

// LockRect locks a rectangle on a surface.
func (obj *Surface) LockRect(
	rect *RECT,
	flags uint32,
) (lockedRect LOCKED_RECT, err Error) {
	ret, _, _ := syscall.Syscall6(
		obj.vtbl.LockRect,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&lockedRect)),
		uintptr(unsafe.Pointer(rect)),
		uintptr(flags),
		0,
		0,
	)
	err = toErr(ret)
	return
}

func (ppw *PixelWindow) CopyFrameToFrontBuffer() {
	ppw.TheLockedR, _ = ppw.PFrontBuffer.LockRect(nil, 0)
	var a byte
	for i := 0; i < ppw.Xpixsize*ppw.Ypixsize*4; i++ {
		a = *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&ppw.MYBUF.MyPixelBuffer.Pixels[ppw.MYBUF.MyPixelBuffer.Tail][0])) + uintptr(i)))
		*(*uint8)(unsafe.Pointer(ppw.TheLockedR.PBits + uintptr(i))) = a
		a++
	}
	ppw.PFrontBuffer.UnlockRect()
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

func GetClientRect(hwnd HWND) *RECT {
	var rect RECT
	ret, _, _ := procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		panic(fmt.Sprintf("GetClientRect(%d) failed", hwnd))
	}

	return &rect
}

func GetWindowRect(hwnd HWND) *RECT {
	var rect RECT
	procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	return &rect
}

func (p *PixelWindow) CalculateExactRect(cx int32, cy int32, rect *RECT) {
	var rcWindow *RECT = GetWindowRect(p.H)
	var rcClient *RECT = GetClientRect(p.H)
	cx += (rcWindow.Right - rcWindow.Left) - rcClient.Right
	cy += (rcWindow.Bottom - rcWindow.Top) - rcClient.Bottom
	rect.Left = rcWindow.Left
	rect.Top = rcWindow.Top
	rect.Right = rect.Left + cx
	rect.Bottom = rect.Top + cy
}

type BOOL int

func BoolToBOOL(b bool) BOOL {
	if b {
		return 1
	}
	return 0
}
func MoveWindow(hwnd HWND, x, y, width, height int, repaint bool) bool {
	ret, _, _ := procMoveWindow.Call(
		uintptr(hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(BoolToBOOL(repaint)))

	return ret != 0

}

type POOL uint32
type Surface struct {
	vtbl *surfaceVtbl
}

type surfaceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetDevice       uintptr
	SetPrivateData  uintptr
	GetPrivateData  uintptr
	FreePrivateData uintptr
	SetPriority     uintptr
	GetPriority     uintptr
	PreLoad         uintptr
	GetType         uintptr
	GetContainer    uintptr
	GetDesc         uintptr
	LockRect        uintptr
	UnlockRect      uintptr
	GetDC           uintptr
	ReleaseDC       uintptr
}

// CreateOffscreenPlainSurface creates an off-screen surface.
func (obj *Device) CreateOffscreenPlainSurface(
	width uint,
	height uint,
	format FORMAT,
	pool POOL,
	sharedHandle uintptr,
) (*Surface, Error) {
	var surface *Surface
	ret, _, _ := syscall.Syscall9(
		obj.vtbl.CreateOffscreenPlainSurface,
		7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(width),
		uintptr(height),
		uintptr(format),
		uintptr(pool),
		uintptr(unsafe.Pointer(&surface)),
		sharedHandle,
		0,
		0,
	)
	return surface, toErr(ret)
}

const D3DPRESENT_RATE_DEFAULT = 0x00000000

const PRESENTFLAG_LOCKABLE_BACKBUFFER = 0x00000001
const SWAPEFFECT_DISCARD = 1

type MULTISAMPLE_TYPE uint32
type SWAPEFFECT uint32
type FORMAT uint32 // ??? BOH

type PRESENT_PARAMETERS struct {
	BackBufferWidth            uint32
	BackBufferHeight           uint32
	BackBufferFormat           FORMAT
	BackBufferCount            uint32
	MultiSampleType            MULTISAMPLE_TYPE
	MultiSampleQuality         uint32
	SwapEffect                 SWAPEFFECT
	HDeviceWindow              HWND
	Windowed                   int32
	EnableAutoDepthStencil     int32
	AutoDepthStencilFormat     FORMAT
	Flags                      uint32
	FullScreen_RefreshRateInHz uint32
	PresentationInterval       uint32
}

func Create(version uint) (*Direct3D, error) {
	obj, _, _ := direct3DCreate9.Call(uintptr(version))
	if obj == 0 {
		return nil, errors.New("Direct3DCreate9 returned nil")
	}
	return (*Direct3D)(unsafe.Pointer(obj)), nil
}

type Device struct {
	vtbl *deviceVtbl
}

type deviceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	TestCooperativeLevel        uintptr
	GetAvailableTextureMem      uintptr
	EvictManagedResources       uintptr
	GetDirect3D                 uintptr
	GetDeviceCaps               uintptr
	GetDisplayMode              uintptr
	GetCreationParameters       uintptr
	SetCursorProperties         uintptr
	SetCursorPosition           uintptr
	ShowCursor                  uintptr
	CreateAdditionalSwapChain   uintptr
	GetSwapChain                uintptr
	GetNumberOfSwapChains       uintptr
	Reset                       uintptr
	Present                     uintptr
	GetBackBuffer               uintptr
	GetRasterStatus             uintptr
	SetDialogBoxMode            uintptr
	SetGammaRamp                uintptr
	GetGammaRamp                uintptr
	CreateTexture               uintptr
	CreateVolumeTexture         uintptr
	CreateCubeTexture           uintptr
	CreateVertexBuffer          uintptr
	CreateIndexBuffer           uintptr
	CreateRenderTarget          uintptr
	CreateDepthStencilSurface   uintptr
	UpdateSurface               uintptr
	UpdateTexture               uintptr
	GetRenderTargetData         uintptr
	GetFrontBufferData          uintptr
	StretchRect                 uintptr
	ColorFill                   uintptr
	CreateOffscreenPlainSurface uintptr
	SetRenderTarget             uintptr
	GetRenderTarget             uintptr
	SetDepthStencilSurface      uintptr
	GetDepthStencilSurface      uintptr
	BeginScene                  uintptr
	EndScene                    uintptr
	Clear                       uintptr
	SetTransform                uintptr
	GetTransform                uintptr
	MultiplyTransform           uintptr
	SetViewport                 uintptr
	GetViewport                 uintptr
	SetMaterial                 uintptr
	GetMaterial                 uintptr
	SetLight                    uintptr
	GetLight                    uintptr
	LightEnable                 uintptr
	GetLightEnable              uintptr
	SetClipPlane                uintptr
	GetClipPlane                uintptr
	SetRenderState              uintptr
	GetRenderState              uintptr
	CreateStateBlock            uintptr
	BeginStateBlock             uintptr
	EndStateBlock               uintptr
	SetClipStatus               uintptr
	GetClipStatus               uintptr
	GetTexture                  uintptr
	SetTexture                  uintptr
	GetTextureStageState        uintptr
	SetTextureStageState        uintptr
	GetSamplerState             uintptr
	SetSamplerState             uintptr
	ValidateDevice              uintptr
	SetPaletteEntries           uintptr
	GetPaletteEntries           uintptr
	SetCurrentTexturePalette    uintptr
	GetCurrentTexturePalette    uintptr
	SetScissorRect              uintptr
	GetScissorRect              uintptr
	SetSoftwareVertexProcessing uintptr
	GetSoftwareVertexProcessing uintptr
	SetNPatchMode               uintptr
	GetNPatchMode               uintptr
	DrawPrimitive               uintptr
	DrawIndexedPrimitive        uintptr
	DrawPrimitiveUP             uintptr
	DrawIndexedPrimitiveUP      uintptr
	ProcessVertices             uintptr
	CreateVertexDeclaration     uintptr
	SetVertexDeclaration        uintptr
	GetVertexDeclaration        uintptr
	SetFVF                      uintptr
	GetFVF                      uintptr
	CreateVertexShader          uintptr
	SetVertexShader             uintptr
	GetVertexShader             uintptr
	SetVertexShaderConstantF    uintptr
	GetVertexShaderConstantF    uintptr
	SetVertexShaderConstantI    uintptr
	GetVertexShaderConstantI    uintptr
	SetVertexShaderConstantB    uintptr
	GetVertexShaderConstantB    uintptr
	SetStreamSource             uintptr
	GetStreamSource             uintptr
	SetStreamSourceFreq         uintptr
	GetStreamSourceFreq         uintptr
	SetIndices                  uintptr
	GetIndices                  uintptr
	CreatePixelShader           uintptr
	SetPixelShader              uintptr
	GetPixelShader              uintptr
	SetPixelShaderConstantF     uintptr
	GetPixelShaderConstantF     uintptr
	SetPixelShaderConstantI     uintptr
	GetPixelShaderConstantI     uintptr
	SetPixelShaderConstantB     uintptr
	GetPixelShaderConstantB     uintptr
	DrawRectPatch               uintptr
	DrawTriPatch                uintptr
	DeletePatch                 uintptr
	CreateQuery                 uintptr
}

type DEVTYPE uint32

// Error is returned by all Direct3D9 functions. It encapsulates the error code
// returned by Direct3D. If a function succeeds it will return nil as the Error
// and if it fails you can retrieve the error code using the Code() function.
// You can check the result against the predefined error codes (like
// ERR_DEVICELOST, E_OUTOFMEMORY etc).
type Error interface {
	error
	// Code returns the Direct3D error code for a function. Call this function
	// only if the Error is not nil, if the error code is D3D_OK or any other
	// code that signifies success, a function will return nil as the Error
	// instead of a non-nil error with that code in it. This way, functions
	// behave in a standard Go way, returning nil as the error in case of
	// success and only returning non-nil errors if something went wrong.
	Code() int32
}

// CreateDevice creates a device to represent the display adapter.
func (obj *Direct3D) CreateDevice(
	adapter uint,
	deviceType DEVTYPE,
	focusWindow HWND,
	behaviorFlags uint32,
	params PRESENT_PARAMETERS,
) (*Device, PRESENT_PARAMETERS, Error) {
	var device *Device
	ret, _, _ := syscall.Syscall9(
		obj.vtbl.CreateDevice,
		7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(adapter),
		uintptr(deviceType),
		uintptr(focusWindow),
		uintptr(behaviorFlags),
		uintptr(unsafe.Pointer(&params)),
		uintptr(unsafe.Pointer(&device)),
		0,
		0,
	)
	return device, params, toErr(ret)
}

func toErr(result uintptr) Error {
	res := hResultError(result) // cast to signed int
	if res >= 0 {
		return nil
	}
	return res
}
func (r hResultError) Code() int32 { return int32(r) }

func (r hResultError) Error() string {
	switch r {
	/*	case ERR_CONFLICTINGRENDERSTATE:
			return "conflicting render state"
		case ERR_CONFLICTINGTEXTUREFILTER:
			return "conflicting texture filter"
		case ERR_CONFLICTINGTEXTUREPALETTE:
			return "conflicting texture palette"
		case ERR_DEVICEHUNG:
			return "device hung"
		case ERR_DEVICELOST:
			return "device lost"
		case ERR_DEVICENOTRESET:
			return "device not reset"
		case ERR_DEVICEREMOVED:
			return "device removed"
		case ERR_DRIVERINTERNALERROR:
			return "driver internal error"
		case ERR_DRIVERINVALIDCALL:
			return "driver invalid call"
		case ERR_INVALIDCALL:
			return "invalid call"
		case ERR_INVALIDDEVICE:
			return "invalid device"
		case ERR_MOREDATA:
			return "more data"
		case ERR_NOTAVAILABLE:
			return "not available"
		case ERR_NOTFOUND:
			return "not found"
		case ERR_OUTOFVIDEOMEMORY:
			return "out of video memory"
		case ERR_TOOMANYOPERATIONS:
			return "too many operations"
		case ERR_UNSUPPORTEDALPHAARG:
			return "unsupported alpha argument"
		case ERR_UNSUPPORTEDALPHAOPERATION:
			return "unsupported alpha operation"
		case ERR_UNSUPPORTEDCOLORARG:
			return "unsupported color argument"
		case ERR_UNSUPPORTEDCOLOROPERATION:
			return "unsupported color operation"
		case ERR_UNSUPPORTEDFACTORVALUE:
			return "unsupported factor value"
		case ERR_UNSUPPORTEDTEXTUREFILTER:
			return "unsupported texture filter"
		case ERR_WASSTILLDRAWING:
			return "was still drawing"
		case ERR_WRONGTEXTUREFORMAT:
			return "wrong texture format"
		case ERR_UNSUPPORTEDOVERLAY:
			return "unsupported overlay"
		case ERR_UNSUPPORTEDOVERLAYFORMAT:
			return "unsupported overlay format"
		case ERR_CANNOTPROTECTCONTENT:
			return "cannot protect content"
		case ERR_UNSUPPORTEDCRYPTO:
			return "unsupported crypto"

		case E_FAIL:
			return "fail"
		case E_INVALIDARG:
			return "invalid argument"
		case E_NOINTERFACE:
			return "no interface"
		case E_NOTIMPL:
			return "not implemented"
		case E_OUTOFMEMORY:
			return "out of memory"

		case S_NOT_RESIDENT:
			return "not resident"
		case S_RESIDENT_IN_SHARED_MEMORY:
			return "resident in shared memory"
	*/
	default:
		return "unknown error code " + strconv.Itoa(int(r))
	}
}

type hResultError int32

type Direct3D struct {
	vtbl *direct3DVtbl
}

type direct3DVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	RegisterSoftwareDevice      uintptr
	GetAdapterCount             uintptr
	GetAdapterIdentifier        uintptr
	GetAdapterModeCount         uintptr
	EnumAdapterModes            uintptr
	GetAdapterDisplayMode       uintptr
	CheckDeviceType             uintptr
	CheckDeviceFormat           uintptr
	CheckDeviceMultiSampleType  uintptr
	CheckDepthStencilMatch      uintptr
	CheckDeviceFormatConversion uintptr
	GetDeviceCaps               uintptr
	GetAdapterMonitor           uintptr
	CreateDevice                uintptr
}

var (
	dll             = syscall.NewLazyDLL("d3d9.dll")
	direct3DCreate9 = dll.NewProc("Direct3DCreate9")
)

type LDAPIXELWINDOWHANDLE int64

type MyBuffer struct {
	MyPixelBuffer PixelBuffer
	Bufmutex      sync.RWMutex
	Bufcondvar    *sync.Cond
}

type PixelWindow struct {
	H            HWND
	ThePointer   uintptr
	Title        string
	Xpixsize     int
	Ypixsize     int
	VSync        bool
	Width        int
	Height       int
	MYBUF        MyBuffer
	TheLockedR   LOCKED_RECT
	PFrontBuffer *Surface
	PBackBuffer  *Surface
	P_device     *Device
	Rect         RECT
	IsRectUsed   bool
}

func (pw *PixelWindow) LDAPIXELWindowDisplayBuffer(bfr *byte) {
	pw.DisplayBuffer(bfr)
}

const PIXELWINDOW_BUFFER_SIZE = 4
const maximgsizebytes = 640 * 480 * 4

type PixelBuffer struct {
	Pixels [PIXELWINDOW_BUFFER_SIZE][maximgsizebytes]byte
	Head   int
	Tail   int
	Size   int
}

func (pw *PixelWindow) DisplayBuffer(b *byte) {
	//fmt.Println("DisplayBuffer start")
	pw.MYBUF.Bufmutex.Lock()
	//fmt.Println("DisplayBuffer dbmutlock afterlock", pw.Height, pw.Width, pw.MYBUF.MyPixelBuffer.Head)
	for pw.MYBUF.MyPixelBuffer.Size == PIXELWINDOW_BUFFER_SIZE {
		//fmt.Println("DisplayBuffer dbmutlock size %d", pw.MYBUF.MyPixelBuffer.Size)
		pw.MYBUF.Bufcondvar.Wait()
	}
	//fmt.Println("DisplayBuffer after for")
	for i := 0; i < pw.Height*pw.Width*4; i++ {
		pw.MYBUF.MyPixelBuffer.Pixels[pw.MYBUF.MyPixelBuffer.Head][i] =
			*(*byte)(unsafe.Pointer(
				uintptr(unsafe.Pointer(b)) + uintptr(i)))
	}
	//fmt.Println("DisplayBuffer before unlock")

	pw.MYBUF.Bufmutex.Unlock()
	pw.MYBUF.MyPixelBuffer.Head++
	pw.MYBUF.MyPixelBuffer.Head %= PIXELWINDOW_BUFFER_SIZE
	pw.MYBUF.MyPixelBuffer.Size++
	//fmt.Println("DisplayBuffer after size++")
	if pw.MYBUF.MyPixelBuffer.Size > 0 {
		//fmt.Println("DisplayBuffer prima di signal")
		pw.MYBUF.Bufcondvar.Signal()
	}
	//fmt.Println("DisplayBuffer dbmutlock sizebeforeunlock ", pw.MYBUF.MyPixelBuffer.Size, pw.H)

}

func TheMessagePump(exit bool) int {
	var msgg MSG
	for true {
		//fmt.Println("TheMessagePump begin")
		var retval = GetMessage(&msgg, 0, 0, 0)
		//fmt.Println("TheMessagePump Getmsg retval MSG", retval, msgg)
		if retval == 0 {
			//fmt.Println("TheMessagePump Getmsgxit retvalZero Hwprm=", retval, msgg.WParam)
			break
		} else {
			//fmt.Println("TheMessagePump GetmsgxitT retval= Hwprm=", retval, msgg.WParam)

		}
		TranslateMessage(&msgg)
		DispatchMessage(&msgg)
		if exit {
			return int(msgg.WParam)
		}
		//fmt.Println("TheMessagePump Xit  retval= Hwprm=", retval, msgg.WParam)
	}
	return int(msgg.WParam)
}
