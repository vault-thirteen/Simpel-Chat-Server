package win32

import "golang.org/x/sys/windows"

// #define ENABLE_VIRTUAL_TERMINAL_PROCESSING  0x0004
const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 4

// #define INVALID_HANDLE_VALUE ((HANDLE)(LONG_PTR)-1)
const INVALID_HANDLE_VALUE = windows.InvalidHandle

// #define STD_OUTPUT_HANDLE ((DWORD)-11) // WinBase.h:820.
const STD_OUTPUT_HANDLE = 1<<32 - 11
