$date = (Get-Date).AddDays(-1).ToString('MM-dd-yyyy_HH:mm:ss')

$initCommand = 'go get'

iex $initCommand

echo "Building for Windows amd64"
$env:GOOS = "windows"

$winCommand = 'go build -v github.com/raggaer/castro -o buildOutput\castro_win_amd64.exe -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

iex $winCommand

echo "Building for Linux amd64"
$env:GOOS = "linux"

$linuxCommand = 'go build -v github.com/raggaer/castro -o buildOutput\castro_linux_amd64 -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

iex $linuxCommand

$linuxCommand = 'go build -v github.com/raggaer/castro -o buildOutput\castro_linux_arm64 -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

echo "Building for Linux arm64"
$env:GOOS = "linux"
$env:GOARCH = "arm64"

iex $linuxCommand




