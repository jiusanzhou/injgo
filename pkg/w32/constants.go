package w32

// +build windows

const (
	NO_ERROR      = 0
	ERROR_SUCCESS = 0
)

// CreateToolhelp32Snapshot flags
const (
	TH32CS_SNAPHEAPLIST = 0x00000001
	TH32CS_SNAPPROCESS  = 0x00000002
	TH32CS_SNAPTHREAD   = 0x00000004
	TH32CS_SNAPMODULE   = 0x00000008
	TH32CS_SNAPMODULE32 = 0x00000010
	TH32CS_INHERIT      = 0x80000000
	TH32CS_SNAPALL      = TH32CS_SNAPHEAPLIST | TH32CS_SNAPMODULE | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD
)

const (
	MAX_MODULE_NAME32 = 255
	MAX_PATH          = 260
)

const (
	MEM_COMMIT     = 0x00001000
	MEM_RESERVE    = 0x00002000
	MEM_RESET      = 0x00080000
	MEM_RESET_UNDO = 0x1000000

	MEM_LARGE_PAGES = 0x20000000
	MEM_PHYSICAL    = 0x00400000
	MEM_TOP_DOWN    = 0x00100000

	MEM_DECOMMIT = 0x4000
	MEM_RELEASE  = 0x8000
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366786(v=vs.85).aspx
const (
	PAGE_EXECUTE           = 0x10
	PAGE_EXECUTE_READ      = 0x20
	PAGE_EXECUTE_READWRITE = 0x40
	PAGE_EXECUTE_WRITECOPY = 0x80
	PAGE_NOACCESS          = 0x01
	PAGE_READWRITE         = 0x04
	PAGE_WRITECOPY         = 0x08
	PAGE_TARGETS_INVALID   = 0x40000000
	PAGE_TARGETS_NO_UPDATE = 0x40000000
)

//Process Access Rights
//https://msdn.microsoft.com/en-us/library/windows/desktop/ms684880(v=vs.85).aspx
const (
	PROCESS_CREATE_PROCESS            = 0x0080  //Required to create a process.
	PROCESS_CREATE_THREAD             = 0x0002  //Required to create a thread.
	PROCESS_DUP_HANDLE                = 0x0040  //Required to duplicate a handle using DuplicateHandle.
	PROCESS_QUERY_INFORMATION         = 0x0400  //Required to retrieve certain information about a process, such as its token, exit code, and priority class (see OpenProcessToken).
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000  //Required to retrieve certain information about a process (see GetExitCodeProcess, GetPriorityClass, IsProcessInJob, QueryFullProcessImageName). A handle that has the PROCESS_QUERY_INFORMATION access right is automatically granted
	PROCESS_SET_INFORMATION           = 0x0200  //Required to set certain information about a process, such as its priority class (see SetPriorityClass).
	PROCESS_SET_QUOTA                 = 0x0100  //Required to set memory limits using SetProcessWorkingSetSize.
	PROCESS_SUSPEND_RESUME            = 0x0800  //Required to suspend or resume a process.
	PROCESS_TERMINATE                 = 0x0001  //Required to terminate a process using TerminateProcess.
	PROCESS_VM_OPERATION              = 0x0008  //Required to perform an operation on the address space of a process (see VirtualProtectEx and WriteProcessMemory).
	PROCESS_VM_READ                   = 0x0010  //Required to read memory in a process using ReadProcessMemory.
	PROCESS_VM_WRITE                  = 0x0020  //Required to write to memory in a process using WriteProcessMemory.
	PROCESS_ALL_ACCESS                = 2035711 //This is not recommended.
	SYNCHRONIZE                       = 0x00100000
)