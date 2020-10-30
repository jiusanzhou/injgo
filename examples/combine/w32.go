package main

import (
	"errors"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")

	procEnumWindows              = user32.MustFindProc("EnumWindows")
	procIsWindow                 = user32.MustFindProc("IsWindow")
	procFindWindow               = user32.MustFindProc("FindWindowW")
	procGetWindowTextW           = user32.MustFindProc("GetWindowTextW")
	procSetWindowsHookEx         = user32.MustFindProc("SetWindowsHookExW")
	procSetWinEventHook          = user32.MustFindProc("SetWinEventHook")
	procUnhookWinEvent           = user32.MustFindProc("UnhookWinEvent")
	procGetWindowThreadProcessID = user32.MustFindProc("GetWindowThreadProcessId")
	procSetFocus                 = user32.MustFindProc("SetFocus")
	procSetForegroundWindow      = user32.MustFindProc("SetForegroundWindow")
	procEnableWindow             = user32.MustFindProc("EnableWindow")
	procShowWindow               = user32.MustFindProc("ShowWindow")
	procMoveWindow               = user32.MustFindProc("MoveWindow")
)

// ...
const (
	WINEVENT_OUTOFCONTEXT = 0x0000
	WM_SYSCOMMAND         = 0x0112
	WM_SETREDRAW          = 11
	SC_MOVE               = 61456
	HTCAPTION             = 2

	EVENT_MIN                  = 0x00000001
	EVENT_MAX                  = 0x7FFFFFFF
	EVENT_OBJECT_IME_CHANGE    = 0x8029
	EVENT_SYSTEM_MOVESIZESTART = 0x000A
	EVENT_SYSTEM_SCROLLINGEND  = 0x0013
	EVENT_SYSTEM_FOREGROUND    = 0x0003
)

// BoolToBOOL ...
func BoolToBOOL(value bool) int {
	if value {
		return 1
	}

	return 0
}

// EnumWindows ...
func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// GetWindowText ...
func GetWindowText(hwnd syscall.Handle) (str string, err error) {
	b := make([]uint16, 200)
	maxCount := int32(len(b))
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(&b[0])), uintptr(maxCount))
	len := int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return
	}
	str = syscall.UTF16ToString(b)
	return
}

// GetThread ...
func GetThread(hwnd syscall.Handle) (syscall.Handle, uint) {
	var id uint
	ret, _, _ := procGetWindowThreadProcessID.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&id)))

	return syscall.Handle(ret), id
}

