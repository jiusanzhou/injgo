package w32

// +build windows

import (
	"encoding/binary"
	"encoding/hex"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func LOWORD(dw uint32) uint16 {
	return uint16(dw)
}

func HIWORD(dw uint32) uint16 {
	return uint16(dw >> 16 & 0xffff)
}

func LOBYTE(word uint16) uint8 {
	return uint8(word)
}

func HIBYTE(word uint16) uint8 {
	return uint8(word >> 8 & 0xff)
}

func BoolToBOOL(value bool) BOOL {
	if value {
		return 1
	}

	return 0
}

func UTF16PtrToString(cstr *uint16) string {
	if cstr != nil {
		us := make([]uint16, 0, 256)
		for p := uintptr(unsafe.Pointer(cstr)); ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}

	return ""
}

func UTF16ToStringArray(s []uint16) []string {
	var ret []string
begin:
	for i, v := range s {
		if v == 0 {
			tmp := s[0:i]
			ret = append(ret, string(utf16.Decode(tmp)))
			if i+2 < len(s) && s[i+1] != 0 {
				s = s[i+1:]
				goto begin
			} else {
				break
			}
		}
	}
	return ret
}

// Convert a hex string to uint32
func HexToUint32(hexString string) (result uint32, err error) {
	var data []byte
	data, err = hex.DecodeString(hexString)
	if err == nil {
		result = binary.BigEndian.Uint32(data)
		return
	}
	if err != hex.ErrLength {
		return
	}
	hexString = "0" + hexString
	data, err = hex.DecodeString(hexString)
	if err == nil {
		result = binary.BigEndian.Uint32(data)
	}
	return
}

// IsErrSuccess checks if an "error" returned is actually the
// success code 0x0 "The operation completed successfully."
//
// This is the optimal approach since the error messages are
// localized depending on the OS language.
func IsErrSuccess(err error) bool {
	if errno, ok := err.(syscall.Errno); ok {
		if errno == 0 {
			return true
		}
	}
	return false
}
