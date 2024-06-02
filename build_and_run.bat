@echo off

REM Build the Go project
go build .

REM Check if the build was successful
if %ERRORLEVEL% == 0 (
    echo Build successful, running asteroids.exe...
    REM Execute the resulting executable
    asteroids.exe
) else (
    echo Build failed, not running asteroids.exe.
)