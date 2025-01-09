# PowerShell Script to Set Environment Variables

# powershell -ExecutionPolicy Bypass -File .\setenv.ps1


# Define the environment variables
$envVariables = @{
    "MAIL_PORT"       = "1025"
    "MAIL_DOMAIN"     = "localhost"
    "MAIL_HOST"       = "mail"
    "MAIL_ENCRYPTION" = "none"
    "MAIL_USERNAME"   = ""
    "MAIL_PASSWORD"   = ""
    "MAIL_FROMNAME"   = "John Smith"
    "MAIL_FROMADDR"   = "JohnSmith@example.com"
}

# Loop through each variable and set it as an environment variable
foreach ($key in $envVariables.Keys) {
    [System.Environment]::SetEnvironmentVariable($key, $envVariables[$key], [System.EnvironmentVariableTarget]::User)
    Write-Host "Set environment variable $key to $($envVariables[$key])"
}

Write-Host "Environment variables set successfully. Please restart your session or use `RefreshEnv` to apply."

#dir env: