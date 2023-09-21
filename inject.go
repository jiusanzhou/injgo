// Package injgo is a package for injecting in golang.
package injgo

import (
	"syscall"
	"unsafe"
)

func ptr(val interface{}) uintptr {
	switch val.(type) {
	case byte:
		return uintptr(val.(byte))
	case bool:
		isTrue := val.(bool)
		if isTrue {
			return uintptr(1)
		}
		return uintptr(0)
	case string:
		bytePtr, _ := syscall.BytePtrFromString(val.(string))
		return uintptr(unsafe.Pointer(bytePtr))
	case int:
		return uintptr(val.(int))
	case uint:
		return uintptr(val.(uint))
	case uintptr:
		return val.(uintptr)
	default:
		return uintptr(0)
	}
}
