$date = (Get-Date).AddDays(-1).ToString('MM-dd-yyyy_HH:mm:ss')

$initCommand = 'go get'

iex $initCommand

echo "Building for Windows amd64"
$env:GOOS = "windows"

$winCommand = 'go build -o buildOutput\castro_win_amd64.exe -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

iex $winCommand

echo "Building for Linux amd64"
$env:GOOS = "linux"

$linuxCommand = 'go build -o buildOutput\castro_linux_amd64 -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

iex $linuxCommand

echo "Building for Linux arm64"
$env:GOOS = "linux"
$env:GOARCH = "arm64"

iex $linuxCommand




