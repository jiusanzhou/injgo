package injgo

import (
	"errors"
	"syscall"
	"unsafe"

	"go.zoe.im/injgo/pkg/w32"
)

// Process ...
type Process struct {
	ProcessID int
	Name      string
	ExePath   string
}

var (
	// ErrProcessNotFound ...
	ErrProcessNotFound = errors.New("process not found")
	// ErrCreateSnapshot ...
	ErrCreateSnapshot = errors.New("create snapshot error")
)

// FindProcessByName get process information by name
func FindProcessByName(name string) (*Process, error) {
	handle, _ := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if handle == 0 {
		return nil, ErrCreateSnapshot
	}
	defer syscall.CloseHandle(handle)

	var entry = syscall.ProcessEntry32{}
	entry.Size = uint32(unsafe.Sizeof(entry))
	var process Process

	for true {
		if nil != syscall.Process32Next(handle, &entry) {
			break
		}

		_exeFile := w32.UTF16PtrToString(&entry.ExeFile[0])
		if name == _exeFile {
			process.Name = _exeFile
			process.ProcessID = int(entry.ProcessID)
			// TODO: 找到路径
			process.ExePath = _exeFile
			return &process, nil
		}

	}
	return nil, ErrProcessNotFound
}

// CreateProcess create a new process
//
// exePath execute path
func CreateProcess(exePath string) (*Process, error) {
	// TODO:
	var sI syscall.StartupInfo
	var pI syscall.ProcessInformation

	argv := syscall.StringToUTF16Ptr(exePath)

	err := syscall.CreateProcess(
		nil, argv, nil,
		nil, true, 0,
		nil, nil, &sI, &pI,
	)
	if err != nil {
		return nil, err
	}
	return &Process{
		ProcessID: int(pI.ProcessId),
		// TODO:
	}, nil
}
