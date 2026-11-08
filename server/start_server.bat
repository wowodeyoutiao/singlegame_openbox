@echo off
setlocal
cd /d "%~dp0nanoserver"

set "PY=C:\Users\Administrator\.workbuddy\binaries\python\versions\3.14.3\python.exe"
if not exist "%PY%" set "PY=python"
set "REDIS=E:\DevelopTools\Xoreax\IncrediBuild\redis-server.exe"

echo ==========================================
echo   Start WLY Local Server (single server 9)
echo ==========================================

echo [0/3] checking port 7000 ...
netstat -ano | findstr /c:":7000 " | findstr "LISTENING" >nul
if not errorlevel 1 (
  echo [ERROR] port 7000 is already in use.
  echo         Run stop_server.bat first, then start again.
  goto fail
)

echo [1/3] go-admin mock   127.0.0.1:18080
start "WLY-MOCK" /min "%PY%" local_mock.py
echo [2/3] redis           127.0.0.1:6379
start "WLY-REDIS" /min "%REDIS%" --port 6379 --save "" --appendonly no
echo [3/3] game server     7000 / 12307
echo        first run takes 30-60s to load configs and sync tables.
echo        Press Ctrl+C here to stop the game server.
echo.
set "GOMEMLIMIT=2GiB"
build\server.exe -c .\configs\config.toml

echo.
echo [game server stopped]
echo Close the WLY-MOCK and WLY-REDIS windows to stop the rest.
pause
exit /b 0

:fail
pause
exit /b 1
