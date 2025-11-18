Write-Host "=== godror Windows Fix Script ===" -ForegroundColor Cyan

# ---------------------------------------------------
# 1. Remove MinGW (gcc) – causes CGo object parse errors
# ---------------------------------------------------
Write-Host "`n[1] Checking for MinGW/GCC..." -ForegroundColor Yellow
$gcc = Get-Command gcc -ErrorAction SilentlyContinue

if ($gcc) {
    Write-Host "Found GCC at $($gcc.Source). Removing from PATH..." -ForegroundColor Red
    
    $oldPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
    $newPath = ($oldPath -split ';' | Where-Object {$_ -notmatch "mingw"} ) -join ';'

    [Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
    Write-Host "MinGW removed from PATH." -ForegroundColor Green
} else {
    Write-Host "GCC not found 👍" -ForegroundColor Green
}

# ---------------------------------------------------
# 2. Configure Oracle Instant Client
# ---------------------------------------------------
Write-Host "`n[2] Setting Oracle Instant Client paths..." -ForegroundColor Yellow

$oraclePath = "C:\oracle\instantclient_19_23"

if (!(Test-Path $oraclePath)) {
    Write-Host "Oracle Instant Client folder not found at $oraclePath" -ForegroundColor Red
    Write-Host "Please download from: https://www.oracle.com/database/technologies/instant-client/downloads.html"
    exit
}

# Add to PATH
$path = [Environment]::GetEnvironmentVariable("PATH", "Machine")
if ($path -notmatch [regex]::Escape($oraclePath)) {
    [Environment]::SetEnvironmentVariable("PATH", "$path;$oraclePath", "Machine")
    Write-Host "InstantClient added to PATH" -ForegroundColor Green
} else {
    Write-Host "InstantClient already in PATH 👍" -ForegroundColor Green
}

# Set ORACLE_HOME & LD_LIBRARY_PATH
[Environment]::SetEnvironmentVariable("ORACLE_HOME", $oraclePath, "Machine")
[Environment]::SetEnvironmentVariable("LD_LIBRARY_PATH", $oraclePath, "Machine")

Write-Host "Oracle env variables set ✔" -ForegroundColor Green

# ---------------------------------------------------
# 3. Enable CGO
# ---------------------------------------------------
Write-Host "`n[3] Enabling CGO..." -ForegroundColor Yellow
[Environment]::SetEnvironmentVariable("CGO_ENABLED", "1", "Machine")
Write-Host "CGO_ENABLED=1 ✔" -ForegroundColor Green

# ---------------------------------------------------
# 4. Make sure MSVC (cl.exe) exists
# ---------------------------------------------------
Write-Host "`n[4] Checking for Microsoft C++ compiler (cl.exe)..." -ForegroundColor Yellow

$cl = Get-Command cl -ErrorAction SilentlyContinue

if (!$cl) {
    Write-Host "❌ cl.exe NOT found!" -ForegroundColor Red
    Write-Host "Install MSVC Build Tools: https://visualstudio.microsoft.com/visual-cpp-build-tools/" -ForegroundColor White
    exit
} else {
    Write-Host "Found cl.exe at $($cl.Source) ✔" -ForegroundColor Green
}

# ---------------------------------------------------
# 5. Run test build
# ---------------------------------------------------
Write-Host "`n[5] Testing Go build with godror..." -ForegroundColor Yellow

try {
    go build -tags oracle ./... 2>&1 | Out-String | Write-Host
    Write-Host "`nBuild OK ✔ godror is working!" -ForegroundColor Green
}
catch {
    Write-Host "`n❌ Build still failing. Send output to ChatGPT." -ForegroundColor Red
}

Write-Host "`n=== Done ===" -ForegroundColor Cyan
