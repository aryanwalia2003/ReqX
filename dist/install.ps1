$installDir = join-path $env:LOCALAPPDATA "postman-cli"
$exeName = "postman-cli.exe"
$sourceExe = join-path $PSScriptRoot $exeName

# 1. Create install directory if it doesn't exist
if (-not (test-path $installDir)) {
    write-host "📁 Creating installation directory: $installDir" -foregroundColor Cyan
    new-item -path $installDir -itemType Directory | out-null
}

# 2. Check if source EXE exists
if (test-path $sourceExe) {
    write-host "📦 Copying $exeName to $installDir..." -foregroundColor Cyan
    copy-item -path $sourceExe -destination $installDir -force
} else {
    write-host "❌ Error: $exeName not found in the current directory!" -foregroundColor Red
    write-host "Please run 'go build -o $exeName main.go' first." -foregroundColor Yellow
    exit 1
}

# 3. Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    write-host "🛠️ Adding $installDir to User PATH..." -foregroundColor Cyan
    $newPath = "$userPath;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    write-host "✅ PATH updated successfully!" -foregroundColor Green
} else {
    write-host "ℹ️  $installDir is already in your PATH." -foregroundColor Gray
}

write-host ""
write-host "🎉 postman-cli has been installed successfully!" -foregroundColor Green
write-host "🚀 Restart your terminal and type 'postman-cli --help' to get started." -foregroundColor Yellow
