$initCommand = 'go get'

iex $initCommand

echo "Building for Windows"
$env:GOOS = "windows"

$winCommand = 'go build -o buildOutput\castro.exe'

iex $winCommand

echo "Building for Linux"
$env:GOOS = "linux"

$linuxCommand = 'go build -o buildOutput\castro'

iex $linuxCommand
