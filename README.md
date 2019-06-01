# gup
basic file uploader/remover for [gohst](https://github.com/voidiz/gohst)

## requirements
tested with
- `go 1.11+` (if building from source)

## building and installing
`go get github.com/voidiz/gup`

## prebuilt binaries
https://github.com/voidiz/gup/releases

Substitute `gup` with the path to the executable and run the 
[quick start](#quick-start-and-basic-usage) commands in cmd/PowerShell. 
Alternatively, add the folder where the executable is located to your 
PATH environment variable and proceed as usual.

## quick start and basic usage
1. `gup config <host> <username> <password>` - Creates a configuration file where
your auth token (generated using the supplied username and password) and the 
host domain is stored. E.g. `gup config https://example.com user secretpass`
2. `gup <file>` - Uploads a file using the given filepath. E.g. `gup image.jpg` to
upload "image.jpg" located in the current folder.
3. `gup delete <url>` - Deletes a file using the supplied url pointing to the file.
