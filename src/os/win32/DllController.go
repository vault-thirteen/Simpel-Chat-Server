package win32

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/sys/windows"
)

const (
	OS_Windows = "windows"
)

const DllFuncNamePrefix = ""

const (
	Err_DllFileIsNotSet = "DLL file is not set"
)

const (
	MsgF_LoadingLibrary     = "Loading library: %v."
	Msg_LoadingFunctions    = "Loading functions: "
	MsgF_FunctionNameInList = "[%s] "
)

type DllController struct {
	guard              *sync.Mutex
	windowsApiIsLoaded *atomic.Bool

	// Pointers to DLL functions.
	DllFunc DllFunc

	// Handles and function mappings.
	kernel32 Dll
	user32   Dll
}

func NewDllController() (dc *DllController) {
	dc = &DllController{
		guard: new(sync.Mutex),
	}

	dc.windowsApiIsLoaded = new(atomic.Bool)
	dc.windowsApiIsLoaded.Store(false)

	return dc
}

func (dc *DllController) Init() (err error) {
	dc.guard.Lock()
	defer dc.guard.Unlock()

	os := runtime.GOOS
	log.Println("Operating system:", os)

	switch strings.ToLower(os) {
	case OS_Windows:
		return dc.loadWindowsApi()
	default:
		log.Println("Current operating system is not fully supported. Some features may be unavailable.")
	}

	return nil
}

func (dc *DllController) loadWindowsApi() (err error) {
	if dc.windowsApiIsLoaded.Load() {
		return nil
	}

	err = dc.prepareObjects()
	if err != nil {
		return err
	}

	err = dc.loadWindowsLibraries()
	if err != nil {
		return err
	}

	dc.windowsApiIsLoaded.Store(true)

	return nil
}

func (dc *DllController) prepareObjects() (err error) {
	dc.kernel32 = NewDll(
		[]FuncMapping{
			{&dc.DllFunc.GetLastError, "GetLastError"},
			{&dc.DllFunc.SetLastError, "SetLastError"},
			{&dc.DllFunc.GetStdHandle, "GetStdHandle"},
			{&dc.DllFunc.GetConsoleMode, "GetConsoleMode"},
			{&dc.DllFunc.SetConsoleMode, "SetConsoleMode"},
		})

	dc.user32 = NewDll([]FuncMapping{})

	return nil
}

func (dc *DllController) loadWindowsLibraries() (err error) {
	err = dc.loadLibrary(DllFile_Kernel32, &dc.kernel32.handle, dc.kernel32.funcMapping, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	err = dc.loadLibrary(DllFile_User32, &dc.user32.handle, dc.user32.funcMapping, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DllController) loadLibrary(
	dllFilePath string,
	h *windows.Handle,
	funcMappings []FuncMapping,
	funcNamePrefix string,
) (err error) {
	if len(dllFilePath) == 0 {
		return errors.New(Err_DllFileIsNotSet)
	}

	fmt.Println(fmt.Sprintf(MsgF_LoadingLibrary, dllFilePath))
	*h, err = windows.LoadLibrary(dllFilePath)
	if err != nil {
		return err
	}

	fmt.Print(Msg_LoadingFunctions)
	for _, fm := range funcMappings {
		fmt.Printf(MsgF_FunctionNameInList, fm.FunctionName)
		*(fm.Fn), err = windows.GetProcAddress(*h, funcNamePrefix+fm.FunctionName)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}
