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
//     Copyright (C) 20223 Amos Tibaldi
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

#pragma comment( lib, "d3d9" )

#include <Windows.h>
#include <iostream>
#include <stdio.h>
#include <d3d9.h>

#include "PixelWindow.h"

using namespace std;

extern bool g_app_done;

void LDAPIXELWindowDisplayBuffer(LDAPIXELWINDOWHANDLE win, unsigned char* bfr)
{
	PixelWindowNS::PixelWindow * p = (PixelWindowNS::PixelWindow *) win;
	p->DisplayBuffer(bfr);
}

void LDACallWhenIdle()
{
	MSG msg;
	for (int i = 0; i < 20; i++)
	{
		if (PeekMessage(&msg, NULL, 0, 0, PM_NOREMOVE) != 0)
		{
			if (GetMessage(&msg, NULL, 0, 0) > 0)
			{
				TranslateMessage(&msg);
				DispatchMessage(&msg);
			}
		}
	}
}

namespace PixelWindowNS
{
	LRESULT CALLBACK WndProc(HWND hwnd, UINT message, WPARAM wParam, LPARAM lParam)
	{
		PixelWindow* winptr = (PixelWindow*)GetWindowLongPtr(hwnd, GWLP_USERDATA);
		if (winptr == NULL)
		{
			return DefWindowProc(hwnd, message, wParam, lParam);
		}
		else
		{
			return winptr->MyMsgProc(hwnd, message, wParam, lParam);
		}
	}

