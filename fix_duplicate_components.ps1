# PowerShell script to fix duplicate components: sections in YAML files
$swaggerDir = "docs\swagger"
$yamlFiles = Get-ChildItem -Path $swaggerDir -Filter "*.yaml"

Write-Host "=== Fixing Duplicate Components Sections ===" -ForegroundColor Cyan
Write-Host ""

foreach ($file in $yamlFiles) {
    $filePath = $file.FullName
    $content = Get-Content $filePath -Raw -Encoding UTF8
    
    # Count components: sections
    $componentsMatches = [regex]::Matches($content, "(?m)^components:")
    
    if ($componentsMatches.Count -gt 1) {
        Write-Host "Processing: $($file.Name)" -ForegroundColor Yellow
        Write-Host "  Found $($componentsMatches.Count) components: sections"
        
        # Strategy: Find the second components: section and merge it with the first
        # The pattern is usually:
        # 1. First components: has securitySchemes and schemas
        # 2. Second components: has parameters and responses
        
        # Find the position of "# Common parameters" or similar comment before second components:
        $secondComponentsPattern = "(?ms)(^# Common.*?\r?\n)?^components:\s*\r?\n\s+parameters:"
        
        if ($content -match $secondComponentsPattern) {
            # Remove the second components: declaration and the comment before it
            $content = $content -replace "(?m)^# Common.*?\r?\n^components:\s*\r?\n", ""
            
            # Find the end of schemas section in first components (before security: or tags:)
            # and insert parameters and responses there
            $content = $content -replace "(?ms)(^\s+schemas:.*?)(^security:)", "`$1`n  parameters:`$2"
            
            Write-Host "  Attempting to merge sections..." -ForegroundColor Yellow
            
            # Actually, let's use a different approach - find and extract the parameters and responses
            # from the second location and move them to the first components section
            
            # Read the file again to get fresh content
            $content = Get-Content $filePath -Raw -Encoding UTF8
            
            # Extract parameters section from anywhere after the first components:
            $paramsPattern = "(?ms)^\s+parameters:\s*\r?\n((?:\s+\w+:.*?\r?\n(?:\s{4,}.*?\r?\n)*)*)"
            $paramsMatch = [regex]::Match($content, $paramsPattern)
            
            # Extract responses section
            $responsesPattern = "(?ms)^\s+responses:\s*\r?\n((?:\s+\w+:.*?\r?\n(?:\s{4,}.*?\r?\n)*)*)"
            $responsesMatch = [regex]::Match($content, $responsesPattern)
            
            if ($paramsMatch.Success -or $responsesMatch.Success) {
                # Find the end of the first components section (before security: or tags:)
                $insertPattern = "(?ms)(^components:.*?schemas:.*?)(^security:|^tags:)"
                
                if ($content -match $insertPattern) {
                    $beforeSecurity = $matches[1]
                    $securityOrTags = $matches[2]
                    
                    # Build the new components section
                    $newComponents = $beforeSecurity
                    
                    if ($paramsMatch.Success) {
                        $newComponents += "`n  parameters:`n" + $paramsMatch.Groups[1].Value
                    }
                    
                    if ($responsesMatch.Success) {
                        $newComponents += "`n  responses:`n" + $responsesMatch.Groups[1].Value
                    }
                    
                    # Remove the duplicate components: section at the end
                    $content = $content -replace "(?ms)^# Common.*?\r?\n^components:\s*\r?\n\s+parameters:.*$", ""
                    
                    # Replace the first components section with the merged one
                    $content = $content -replace "(?ms)(^components:.*?)(^security:|^tags:)", "$newComponents`n`$2"
                    
                    # Save the fixed content
                    Set-Content -Path $filePath -Value $content -Encoding UTF8 -NoNewline
                    Write-Host "  Fixed: $($file.Name)" -ForegroundColor Green
                } else {
                    Write-Host "  Could not find insertion point - manual fix needed" -ForegroundColor Red
                }
            } else {
                Write-Host "  Could not extract parameters/responses - manual fix needed" -ForegroundColor Red
            }
        } else {
            Write-Host "  Pattern not recognized - manual fix needed" -ForegroundColor Red
        }
    }
}

Write-Host "`n=== Verification ===" -ForegroundColor Cyan
foreach ($file in $yamlFiles) {
    $content = Get-Content $file.FullName -Raw
    $compCount = ([regex]::Matches($content, "(?m)^components:")).Count
    if ($compCount -gt 1) {
        Write-Host "$($file.Name): Still has $compCount components sections" -ForegroundColor Red
    } else {
        Write-Host "$($file.Name): OK ✓" -ForegroundColor Green
    }
}

Write-Host "`nDone!"
