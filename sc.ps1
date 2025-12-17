Write-Host "=== Stopping Go-related processes ==="

$processes = @("go.exe", "gopls.exe", "dlv.exe")

foreach ($p in $processes) {
    try {
        taskkill /F /IM $p /T 2>$null
        Write-Host "Killed process: $p"
    }
    catch {
        Write-Host "Process not running: $p"
    }
}

Write-Host "`n=== Closing VS Code (if running) ==="

try {
    taskkill /F /IM Code.exe /T 2>$null
    Write-Host "Closed VS Code"
}
catch {
    Write-Host "VS Code not running"
}

Write-Host "`n=== Cleaning Go cache ==="

go clean -cache -testcache

Write-Host "`n=== Cleaning Go modcache ==="
go clean --modcache

Write-Host "`n=== Forcing deletion of locked toolchain folder if needed ==="

$toolchain = "C:\Users\pc\go\pkg\mod\golang.org\toolchain@v0.0.1-go1.25.0.windows-amd64"

if (Test-Path $toolchain) {
    try {
        Remove-Item -Recurse -Force $toolchain
        Write-Host "Deleted locked toolchain folder successfully."
    }
    catch {
        Write-Warning "Failed to delete folder due to lock. Retrying with TAKEOWN + ICACLS..."

        takeown /F $toolchain /R /D Y
        icacls $toolchain /grant administrators:F /T

        Remove-Item -Recurse -Force $toolchain
        Write-Host "Deleted folder after ownership fix."
    }
}
else {
    Write-Host "Toolchain folder not found."
}

Write-Host "`n=== Fix completed successfully ==="
