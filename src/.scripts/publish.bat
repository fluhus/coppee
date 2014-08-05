:: Compiles and publishes the program online.
:: WORKS ONLY ON MY ENVIRONMENT!

:: Run from project root directory.
@echo off

if "%1"=="" goto noargs
goto yesargs

:noargs
echo Usage:
echo publish ^<version number^>
goto end

:yesargs
setlocal

set pubdir=C:\Stuff\Dropbox\coppee\

:: Compile
rmdir /s /q pkg bin
call gocross coppee

:: Mac 32-bit
cd bin\cross\darwin_386
del *.zip
call zip mac_32.zip coppee
move /y *.zip %pubdir%

:: Mac 64-bit
cd ..\darwin_amd64
del *.zip
call zip mac_64.zip coppee
move /y *.zip %pubdir%

:: Linux 32-bit
cd ..\linux_386
del *.zip
call zip linux_32.zip coppee
move /y *.zip %pubdir%

:: Linux 64-bit
cd ..\linux_amd64
del *.zip
call zip linux_64.zip coppee
move /y *.zip %pubdir%

:: Windows 32-bit
cd ..\windows_386
del *.zip
call zip windows_32.zip coppee.exe
move /y *.zip %pubdir%

:: Windows 64-bit
cd ..\windows_amd64
del *.zip
call zip windows_64.zip coppee.exe
move /y *.zip %pubdir%

:: Readme
cd ..\..\..
copy /y README.md %pubdir%

:: Update version
del %pubdir%.Version*
copy nul "%pubdir%.Version %1"

:: Clean up
rmdir /s /q pkg bin

endlocal

:end
