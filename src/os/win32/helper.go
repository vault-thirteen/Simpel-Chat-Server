package win32

import (
	"errors"
	"fmt"
	"runtime"

	"golang.org/x/sys/windows"
)

// MustBeNoCallError checks for a call error.
// If call error is set, it finds a name of the function which started this
// function and panics about that function's failure.
func MustBeNoCallError(callErr windows.Errno) {
	if !errors.Is(callErr, windows.ERROR_SUCCESS) {
		caller, _, _, _ := runtime.Caller(1) // Skip 1 function up in stack: mustBeNoCallError -> caller.
		fn := runtime.FuncForPC(caller)
		panic(fmt.Sprintf("%s syscall failed: %v", fn.Name(), callErr.Error()))
	}
}

func (dc *DllController) EnableConsoleColours() (err error) {
	var hConsole windows.Handle
	hConsole = dc.GetStdHandle(STD_OUTPUT_HANDLE)
	if hConsole == INVALID_HANDLE_VALUE {
		return errors.New("GetStdHandle returned invalid handle")
	}

	var dwMode DWORD
	if !dc.GetConsoleMode(hConsole, &dwMode) {
		return errors.New("GetConsoleMode failed")
	}

	dwMode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if !dc.SetConsoleMode(hConsole, dwMode) {
		return errors.New("SetConsoleMode failed")
	}

	return nil
}
