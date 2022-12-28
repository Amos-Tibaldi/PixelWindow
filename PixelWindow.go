package PixelWindowGo

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"syscall"
	"unsafe"
)

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

type HGDIOBJ HANDLE

func GetStockObject(fnObject int) HGDIOBJ {
	ret, _, _ := procGetStockObject.Call(
		uintptr(fnObject))

	return HGDIOBJ(ret)
}

func CreatePixelWindow(pwg *sync.WaitGroup, mytitle string, xpix int, ypix int, isonsync bool, ppw *PixelWindow) {
	if pwg == nil {
		return
	}
	hInstance := GetModuleHandle("")

	lpszClassName := syscall.StringToUTF16Ptr("CN" + mytitle)

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
		0, lpszClassName, syscall.StringToUTF16Ptr(mytitle),
		WS_OVERLAPPEDWINDOW|WS_VISIBLE|WS_SYSMENU|WS_MINIMIZEBOX,
		0, 0, xpix, ypix, 0, 0, hInstance, nil)

	SetWindowLongPtr(hWnd, GWLP_USERDATA, uintptr(unsafe.Pointer(ppw)))

	g_D3D, theerr := Create(D3D_SDK_VERSION)
	fmt.Println(theerr)
	var pp PRESENT_PARAMETERS
	pp.BackBufferCount = 1
	pp.BackBufferWidth = uint32(xpix)
	pp.BackBufferHeight = uint32(ypix)
	pp.MultiSampleType = MULTISAMPLE_NONE
	pp.MultiSampleQuality = 0
	pp.SwapEffect = SWAPEFFECT_DISCARD
	pp.HDeviceWindow = hWnd
	pp.Windowed = 1
	pp.Flags = PRESENTFLAG_LOCKABLE_BACKBUFFER
	pp.FullScreen_RefreshRateInHz = D3DPRESENT_RATE_DEFAULT
	const D3DPRESENT_INTERVAL_DEFAULT = 0x00000000
	const D3DPRESENT_INTERVAL_IMMEDIATE = 0x80000000
	if isonsync {
		pp.PresentationInterval = D3DPRESENT_INTERVAL_DEFAULT
	} else {
		pp.PresentationInterval = D3DPRESENT_INTERVAL_IMMEDIATE
	}
	pp.BackBufferFormat = 22      //Display format
	pp.EnableAutoDepthStencil = 0 //No depth/stencil buffer
	const D3DADAPTER_DEFAULT = 0
	const D3DDEVTYPE_HAL = 1
	const D3DCREATE_HARDWARE_VERTEXPROCESSING = 0x00000040
	var p_device *Device
	p_device, _, _ = g_D3D.CreateDevice(D3DADAPTER_DEFAULT,
		D3DDEVTYPE_HAL,
		hWnd,
		D3DCREATE_HARDWARE_VERTEXPROCESSING, // D3DCREATE_SOFTWARE_VERTEXPROCESSING,
		pp)

	var FrontBuffer Surface
	var prova = uintptr(unsafe.Pointer(&FrontBuffer))
	p_device.CreateOffscreenPlainSurface(
		uint(xpix),
		uint(ypix),
		22,
		0, //D3DPOOL_SYSTEMMEM,
		prova,
	)

	//ResizeWindow(xpix, ypix)

	const SW_SHOW = 5
	ShowWindow(hWnd, SW_SHOW)
	UpdateWindow(hWnd)

	//CreateThread(NULL, NULL, (LPTHREAD_START_ROUTINE)&PixelWindow::run, this, NULL, NULL);

	theMessagePump()
	pwg.Done()
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

func CalculateExactRect(a int, b int, r RECT) {

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

func (ppw *PixelWindow) ResizeWindow(width, height int) {
	var rect RECT = RECT{0, 0, 0, 0}
	CalculateExactRect(width, height, rect)
	MoveWindow(ppw.H,
		100, 100,
		200,
		200,
		true)
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
