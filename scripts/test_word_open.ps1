try {
    $path = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.html"
    Write-Host "Creating Word Object..."
    $word = New-Object -ComObject Word.Application
    $word.Visible = $false
    
    Write-Host "Opening Document: $path"
    $doc = $word.Documents.Open($path)
    Write-Host "Document Opened."
    
    $doc.Close()
    $word.Quit()
} catch {
    Write-Error $_
}
