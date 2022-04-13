// +build windows
package Log

import (
	"syscall"
)

var kernel32  *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
var SetColorProc   *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
var CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)

func PrintColorText(s string, color int) {
	if instance.curColor == color {
		print(s)
		return
	} else {
		handle, _, _ := SetColorProc.Call(uintptr(syscall.Stdout), uintptr(color))
		print(s)
		CloseHandle.Call(handle)
		instance.curColor = color
	}
}
