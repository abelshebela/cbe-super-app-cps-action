try {
    Write-Host "Creating Word Object..."
    $word = New-Object -ComObject Word.Application
    Write-Host "Word Object Created."
    $word.Quit()
} catch {
    Write-Error $_
}
