@echo off
echo Running OTA data seeder...
go run scripts/seed-ota-data.go
if %ERRORLEVEL% EQU 0 (
    echo Seeder completed successfully!
) else (
    echo Seeder encountered an error. Please check the logs.
)
pause 