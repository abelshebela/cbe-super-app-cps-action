# PowerShell script to fix YAML files with duplicate path issues
$swaggerDir = "docs\swagger"
$basePath = "/api/v1/cbesuperapp/cps_action"

# Get all YAML files except account_block.yaml (already correct)
$yamlFiles = Get-ChildItem -Path $swaggerDir -Filter "*.yaml" | Where-Object { $_.Name -ne "account_block.yaml" -and $_.Name -ne "index.html" }

foreach ($file in $yamlFiles) {
    Write-Host "`nProcessing: $($file.Name)"
    $filePath = $file.FullName
    $content = Get-Content $filePath -Raw -Encoding UTF8
    
    # Check if file has the base path in paths section
    if ($content -match "paths:\s+/api/v1/cbesuperapp/cps_action") {
        Write-Host "  Found base path in paths section - fixing..."
        
        # Remove basePath field if it exists (OpenAPI 3.0 doesn't use it)
        $content = $content -replace "(?m)^# Base path.*\r?\n", ""
        $content = $content -replace "(?m)^basePath:.*\r?\n", ""
        
        # Replace paths that include the full base path
        $content = $content -replace "(?m)^(\s{2})/api/v1/cbesuperapp/cps_action/", "`$1/"
        
        # Check for duplicate components: sections
        $componentsCount = ([regex]::Matches($content, "(?m)^components:")).Count
        if ($componentsCount -gt 1) {
            Write-Host "  WARNING: Multiple 'components:' sections found - manual review needed"
        }
        
        # Save the fixed content
        Set-Content -Path $filePath -Value $content -Encoding UTF8 -NoNewline
        Write-Host "  Fixed: $($file.Name)" -ForegroundColor Green
    } else {
        Write-Host "  Skipped: Already correct" -ForegroundColor Yellow
    }
}

Write-Host "`n=== Validation ==="
Write-Host "Checking for remaining issues..."

foreach ($file in $yamlFiles) {
    $content = Get-Content $file.FullName -Raw
    
    # Check for duplicate paths
    $pathLines = $content -split "`n" | Where-Object { $_ -match "^\s{2}/" }
    $duplicates = $pathLines | Group-Object | Where-Object { $_.Count -gt 1 }
    
    if ($duplicates) {
        Write-Host "`nWARNING: $($file.Name) has duplicate paths:" -ForegroundColor Red
        foreach ($dup in $duplicates) {
            Write-Host "  $($dup.Name) (appears $($dup.Count) times)" -ForegroundColor Red
        }
    }
}

Write-Host "`nDone! All YAML files have been processed."
