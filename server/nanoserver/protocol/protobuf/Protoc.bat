@echo off
:: 1. 设好路径
set PROTOC=protoc.exe
set SRC_DIR=proto
set OUT_DIR=cs

:: 2. 创建输出目录
if not exist "%OUT_DIR%" mkdir "%OUT_DIR%"

:: 3. 逐个编译
for %%f in ("%SRC_DIR%\*.proto") do (
    echo [protoc] %%~nf.cs
    "%PROTOC%" -I="%SRC_DIR%" --csharp_out="%OUT_DIR%" "%%f"
)

echo 全部完成！
pause