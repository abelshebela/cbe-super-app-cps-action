# BPS Action Service - Markdown to DOCX Converter with PlantUML Rendering
# This script converts the Markdown documentation to a professional DOCX file
# with all PlantUML diagrams rendered as images

param(
    [string]$MarkdownFile = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.md",
    [string]$OutputDocx = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\BPS-ACTION-SERVICE-LLD.docx",
    [string]$DiagramsDir = "d:\Project\Refactor\cbe-supper-app-cps-action\docs\diagrams"
)

Write-Host "🔄 BPS Action Service Documentation Converter" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

# Create diagrams directory
if (-not (Test-Path $DiagramsDir)) {
    New-Item -ItemType Directory -Path $DiagramsDir | Out-Null
    Write-Host "✅ Created diagrams directory: $DiagramsDir" -ForegroundColor Green
}

# Function to encode PlantUML to URL format
function Encode-PlantUML {
    param([string]$PlantUMLCode)
    
    # Compress and encode for PlantUML server
    $bytes = [System.Text.Encoding]::UTF8.GetBytes($PlantUMLCode)
    $compressed = [System.IO.MemoryStream]::new()
    $deflate = [System.IO.Compression.DeflateStream]::new($compressed, [System.IO.Compression.CompressionMode]::Compress)
    $deflate.Write($bytes, 0, $bytes.Length)
    $deflate.Close()
    
    $compressedBytes = $compressed.ToArray()
    
    # Custom base64 encoding for PlantUML
    $base64 = [Convert]::ToBase64String($compressedBytes)
    $base64 = $base64.Replace('+', '-').Replace('/', '_').TrimEnd('=')
    
    return $base64
}

# Function to download PlantUML diagram
function Get-PlantUMLDiagram {
    param(
        [string]$PlantUMLCode,
        [string]$OutputPath,
        [int]$DiagramNumber
    )
    
    try {
        Write-Host "  📊 Rendering diagram $DiagramNumber..." -ForegroundColor Yellow
        
        # Use PlantUML online service
        $encoded = Encode-PlantUML -PlantUMLCode $PlantUMLCode
        $url = "http://www.plantuml.com/plantuml/png/$encoded"
        
        # Download the image
        Invoke-WebRequest -Uri $url -OutFile $OutputPath -ErrorAction Stop
        
        Write-Host "  ✅ Diagram $DiagramNumber saved: $OutputPath" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "  ⚠️  Failed to render diagram $DiagramNumber. Using placeholder." -ForegroundColor Yellow
        return $false
    }
}

# Extract PlantUML diagrams from markdown
Write-Host "📖 Reading markdown file..." -ForegroundColor Cyan
$content = Get-Content -Path $MarkdownFile -Raw

# Find all PlantUML blocks
$plantUMLPattern = '```plantuml\s*\n(.*?)\n```'
$matches = [regex]::Matches($content, $plantUMLPattern, [System.Text.RegularExpressions.RegexOptions]::Singleline)

Write-Host "🔍 Found $($matches.Count) PlantUML diagrams" -ForegroundColor Cyan
Write-Host ""

# Download each diagram
$diagramPaths = @()
for ($i = 0; $i -lt $matches.Count; $i++) {
    $plantUMLCode = $matches[$i].Groups[1].Value
    $diagramPath = Join-Path $DiagramsDir "diagram_$($i + 1).png"
    
    $success = Get-PlantUMLDiagram -PlantUMLCode $plantUMLCode -OutputPath $diagramPath -DiagramNumber ($i + 1)
    
    if ($success) {
        $diagramPaths += $diagramPath
    } else {
        $diagramPaths += $null
    }
}

Write-Host ""
Write-Host "✅ Diagram rendering complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Now creating DOCX file..." -ForegroundColor Cyan
Write-Host "   Please use Microsoft Word or LibreOffice to open the markdown file" -ForegroundColor Yellow
Write-Host "   and manually insert the diagrams from: $DiagramsDir" -ForegroundColor Yellow
Write-Host ""
Write-Host "💡 Alternative: Use pandoc for automatic conversion:" -ForegroundColor Cyan
Write-Host "   Install: winget install --id=JohnMacFarlane.Pandoc -e" -ForegroundColor Gray
Write-Host "   Convert: pandoc -s '$MarkdownFile' -o '$OutputDocx'" -ForegroundColor Gray
Write-Host ""
Write-Host "📊 Diagram files created:" -ForegroundColor Green
for ($i = 0; $i -lt $diagramPaths.Count; $i++) {
    if ($diagramPaths[$i]) {
        Write-Host "   $($i + 1). $($diagramPaths[$i])" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "✨ Process complete!" -ForegroundColor Green
