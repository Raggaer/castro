$date = (Get-Date).AddDays(-1).ToString('MM-dd-yyyy_HH:mm:ss')

$initCommand = 'go get'

iex $initCommand

echo "Building for Windows"
$env:GOOS = "windows"

$winCommand = 'go build -o buildOutput\castro.exe -ldflags "-X main.VERSION=$env:APPVEYOR_BUILD_VERSION -X main.BUILD_DATE=$date"'

iex $winCommand

echo "Building for Linux"
$env:GOOS = "linux"

$linuxCommand = 'go build -o buildOutput\castro'

iex $linuxCommand
