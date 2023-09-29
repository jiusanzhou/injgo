//go:build windows
// +build windows

package w32

import (
	"encoding/binary"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCloseHandle              = modkernel32.NewProc("CloseHandle")
	procCreateProcessA           = modkernel32.NewProc("CreateProcessA")
	procCreateProcessW           = modkernel32.NewProc("CreateProcessW")
	procGetModuleHandleA         = modkernel32.NewProc("GetModuleHandleA")
	procLoadLibraryA             = modkernel32.NewProc("LoadLibraryA")
	ProcFreeLibrary              = modkernel32.NewProc("FreeLibrary")
	procCreateRemoteThread       = modkernel32.NewProc("CreateRemoteThread")
	procCreateToolhelp32Snapshot = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procTerminateProcess         = modkernel32.NewProc("TerminateProcess")
	procOpenProcess              = modkernel32.NewProc("OpenProcess")
	procVirtualAlloc             = modkernel32.NewProc("VirtualAlloc")
	procVirtualAllocEx           = modkernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx            = modkernel32.NewProc("VirtualFreeEx")
	procReadProcessMemory        = modkernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory       = modkernel32.NewProc("WriteProcessMemory")
	procProcess32First           = modkernel32.NewProc("Process32FirstW")
	procProcess32Next            = modkernel32.NewProc("Process32NextW")
	procModule32First            = modkernel32.NewProc("Module32FirstW")
	procModule32Next             = modkernel32.NewProc("Module32NextW")
	procGetProcAddress           = modkernel32.NewProc("GetProcAddress")
	procWaitForSingleObj         = modkernel32.NewProc("WaitForSingleObject")
)

func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uintptr) (handle HANDLE, err error) {
	inherit := 0
	if inheritHandle {
		inherit = 1
	}

	ret, _, err := procOpenProcess.Call(
		uintptr(desiredAccess),
		uintptr(inherit),
		processId)
	if err != nil && IsErrSuccess(err) {
		err = nil
	}
	handle = HANDLE(ret)
	return
}

func TerminateProcess(hProcess HANDLE, uExitCode uint) bool {
	ret, _, _ := procTerminateProcess.Call(
		uintptr(hProcess),
		uintptr(uExitCode))
	return ret != 0
}

func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}

func CreateToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

func Process32First(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	if pe.Size == 0 {
		pe.Size = uint32(unsafe.Sizeof(*pe))
	}
	ret, _, _ := procProcess32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)))

	return ret != 0
}

func Process32Next(snapshot HANDLE, pe *PROCESSENTRY32) bool {
	if pe.Size == 0 {
		pe.Size = uint32(unsafe.Sizeof(*pe))
	}
	ret, _, _ := procProcess32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(pe)))

	return ret != 0
}

func Module32First(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func GetModuleHandleA(modulename string) uintptr {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		bytes := []byte(modulename)
		mn = uintptr(unsafe.Pointer(&bytes[0]))
	}
	ret, _, _ := procGetModuleHandleA.Call(mn)
	return ret
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682425(v=vs.85).aspx
func CreateProcessA(lpApplicationName *string,
	lpCommandLine string,
	lpProcessAttributes *syscall.SecurityAttributes,
	lpThreadAttributes *syscall.SecurityAttributes,
	bInheritHandles bool,
	dwCreationFlags uint32,
	lpEnvironment *string,
	lpCurrentDirectory *uint16,
	lpStartupInfo *syscall.StartupInfo,
	lpProcessInformation *syscall.ProcessInformation) {

	inherit := 0
	if bInheritHandles {
		inherit = 1
	}

	procCreateProcessA.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(*lpApplicationName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpCommandLine))),
		uintptr(unsafe.Pointer(lpProcessAttributes)),
		uintptr(unsafe.Pointer(lpThreadAttributes)),
		uintptr(inherit),
		uintptr(dwCreationFlags),
		uintptr(unsafe.Pointer(lpEnvironment)),
		uintptr(unsafe.Pointer(lpCurrentDirectory)),
		uintptr(unsafe.Pointer(lpStartupInfo)),
		uintptr(unsafe.Pointer(lpProcessInformation)))
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366890(v=vs.85).aspx
func VirtualAllocEx(hProcess HANDLE, lpAddress uintptr, dwSize uintptr, flAllocationType uintptr, flProtect uintptr) (addr uintptr, err error) {
	ret, _, err := procVirtualAllocEx.Call(
		uintptr(hProcess), // The handle to a process.
		lpAddress,         // The pointer that specifies a desired starting address for the region of pages that you want to allocate.
		dwSize,            // The size of the region of memory to allocate, in bytes.
		flAllocationType,
		flProtect)
	if int(ret) == 0 {
		return ret, err
	}
	return ret, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366887(v=vs.85).aspx
func VirtualAlloc(lpAddress int, dwSize int, flAllocationType int, flProtect int) (addr uintptr, err error) {
	ret, _, err := procVirtualAlloc.Call(
		uintptr(lpAddress), // The starting address of the region to allocate
		uintptr(dwSize),    // The size of the region of memory to allocate, in bytes.
		uintptr(flAllocationType),
		uintptr(flProtect))
	if int(ret) == 0 {
		return ret, err
	}
	return ret, nil
}

// https://github.com/AllenDang/w32/pull/62/commits/08a52ff508c6b2b9b9bf5ee476109b903dcf219d
func VirtualFreeEx(hProcess HANDLE, lpAddress, dwSize uintptr, dwFreeType uint32) (uintptr, error) {
	ret, _, err := procVirtualFreeEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(dwFreeType),
	)
	if IsErrSuccess(err) {
		return ret, nil
	}
	return ret, err
}

