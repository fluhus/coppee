:: Updates version number in all relevant places, according
:: to the version in file "version".
:: DESIGNED TO WORK ONLY ON MY ENVIRONMENT!

:: Execute from project root directory.
@echo off

setlocal

:: Read version from file
set /p version=<version

:: Update readme
echo Coppee (%version%)>README_new.md
tail -n+2 README.md>>README_new.md
move /y README_new.md README.md

:: Update code
set gofile=src\coppee\version_info.go
echo package main>%gofile%
echo.>>%gofile%
echo // Auto-generated version info constant.>>%gofile%
echo const version="%version%">>%gofile%

endlocal
