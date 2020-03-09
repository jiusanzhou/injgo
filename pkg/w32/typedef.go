package w32

// +build windows

type (
	HANDLE uintptr
	HKEY   HANDLE
	BOOL   int32
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms684225.aspx
type MODULEENTRY32 struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  *uint8
	ModBaseSize  uint32
	HModule      HMODULE
	SzModule     [MAX_MODULE_NAME32 + 1]uint16
	SzExePath    [MAX_PATH]uint16
}

// https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/ns-tlhelp32-tagprocessentry32
type PROCESSENTRY32 struct {
	Size            uint32
	CntUsage        uint32
	ProcessID       uint32
	DefaultHeapID   ULONG_PTR
	ModuleID        uint32
	Threads         uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [MAX_PATH]uint16
}

// VS_FIXEDFILEINFO presents file information
type VS_FIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

// WinVersion window version struct
type WinVersion struct {
	Major uint32
	Minor uint32
	Patch uint32
	Build uint32
}
