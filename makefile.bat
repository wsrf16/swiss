SETLOCAL

set CGO_ENABLED=0

set GOARCH=amd64
set GOOS=linux

@REM set GOARCH=arm64
@REM set GOOS=linux

@REM set GOARCH=amd64
@REM set GOOS=windows

@REM set GOARM=7

set CC=x86_64-w64-mingw32-gcc
set CXX=CXX=x86_64-w64-mingw32-g++

go build -o .\output\swiss main.go

ENDLOCAL