// SetForegroundWindow ...
func SetForegroundWindow(hwnd syscall.Handle) bool {
	ret, _, _ := procSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

// SetFocus ...
func SetFocus(hwnd syscall.Handle) syscall.Handle {
	ret, _, _ := procSetFocus.Call(uintptr(hwnd))
	return syscall.Handle(ret)
}

// EnableWindow ...
func EnableWindow(hwnd syscall.Handle, b bool) bool {
	ret, _, _ := procEnableWindow.Call(uintptr(hwnd), uintptr(BoolToBOOL(b)))
	return ret != 0
}

// ShowWindow ...
func ShowWindow(hwnd syscall.Handle, cmdshow bool) bool {
	ret, _, _ := procShowWindow.Call(uintptr(hwnd), uintptr(BoolToBOOL(cmdshow)))

	return ret != 0
}

// MoveWindow ...
func MoveWindow(hwnd syscall.Handle, x, y, width, height int, repaint bool) bool {
	ret, _, _ := procMoveWindow.Call(
		uintptr(hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(BoolToBOOL(repaint)))

	return ret != 0
}

// Window ...
type Window struct {
	hwnd      syscall.Handle
	title     string
	wevthooks map[interface{}]syscall.Handle
	thread    *Thread
}

// Thread ...
type Thread struct {
	hwnd syscall.Handle
	id   uint
}

// NewThread ...
func NewThread(h syscall.Handle, i uint) *Thread {
	return &Thread{
		h, i,
	}
}

// NewWindow ...
func NewWindow(hwnd syscall.Handle) *Window {
	title, _ := GetWindowText(hwnd)
	return &Window{
		hwnd:      hwnd,
		title:     title,
		thread:    NewThread(GetThread(hwnd)), // get the thread
		wevthooks: make(map[interface{}]syscall.Handle),
	}
}

// Move ...
func (w *Window) Move(x, y, width, height int, repaint bool) bool {
	return MoveWindow(w.hwnd, x, y, width, height, repaint)
}

// Show ...
func (w *Window) Show() bool {
	return ShowWindow(w.hwnd, true)
}

// Hidden ...
func (w *Window) Hidden() bool {
	return ShowWindow(w.hwnd, false)
}

// Enable ...
func (w *Window) Enable() bool {
	return EnableWindow(w.hwnd, true)
}

// Disable ...
func (w *Window) Disable() bool {
	return EnableWindow(w.hwnd, false)
}

// SetForeground ...
func (w *Window) SetForeground() bool {
	return SetForegroundWindow(w.hwnd)
}

// UnhookWinEvent ...
func (w *Window) UnhookWinEvent(fn interface{}) bool {
	fnkey := reflect.ValueOf(fn)
	v, ok := w.wevthooks[fnkey]
	if !ok {
		return ok
	}
	ret, _, _ := procUnhookWinEvent.Call(uintptr(v))
	return ret != 0
}

// Event ...
type Event struct {
	Hook     syscall.Handle
	HWND     syscall.Handle
	Type     int
	ObjectID int32
	ChildID  int32
	CreateAt uint32
}

// SetWinEventHook ... TODO: use options
func (w *Window) SetWinEventHook(fn func(evt *Event) error, evts ...int) (syscall.Handle, error) {

	fnkey := reflect.ValueOf(fn)
	if v, ok := w.wevthooks[fnkey]; ok {
		return v, nil
	}

	// create a new fn
	ofn := func(hook syscall.Handle, evt uint32, hwnd syscall.Handle, idObject int32, idChild int32, dwEventThread uint32, dwmsEventTime uint32) syscall.Handle {
		_ = fn(&Event{
			Hook:     hook,
			HWND:     hwnd,
			Type:     int(evt),
			ObjectID: idObject,
			ChildID:  idChild,
			CreateAt: dwEventThread,
		})
		return 1
	}

	evtMin := EVENT_MIN
	evtMax := EVENT_MAX

	if len(evts) == 1 {
		evtMin = evts[0]
		evtMax = evts[0]
	} else {
		evtMin = evts[0]
		evtMax = evts[1]
	}

	ret, _, err := procSetWinEventHook.Call(
		uintptr(evtMin), uintptr(evtMax),
		uintptr(0), // ??? dll without
		syscall.NewCallback(ofn),
		uintptr(w.thread.id), // ??? process id
		uintptr(0),           // ??? thread id
		uintptr(WINEVENT_OUTOFCONTEXT),
	)

	if !strings.Contains(err.Error(), "successfully") {
		return 0, err
	}

	// store
	w.wevthooks[fnkey] = syscall.Handle(ret)

	return syscall.Handle(ret), nil
}

// FindWindow ...
func FindWindow(title string) (*Window, error) {

	ret, _, _ := syscall.Syscall(procFindWindow.Addr(), 2,
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		0)

	hwnd := syscall.Handle(ret)
	if hwnd == 0 {
		return nil, errors.New("can't found")
	}

	return NewWindow(syscall.Handle(uintptr(hwnd))), nil
}

// ListWindows ...
func ListWindows() ([]*Window, error) {
	ws := []*Window{}
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		is, _, _ := procIsWindow.Call(uintptr(hwnd))
		if is == 0 {
			return 1
		}
		ws = append(ws, NewWindow(hwnd))
		return 1 // continue enumeration
	})
	return ws, EnumWindows(cb, 0)
}