	LRESULT CALLBACK PixelWindow::MyMsgProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam)
	{
		switch (msg)
		{
		case WM_MOUSEWHEEL:
			if (wParam & 0x80000000)
				ReduceWindowSize();
			else
				IncreaseWindowSize();
			break;
		case WM_KEYDOWN:
			switch (wParam)
			{
			case VK_ADD:
			case VK_OEM_PLUS:
				IncreaseWindowSize();
				break;
			case VK_SUBTRACT:
			case VK_OEM_MINUS:
				ReduceWindowSize();
				break;
			case VK_ESCAPE:
				g_app_done = true;
				break;
			default:
				break;
			}
			break;
		case WM_CLOSE:
			g_app_done = true;
			DestroyWindow(hwnd);
			break;
		case WM_DESTROY:
			g_app_done = true;
			PostQuitMessage(0);
			break;
		default:
			return DefWindowProc(hwnd, msg, wParam, lParam);
		}
		return 0;
	}

	DWORD PixelWindow::run(LPVOID arg)
	{
		PixelWindow* pDXW = (PixelWindow*)arg;
		while (!g_app_done)
		{
			EnterCriticalSection(&pDXW->PixelBuffer.mutex);
			while (pDXW->PixelBuffer.size == 0)
			{
				SleepConditionVariableCS(&pDXW->PixelBuffer.cond_not_empty, &pDXW->PixelBuffer.mutex, INFINITE);
			}
			pDXW->CopyFrameToFrontBuffer();
			pDXW->PutFrontBufferOntoScreen();
			LeaveCriticalSection(&pDXW->PixelBuffer.mutex);
			WakeConditionVariable(&pDXW->PixelBuffer.cond_not_full);
		}
		return 0;
	}

	void PixelWindow::SetRGBAtXY(int r, int g, int b, int x, int y)
	{
		if ((x < 0) || (y < 0) || (x >= buffer_width) || (y >= buffer_height)) return;
		DWORD theColor = (0x00FF0000 & (r << 16)) | (0x0000FF00 & (g << 8)) | (0x000000FF & b);
		int offset = (x + buffer_width * y);
		PixelBuffer.Pixels[PixelBuffer.head][offset] = theColor;
	}

	PixelWindow::~PixelWindow()
	{
		mycond = true;

		// thrd->join();

		if (g_D3D)
		{
			g_D3D->Release();
			g_D3D = NULL;
		}

		DestroyWindow(p_window);

		MSG msg;
		while (PeekMessage(&msg, NULL, 0, 0, PM_REMOVE))
		{
			TranslateMessage(&msg);
			DispatchMessage(&msg);
		}

		UnregisterClass((LPCWSTR)&class_name, instance);

		for (int q = 0; q < PIXELWINDOW_BUFFER_SIZE; q++)
		{
			delete PixelBuffer.Pixels[q];
		}

		DeleteCriticalSection(&PixelBuffer.mutex);
	}

	PixelWindow::PixelWindow(wchar_t * theWindowTitle, int ww, int hh, bool WaitVSync)
	{
		InitializeCriticalSection(&PixelBuffer.mutex);
		InitializeConditionVariable(&PixelBuffer.cond_not_empty);
		InitializeConditionVariable(&PixelBuffer.cond_not_full);

		PixelBuffer.head = PixelBuffer.tail = PixelBuffer.size = 0;

		mycond = false;
		is_app_fullscreen = false;
		style = 0;
		p_window = 0;
		g_D3D = NULL;
		format = D3DFMT_X8R8G8B8;
		p_device = NULL;
		FrontBuffer = NULL;
		BackBuffer = NULL;

		nCmdShow = SW_SHOW;
		instance = GetModuleHandle(NULL);

		wsprintf(class_name, L"CN%s", theWindowTitle);
	
		window_class.style = CS_OWNDC;
		window_class.cbClsExtra = 0;
		window_class.cbWndExtra = 0;
		window_class.hInstance = instance;
		window_class.hIcon = LoadIcon(NULL, IDI_APPLICATION);
		window_class.hCursor = LoadCursor(NULL, IDC_ARROW);
		window_class.hbrBackground = (HBRUSH)GetStockObject(BLACK_BRUSH);
		window_class.lpszMenuName = NULL;
		window_class.lpszClassName = class_name;
		window_class.lpfnWndProc = (WNDPROC)WndProc;

		if (!RegisterClass(&window_class)) {
			MessageBox(p_window, L"Error during registration", L"Error", MB_OK);
			return;
		}

		buffer_width = ww;  buffer_height = hh;
		buffer_ratio = (float)buffer_width / (float)buffer_height;

		for (int q = 0; q < PIXELWINDOW_BUFFER_SIZE; q++)
		{
			PixelBuffer.Pixels[q] = new DWORD[buffer_width * buffer_height];
		}

		if (is_app_fullscreen) {
			window_width = GetSystemMetrics(SM_CXSCREEN);
			window_height = GetSystemMetrics(SM_CYSCREEN);
			style = WS_POPUP | WS_VISIBLE;
		}
		else {
			window_width = ww;
			window_height = hh;
			style = WS_OVERLAPPED | WS_SYSMENU | WS_VISIBLE | WS_MINIMIZEBOX;
		}
		p_window = CreateWindow(class_name,   //name of our registered class
			theWindowTitle,  //Window name/title
			style,        //Style flags
			0,            //X position
			0,            //Y position
			window_width, //width of window
			window_height,//height of window
			NULL,         //Parent window
			NULL,         //Menu
			instance,     //application instance handle
			NULL);        //pointer to window-creation data

		if (!p_window) {
			MessageBox(NULL, L"Cannot create window", NULL, MB_OK);
			exit(-1);
		}

		SetWindowLongPtr(p_window, GWLP_USERDATA, (long long)this);

		if (is_app_fullscreen) {
			ShowCursor(FALSE);
		}

		g_D3D = Direct3DCreate9(D3D_SDK_VERSION);
		if (!g_D3D) {
			MessageBox(p_window, L"Error creating Direct3DCreate9", L"Error", MB_OK);
			exit(-1);
		}

		ZeroMemory(&pp, sizeof(D3DPRESENT_PARAMETERS));
		pp.BackBufferCount = 1;  //We only need a single back buffer
		pp.BackBufferWidth = buffer_width;
		pp.BackBufferHeight = buffer_height;
		pp.MultiSampleType = D3DMULTISAMPLE_NONE; //No multi-sampling
		pp.MultiSampleQuality = 0;                //No multi-sampling
		pp.SwapEffect = D3DSWAPEFFECT_DISCARD;  // Throw away previous frames, we don't need them
		pp.hDeviceWindow = p_window;  //This is our main (and only) window
		pp.Windowed = TRUE;
		pp.Flags = D3DPRESENTFLAG_LOCKABLE_BACKBUFFER;            //No flags to set
		pp.FullScreen_RefreshRateInHz = D3DPRESENT_RATE_DEFAULT; //Default Refresh Rate
		pp.PresentationInterval =
			WaitVSync ? D3DPRESENT_INTERVAL_DEFAULT : D3DPRESENT_INTERVAL_IMMEDIATE;
		pp.BackBufferFormat = format;      //Display format
		pp.EnableAutoDepthStencil = FALSE; //No depth/stencil buffer

		hr = g_D3D->CreateDevice(D3DADAPTER_DEFAULT,
			D3DDEVTYPE_HAL,
			p_window,
			D3DCREATE_HARDWARE_VERTEXPROCESSING, // D3DCREATE_SOFTWARE_VERTEXPROCESSING,
			&pp,
			&p_device);
		if (FAILED(hr)) {
			MessageBox(p_window, L"Error creating Direct3D device with hardware acceleration", L"Error", MB_OK);
			exit(-1);
		}

		Result = p_device->CreateOffscreenPlainSurface(
			buffer_width,
			buffer_height,
			format,
			D3DPOOL_SYSTEMMEM,
			&FrontBuffer,
			NULL
		);

		if (FAILED(Result))
		{
			MessageBox(p_window, L"FrontBuffer failed.", L"Error", MB_OK);
			exit(-1);
		}

		p_device->GetBackBuffer(0, 0, D3DBACKBUFFER_TYPE_MONO, &BackBuffer);

		ResizeWindow(window_width, window_height);

		ShowWindow(p_window, nCmdShow);
		UpdateWindow(p_window);

		CreateThread(NULL, NULL, (LPTHREAD_START_ROUTINE)&PixelWindow::run, this, NULL, NULL);
	}

	void PixelWindow::CalculateExactRect(int cx, int cy, RECT& rect)
	{
		RECT rcWindow;
		GetWindowRect(p_window, &rcWindow);
		RECT rcClient;
		GetClientRect(p_window, &rcClient);
		cx += (rcWindow.right - rcWindow.left) - rcClient.right;
		cy += (rcWindow.bottom - rcWindow.top) - rcClient.bottom;

		rect.left = rcWindow.left;
		rect.top = rcWindow.top;
		rect.right = rect.left + cx;
		rect.bottom = rect.top + cy;
	}

	void PixelWindow::ReduceWindowSize()
	{
		if (window_width <= 150) return;
		window_width -= 20;
		window_height = (int)((float)window_width / (float)buffer_ratio);
		ResizeWindow(window_width, window_height);
	}

	void PixelWindow::IncreaseWindowSize()
	{
		if (window_width >= 1000) return;
		window_width += 20;
		window_height = (int)((float)window_width / (float)buffer_ratio);
		ResizeWindow(window_width, window_height);
	}

	void PixelWindow::ResizeWindow(int width, int height)
	{
		RECT rect = { 0 };
		CalculateExactRect(width, height, rect);
		MoveWindow(p_window,
			rect.left, rect.top,
			rect.right - rect.left,
			rect.bottom - rect.top,
			TRUE);
	}

	void PixelWindow::DisplayBuffer(unsigned char* bfr)
	{
		EnterCriticalSection(&PixelBuffer.mutex);
		while (PixelBuffer.size >= PIXELWINDOW_BUFFER_SIZE)
		{
			SleepConditionVariableCS(&PixelBuffer.cond_not_full, &PixelBuffer.mutex, INFINITE);
		}

		memcpy(PixelBuffer.Pixels[PixelBuffer.head], bfr, buffer_width * buffer_height * 4);

		PixelBuffer.head++;
		PixelBuffer.head %= PIXELWINDOW_BUFFER_SIZE;
		PixelBuffer.size++;
		LeaveCriticalSection(&PixelBuffer.mutex);
		WakeConditionVariable(&PixelBuffer.cond_not_empty);
	}

	void PixelWindow::CopyFrameToFrontBuffer()
	{
		FrontBuffer->LockRect(&lr, NULL, 0);
		CopyMemory(lr.pBits, PixelBuffer.Pixels[PixelBuffer.tail], buffer_width * buffer_height * 4);
		PixelBuffer.tail++;
		PixelBuffer.tail %= PIXELWINDOW_BUFFER_SIZE;
		PixelBuffer.size--;
		FrontBuffer->UnlockRect();
	}

	void PixelWindow::PutFrontBufferOntoScreen()
	{
		hr = p_device->Clear(0,  //Number of rectangles to clear, we're clearing everything so set it to 0
			NULL, //Pointer to the rectangles to clear, NULL to clear whole display
			D3DCLEAR_TARGET,   //What to clear.  We don't have a Z Buffer or Stencil Buffer
			0x00000000, //Colour to clear to (AARRGGBB)
			1.0f,  //Value to clear ZBuffer to, doesn't matter since we don't have one
			0);   //Stencil clear value, again, we don't have one, this value doesn't matter
		if (FAILED(hr)) {
			return;
		}

		hr = p_device->BeginScene();
		if (FAILED(hr)) {
			return;
		}

		p_device->UpdateSurface(FrontBuffer, NULL, BackBuffer, NULL);

		p_device->EndScene();

		hr = p_device->Present(NULL,  //Source rectangle to display, NULL for all of it
			NULL,  //Destination rectangle, NULL to fill whole display
			NULL,  //Target window, if NULL uses device window set in CreateDevice
			NULL);//Dirty Region, set it to NULL

		if (FAILED(hr)) {
			return;
		}

	}
}
