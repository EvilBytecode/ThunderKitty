@Echo off
Setlocal EnableExtensions

:: Check if Go is installed
where go > NUL 2>&1
if %ERRORLEVEL% neq 0 (
    Echo Go is not installed. Please install Go from https://golang.org/dl/ and try again.
    pause
    exit /b 1
)

:: Install required modules
go get -u fyne.io/fyne/v2@latest
go get -u golang.org/x/crypto/pbkdf2
go get -u modernc.org/sqlite
go get -u github.com/shirou/gopsutil/disk
go get -u github.com/EvilBytecode/GoDefender

:: Check if all modules were installed successfully
if %ERRORLEVEL% neq 0 (
    Echo Error installing modules. Please check the output above for errors.
    pause
    exit /b 1
)

Echo Modules installed successfully!
pause