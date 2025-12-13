$ErrorActionPreference = "Stop"

try {
    Write-Host "Starting conversion..."
    
    $htmlPath = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.html"
    $docxPath = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.docx"

    Write-Host "Source: $htmlPath"
    Write-Host "Target: $docxPath"

    $word = New-Object -ComObject Word.Application
    $word.Visible = $false
    $word.DisplayAlerts = 0 # wdAlertsNone

    Write-Host "Opening HTML..."
    # Open(FileName, ConfirmConversions, ReadOnly)
    $doc = $word.Documents.Open($htmlPath, $false, $true)
    
    Write-Host "Saving as DOCX..."
    # SaveAs2(FileName, FileFormat)
    $doc.SaveAs2($docxPath, 12) # wdFormatXMLDocument
    
    $doc.Close()
    $word.Quit()
    
    [System.Runtime.InteropServices.Marshal]::ReleaseComObject($word) | Out-Null
    
    Write-Host "SUCCESS: Created $docxPath" -ForegroundColor Green
}
catch {
    Write-Error "FATAL ERROR: $($_.Exception.Message)"
    if ($word) {
        $word.Quit()
        [System.Runtime.InteropServices.Marshal]::ReleaseComObject($word) | Out-Null
    }
    exit 1
}
