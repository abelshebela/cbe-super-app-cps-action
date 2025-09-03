# PowerShell script to update all Swagger YAML files with correct base path
$basePath = "/api/v1/cbesuperapp/cps_action"
$swaggerDir = "docs\swagger"

# List of files to update (excluding already updated ones)
$filesToUpdate = @(
    "wallet.yaml", "mini_app.yaml", "feedback.yaml", "event.yaml", "customer.yaml",
    "advertisement.yaml", "avatar.yaml", "budget.yaml", "product_code.yaml", 
    "service.yaml", "hq.yaml", "password_rule.yaml", "bulk_service.yaml",
    "mini_app_merchant.yaml", "account_validation.yaml", "bps_user.yaml",
    "fayda.yaml", "portal_card.yaml", "unlink.yaml", "notification.yaml",
    "donation.yaml", "branch.yaml", "budget_category.yaml", "action.yaml",
    "account_search.yaml", "key_generator.yaml", "permission.yaml"
)

foreach ($file in $filesToUpdate) {
    $filePath = Join-Path $swaggerDir $file
    if (Test-Path $filePath) {
        Write-Host "Updating $file..."
        $content = Get-Content $filePath -Raw
        
        # Replace paths that start with "  /" with the base path
        $content = $content -replace '(?m)^(\s{2})/([^/])', "`$1$basePath/`$2"
        
        Set-Content $filePath $content -NoNewline
        Write-Host "Updated $file"
    } else {
        Write-Host "File not found: $file"
    }
}

Write-Host "All Swagger files updated with base path: $basePath"
