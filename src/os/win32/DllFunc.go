package win32

// DllFunc object is used for fast access to functions from DLL files.
type DllFunc struct {
	// kernel32.dll
	GetLastError   uintptr // https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-getlasterror
	SetLastError   uintptr // https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-setlasterror
	GetStdHandle   uintptr // https://learn.microsoft.com/en-us/windows/console/getstdhandle
	GetConsoleMode uintptr // https://learn.microsoft.com/en-us/windows/console/getconsolemode
	SetConsoleMode uintptr // https://learn.microsoft.com/en-us/windows/console/setconsolemode

	// user32.dll
}
