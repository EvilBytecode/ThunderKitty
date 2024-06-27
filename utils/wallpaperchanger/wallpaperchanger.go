package WallpaperChanger

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

const S = 2

var (
	SW = uint32(0x0014)
	U  = uint32(0x01)
	SC = uint32(0x02)
)

func DownloadAndSetWallpaper() {
	url := "https://cataas.com/cat/says/THUNDER%20KITTY%20RUNS%20YOU?filter=custom&brightness=1.5&saturation=50"
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, "Thunder_Kitty.jpg")

	r, _ := http.Get(url)
	defer r.Body.Close()

	p, _ := os.Create(filePath)
	defer p.Close()

	io.Copy(p, r.Body)

	setWallpaper(filePath, S)
	exec.Command("taskkill", "/F", "/IM", "wallpaper32.exe").Run()
}

func setWallpaper(path string, style int) {
	k, _ := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE)
	defer k.Close()
	k.SetStringValue("TileWallpaper", "0")
	k.SetStringValue("WallpaperStyle", fmt.Sprintf("%d", style))
	p, _ := syscall.UTF16PtrFromString(path)
	syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW").Call(uintptr(SW), uintptr(0), uintptr(unsafe.Pointer(p)), uintptr(U|SC))
}
