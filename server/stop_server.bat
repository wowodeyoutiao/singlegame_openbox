@echo off
setlocal
echo ==========================================
echo   Stop WLY Local Server
echo ==========================================

for %%P in (7000 12307 6379 18080) do (
  for /f "tokens=5" %%a in ('netstat -ano ^| findstr /c:":%%P " ^| findstr "LISTENING"') do (
    echo [stop] port %%P - PID %%a
    taskkill /F /PID %%a >nul 2>&1
  )
)

echo.
echo [OK] game server, mock and redis stopped.
echo      MySQL service is NOT touched (shared by other projects).
exit /b 0
