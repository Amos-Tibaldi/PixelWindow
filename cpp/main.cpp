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

#include <Windows.h>
#include <d3d9.h>
#include "PixelWindow.h"

using namespace std;
using namespace PixelWindowNS;

bool g_app_done = false;
unsigned char * redbuffer = 0, * greenbuffer = 0, * bluebuffer = 0;

int main()
{
    redbuffer = (unsigned char*)malloc(640 * 480 * 4);
    greenbuffer = (unsigned char*)malloc(640 * 480 * 4);
    bluebuffer = (unsigned char*)malloc(640 * 480 * 4);

    wchar_t redname[100] = { L"RED" };
    wchar_t greenname[100] = { L"GREEN" };
    wchar_t bluename[100] = { L"BLUE" };

    PixelWindow* pwred = new PixelWindow(redname, 640, 480, true);
	LDAPIXELWINDOWHANDLE hred = (LDAPIXELWINDOWHANDLE)pwred;
    PixelWindow* pwgreen = new PixelWindow(greenname, 640, 480, true);
    LDAPIXELWINDOWHANDLE hgreen = (LDAPIXELWINDOWHANDLE)pwgreen;
    PixelWindow* pwblue = new PixelWindow(bluename, 640, 480, true);
    LDAPIXELWINDOWHANDLE hblue = (LDAPIXELWINDOWHANDLE)pwblue;

    unsigned char level = 100;

	while (!g_app_done)
	{
        for (int i = 0; i < 640; i++)
        {
            for (int j = 0; j < 480; j++)
            {
                unsigned long offset = 4 * (j * 640 + i);

                redbuffer[offset + 0] = 0;
                redbuffer[offset + 1] = 0;
                redbuffer[offset + 2] = level;
                redbuffer[offset + 3] = 0;

                greenbuffer[offset + 0] = 0;
                greenbuffer[offset + 1] = level;
                greenbuffer[offset + 2] = 0;
                greenbuffer[offset + 3] = 0;

                bluebuffer[offset + 0] = level;
                bluebuffer[offset + 1] = 0;
                bluebuffer[offset + 2] = 0;
                bluebuffer[offset + 2] = 0;
            }
        }
        level++;
        level %= 255;
        for (int i = 0; i < 100; i++)
        {
            unsigned long offset = 4 * (i * 640 + i);
            redbuffer[offset + 0] = greenbuffer[offset + 0] = bluebuffer[offset + 0] = 255;
            redbuffer[offset + 1] = greenbuffer[offset + 1] = bluebuffer[offset + 1] = 255;
            redbuffer[offset + 2] = greenbuffer[offset + 2] = bluebuffer[offset + 2] = 255;
            redbuffer[offset + 3] = greenbuffer[offset + 3] = bluebuffer[offset + 3] = 255;
            offset = 4 * ((480 - i - 1) * 640 + i);
            redbuffer[offset + 0] = greenbuffer[offset + 0] = bluebuffer[offset + 0] = 255;
            redbuffer[offset + 1] = greenbuffer[offset + 1] = bluebuffer[offset + 1] = 255;
            redbuffer[offset + 2] = greenbuffer[offset + 2] = bluebuffer[offset + 2] = 255;
            redbuffer[offset + 3] = greenbuffer[offset + 3] = bluebuffer[offset + 3] = 255;
        }
        
        LDAPIXELWindowDisplayBuffer(hred, redbuffer);
        LDAPIXELWindowDisplayBuffer(hgreen, greenbuffer);
        LDAPIXELWindowDisplayBuffer(hblue, bluebuffer);

		LDACallWhenIdle();
		SwitchToThread();
		Sleep(1);
	}

    free((void*)redbuffer);
    free((void*)greenbuffer);
    free((void*)bluebuffer);
	
	return 0;
}
