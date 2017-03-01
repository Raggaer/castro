$date = (Get-Date).AddDays(-1).ToString('MM-dd-yyyy HH:mm:ss')

$initCommand = 'go get'

iex $initCommand

echo "Building for Windows"
$env:GOOS = "windows"

$winCommand = 'go build -o buildOutput\castro.exe -ldflags "-X main.VERSION=1.0.0.{build} -X main.BUILD_DATE=$date"'

iex $winCommand

echo "Building for Linux"
$env:GOOS = "linux"

$linuxCommand = 'go build -o buildOutput\castro'

iex $linuxCommand
