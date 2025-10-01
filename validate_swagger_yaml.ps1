# PowerShell script to validate YAML files for common issues
$swaggerDir = "docs\swagger"

Write-Host "=== Swagger YAML Validation ===" -ForegroundColor Cyan
Write-Host ""

$yamlFiles = Get-ChildItem -Path $swaggerDir -Filter "*.yaml"
$totalFiles = $yamlFiles.Count
$issuesFound = 0

foreach ($file in $yamlFiles) {
    Write-Host "Checking: $($file.Name)" -ForegroundColor White
    $filePath = $file.FullName
    $content = Get-Content $filePath -Raw -Encoding UTF8
    $lines = $content -split "`r?`n"
    $fileIssues = @()
    
    # Check 1: Duplicate paths
    $pathPattern = "^\s{2}/"
    $paths = @{}
    for ($i = 0; $i -lt $lines.Count; $i++) {
        if ($lines[$i] -match $pathPattern) {
            $path = $lines[$i].Trim()
            if ($paths.ContainsKey($path)) {
                $fileIssues += "  Line $($i+1): Duplicate path '$path' (first seen at line $($paths[$path]))"
            } else {
                $paths[$path] = $i + 1
            }
        }
    }
    
    # Check 2: Duplicate components: sections
    $componentsMatches = [regex]::Matches($content, "(?m)^components:")
    if ($componentsMatches.Count -gt 1) {
        $fileIssues += "  Multiple 'components:' sections found ($($componentsMatches.Count) occurrences)"
    }
    
    # Check 3: Duplicate security: sections at root level
    $securityMatches = [regex]::Matches($content, "(?m)^security:")
    if ($securityMatches.Count -gt 1) {
        $fileIssues += "  Multiple root-level 'security:' sections found ($($securityMatches.Count) occurrences)"
    }
    
    # Check 4: Duplicate tags: sections
    $tagsMatches = [regex]::Matches($content, "(?m)^tags:")
    if ($tagsMatches.Count -gt 1) {
        $fileIssues += "  Multiple 'tags:' sections found ($($tagsMatches.Count) occurrences)"
    }
    
    # Check 5: Duplicate paths: sections
    $pathsMatches = [regex]::Matches($content, "(?m)^paths:")
    if ($pathsMatches.Count -gt 1) {
        $fileIssues += "  Multiple 'paths:' sections found ($($pathsMatches.Count) occurrences)"
    }
    
    # Check 6: Invalid basePath (OpenAPI 3.0 doesn't use it)
    if ($content -match "(?m)^basePath:") {
        $fileIssues += "  Invalid 'basePath' field found (not supported in OpenAPI 3.0)"
    }
    
    # Check 7: Paths still containing full base URL
    if ($content -match "(?m)^\s{2}/api/v1/cbesuperapp/cps_action/") {
        $fileIssues += "  Paths contain full base URL (should be relative to server URL)"
    }
    
    # Check 8: Inconsistent indentation in YAML
    $indentIssues = 0
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        if ($line -match "^(\s+)" -and $line -notmatch "^\s*#" -and $line.Trim() -ne "") {
            $indent = $matches[1].Length
            if ($indent % 2 -ne 0) {
                $indentIssues++
            }
        }
    }
    if ($indentIssues -gt 0) {
        $fileIssues += "  $indentIssues lines with odd indentation (should be multiples of 2)"
    }
    
    # Report results for this file
    if ($fileIssues.Count -gt 0) {
        Write-Host "  ❌ Issues found:" -ForegroundColor Red
        foreach ($issue in $fileIssues) {
            Write-Host $issue -ForegroundColor Yellow
        }
        $issuesFound += $fileIssues.Count
    } else {
        Write-Host "  ✓ No issues found" -ForegroundColor Green
    }
    Write-Host ""
}

Write-Host "=== Summary ===" -ForegroundColor Cyan
Write-Host "Total files checked: $totalFiles"
if ($issuesFound -eq 0) {
    Write-Host "All files are valid! ✓" -ForegroundColor Green
} else {
    Write-Host "Total issues found: $issuesFound" -ForegroundColor Yellow
    Write-Host "Please review and fix the issues above." -ForegroundColor Yellow
}
