:: Compiles and publishes the program online.
:: WORKS ONLY ON MY ENVIRONMENT!

:: Run from project root directory.
@echo off

setlocal

set pubdir=C:\Stuff\Dropbox\coppee\

:: Update version
call src\.scripts\update_version.bat
set /p version=<version

:: Compile
rmdir /s /q pkg bin
call gocross coppee

:: Mac 32-bit
cd bin\cross\darwin_386
call zip mac_32.zip coppee
move /y *.zip %pubdir%

:: Mac 64-bit
cd ..\darwin_amd64
call zip mac_64.zip coppee
move /y *.zip %pubdir%

:: Linux 32-bit
cd ..\linux_386
call zip linux_32.zip coppee
move /y *.zip %pubdir%

:: Linux 64-bit
cd ..\linux_amd64
call zip linux_64.zip coppee
move /y *.zip %pubdir%

:: Windows 32-bit
cd ..\windows_386
call zip windows_32.zip coppee.exe
move /y *.zip %pubdir%

:: Windows 64-bit
cd ..\windows_amd64
call zip windows_64.zip coppee.exe
move /y *.zip %pubdir%

:: Readme
cd ..\..\..
copy /y README.md %pubdir%

:: Update version
del %pubdir%.version*
copy nul "%pubdir%.version %version%"

:: Clean up
rmdir /s /q pkg bin

endlocal

:end
