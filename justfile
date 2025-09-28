# set shell := ["pwsh", "", "-CommandWithArgs"]
# set positional-arguments
shebang := 'pwsh.exe'
# Variables
exe_name := "4zupper"
mod_name := "zupper"
ld_flags :="-H=windowsgui -s -w -X spahttp/entity.Mode=production"
dist := ".dist"

default:
  just --list

win64:
    #!{{shebang}}
    $env:Path = "C:\Go\go.124\bin;C:\go\gcc\mingw64\bin;" + $env:Path
    $env:GOARCH = "amd64"
    $env:GOOS = "windows"
    $env:CGO_ENABLED = 1
    if (-Not (Test-Path go.mod)) {
      go mod init {{mod_name}}
    }
    go mod tidy -go 1.24 -v
    if(-Not $?) { exit }
    if (-Not (Test-Path "{{dist}}")) { New-Item -ItemType Directory -Force -Path "{{dist}}" | Out-Null }
    Remove-Item -Force -ErrorAction SilentlyContinue -LiteralPath "{{dist}}\{{exe_name}}.exe","{{dist}}\{{exe_name}}_64.exe"
    go build -ldflags="{{ld_flags}}" -o "{{dist}}\{{exe_name}}_64.exe" ./cmd
    if(-Not $?) { exit }
    upx --force-overwrite -o {{dist}}\{{exe_name}}.exe {{dist}}\{{exe_name}}_64.exe
