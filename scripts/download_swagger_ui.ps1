# Create directories if they don't exist
$swaggerDir = "$PSScriptRoot/../internal/handlers/rest/http/swagger/swagger-ui"
New-Item -ItemType Directory -Force -Path $swaggerDir | Out-Null

# Download Swagger UI files
$swaggerUIVersion = "5.15.0"
$swaggerUIFiles = @(
    "swagger-ui-bundle.js",
    "swagger-ui-es-bundle-core.js",
    "swagger-ui-es-bundle.js",
    "swagger-ui-standalone-preset.js",
    "swagger-ui.css",
    "swagger-ui.js",
    "swagger-ui.css.map",
    "swagger-ui.js.map",
    "favicon-16x16.png",
    "favicon-32x32.png",
    "swagger-ui-bundle.js.map",
    "swagger-ui-standalone-preset.js.map"
)

foreach ($file in $swaggerUIFiles) {
    $url = "https://unpkg.com/swagger-ui-dist@$swaggerUIVersion/$file"
    $output = "$swaggerDir/$file"
    Write-Host "Downloading $url to $output"
    Invoke-WebRequest -Uri $url -OutFile $output -UseBasicParsing
}

# Create the index.html file
$indexContent = @"
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>CPS Action API Documentation</title>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *,
    *:before,
    *:after {
      box-sizing: inherit;
    }
    body {
      margin: 0;
      background: #fafafa;
    }
    .swagger-ui .topbar {
      background-color: #1e3a8a;
      padding: 10px 0;
    }
    .swagger-ui .topbar .topbar-wrapper {
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .swagger-ui .topbar .topbar-wrapper img {
      content: url('data:image/svg+xml;charset=UTF-8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" width="24px" height="24px"><path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/></svg>');
      height: 40px;
      margin-right: 10px;
    }
    .swagger-ui .topbar .topbar-wrapper .topbar-title {
      color: white;
      font-size: 1.5em;
      font-weight: bold;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="./swagger-ui-bundle.js" charset="UTF-8"></script>
  <script src="./swagger-ui-standalone-preset.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      // Get all available API specs
      const apiSpecs = [
        {url: "/swagger/bank.yaml", name: "Bank API"},
        {url: "/swagger/account_block.yaml", name: "Account Block API"},
        {url: "/swagger/account_validation.yaml", name: "Account Validation API"},
        {url: "/swagger/wallet.yaml", name: "Wallet API"},
        {url: "/swagger/notification.yaml", name: "Notification API"},
        // More API specs will be added automatically
      ];

      // Begin Swagger UI call parameters
      const ui = SwaggerUIBundle({
        urls: apiSpecs,
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        validatorUrl: null,
        defaultModelsExpandDepth: -1,
        defaultModelExpandDepth: 3,
        docExpansion: 'none',
        defaultModelExpandDepth: 1,
        defaultModelsExpandDepth: 1,
        filter: true,
        persistAuthorization: true,
        displayRequestDuration: true,
        showCommonExtensions: true,
        showExtensions: true,
        showMutatedRequest: true,
        tagsSorter: 'alpha',
        operationsSorter: 'alpha'
      });

      // Add API key input field
      const authorizeBtn = document.createElement('div');
      authorizeBtn.innerHTML = `
        <div class="scheme-container">
          <div class="auth-wrapper">
            <div class="btn authorize">
              <span>Authorize</span>
            </div>
          </div>
        </div>
      `;
      
      // Add API key input to the top bar
      const topbar = document.querySelector('.topbar');
      if (topbar) {
        topbar.appendChild(authorizeBtn);
      }

      // Handle API key submission
      const authorizeButton = document.querySelector('.btn.authorize');
      if (authorizeButton) {
        authorizeButton.addEventListener('click', function() {
          const token = prompt('Enter your JWT token:');
          if (token) {
            // Set the token in the authorization header
            ui.authActions.authorize({
              BearerAuth: {
                name: 'Authorization',
                schema: {
                  type: 'http',
                  scheme: 'bearer',
                  bearerFormat: 'JWT'
                },
                value: token
              }
            });
            
            // Store the token in localStorage for page refreshes
            localStorage.setItem('jwt_token', token);
          }
        });
      }

      // Load token from localStorage if available
      const savedToken = localStorage.getItem('jwt_token');
      if (savedToken) {
        ui.authActions.authorize({
          BearerAuth: {
            name: 'Authorization',
            schema: {
              type: 'http',
              scheme: 'bearer',
              bearerFormat: 'JWT'
            },
            value: savedToken
          }
        });
      }
    };
  </script>
</body>
</html>
"@

$indexContent | Out-File -FilePath "$swaggerDir/index.html" -Encoding utf8

Write-Host "Swagger UI setup complete!"
