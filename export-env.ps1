# Usage: . .\export-env.ps1   <-- note the dot at the beginning (to persist in the terminal)

param (
    [string]$EnvFile = ".env"
)

# Check if file exists
if (-not (Test-Path $EnvFile)) {
    Write-Host "❌ Error: $EnvFile not found!"
    return
}

# Read and export each variable
Get-Content $EnvFile | ForEach-Object {
    $line = $_.Trim()

    # Skip empty lines or comments
    if ([string]::IsNullOrWhiteSpace($line) -or $line.StartsWith("#")) {
        return
    }

    # Split into key/value (only on the first '=')
    $parts = $line -split '=', 2
    if ($parts.Count -eq 2) {
        $key = $parts[0].Trim()
        $value = $parts[1].Trim()

        # Remove wrapping quotes
        # $value = $value 

        # Export to current terminal environment
        Set-Item -Path "env:$key" -Value $value
        Write-Host "✅ Exported $key"
    }
}

Write-Host "`n✅ All environment variables from $EnvFile are now available in this terminal session!"