func GetProcAddress(h uintptr, name string) (uintptr, error) {
	return syscall.GetProcAddress(syscall.Handle(h), name)
}

func LoadLibraryAddress(libraryPtr uintptr) (uintptr, error) {
	loadLibraryAddress, _, err := procGetProcAddress.Call(modkernel32.Handle(), libraryPtr)
	if !IsErrSuccess(err) {
		return loadLibraryAddress, err
	}
	return loadLibraryAddress, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682437(v=vs.85).aspx
// Credit: https://github.com/contester/runlib/blob/master/win32/win32_windows.go#L577
func CreateRemoteThread(hprocess HANDLE, sa *syscall.SecurityAttributes,
	stackSize uintptr, startAddress uintptr, parameter uintptr, creationFlags uintptr) (HANDLE, int, error) {
	var threadId int
	r1, _, e1 := procCreateRemoteThread.Call(
		uintptr(hprocess),
		uintptr(unsafe.Pointer(sa)),
		stackSize,
		startAddress,
		parameter,
		creationFlags,
		uintptr(unsafe.Pointer(&threadId)))

	if int(r1) == 0 {
		return INVALID_HANDLE, 0, e1
	}
	return HANDLE(r1), threadId, nil
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/synchapi/nf-synchapi-waitforsingleobject
func WaitForSingleObj(hprocess HANDLE, milliseconds int) error {
	_, _, err := procWaitForSingleObj.Call(uintptr(hprocess), uintptr(milliseconds))
	if IsErrSuccess(err) {
		return nil
	}
	return err
}

// Writes data to an area of memory in a specified process. The entire area to be written to must be accessible or the operation fails.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms681674(v=vs.85).aspx
func WriteProcessMemory(hProcess HANDLE, lpBaseAddress uintptr, data uintptr, size uintptr) (err error) {
	var numBytesRead uintptr

	_, _, err = procWriteProcessMemory.Call(uintptr(hProcess),
		lpBaseAddress,
		data,
		size,
		uintptr(unsafe.Pointer(&numBytesRead)))
	if !IsErrSuccess(err) {
		return
	}
	err = nil
	return
}

// Write process memory with a source of uint32
func WriteProcessMemoryAsUint32(hProcess HANDLE, lpBaseAddress uintptr, data uint32) (err error) {
	bData := make([]byte, 4)
	binary.LittleEndian.PutUint32(bData, data)
	err = WriteProcessMemory(hProcess, lpBaseAddress, uintptr(unsafe.Pointer(&bData[0])), 4)
	if err != nil {
		return
	}
	return
}

// Reads data from an area of memory in a specified process. The entire area to be read must be accessible or the operation fails.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms680553(v=vs.85).aspx
func ReadProcessMemory(hProcess HANDLE, lpBaseAddress uintptr, size uintptr) (data []byte, err error) {
	var numBytesRead uintptr
	data = make([]byte, size)

	_, _, err = procReadProcessMemory.Call(uintptr(hProcess),
		lpBaseAddress,
		uintptr(unsafe.Pointer(&data[0])),
		size,
		uintptr(unsafe.Pointer(&numBytesRead)))
	if !IsErrSuccess(err) {
		return
	}
	err = nil
	return
}

// Read process memory and convert the returned data to uint32
func ReadProcessMemoryAsUint32(hProcess HANDLE, lpBaseAddress uintptr) (buffer uint32, err error) {
	data, err := ReadProcessMemory(hProcess, lpBaseAddress, 4)
	if err != nil {
		return
	}
	buffer = binary.LittleEndian.Uint32(data)
	return
}
