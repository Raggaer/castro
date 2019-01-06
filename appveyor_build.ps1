If (Test-Path "buildOutput") {
	Remove-Item "buildOutput" -recurse
}

$date = (Get-Date).AddDays(-1).ToString('MM-dd-yyyy_HH:mm:ss')
$version = git rev-parse HEAD

echo "$date - $version"
echo "Building for Windows amd64"

$env:GOOS = "windows"
$env:GOARCH = "amd64"

$winCommand = 'go build -o buildOutput\castro_win_amd64.exe -ldflags "-X github.com/raggaer/castro/app/util.VERSION=$version -X github.com/raggaer/castro/app/util.BUILD_DATE=$date"'

iex $winCommand

echo "Building for Linux amd64"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

$linuxCommand = 'go build -o buildOutput\castro_linux_amd64 -ldflags "-X github.com/raggaer/castro/app/util.VERSION=$version -X github.com/raggaer/castro/app/util.BUILD_DATE=$date"'

iex $linuxCommand

$linuxCommand = 'go build -o buildOutput\castro_linux_arm64 -ldflags "-X github.com/raggaer/castro/app/util.VERSION=$version -X github.com/raggaer/castro/app/util.BUILD_DATE=$date"'

echo "Building for Linux arm64"
$env:GOOS = "linux"
$env:GOARCH = "arm64"

iex $linuxCommand

echo "Creating data directories"

Copy-Item pages buildOutput\data\pages -recurse
Copy-Item widgets buildOutput\data\widgets -recurse
Copy-Item install buildOutput\data\install -recurse
Copy-Item public buildOutput\data\public -recurse
Copy-Item views buildOutput\data\views -recurse
Copy-Item engine buildOutput\data\engine -recurse
Copy-Item migrations buildOutput\data\migrations -recurse
Copy-Item extensions buildOutput\data\extensions -recurse
Copy-Item i18n buildOutput\data\i18n -recurse
Copy-Item LICENSE buildOutput\LICENSE

echo "Compressing data directories"

$files = Get-ChildItem -Path "buildOutput\data\*"

Compress-Archive -Path $files -CompressionLevel Optimal -DestinationPath buildOutput\release.zip


