@echo off
echo 启动考古家数据采集系统...
echo.

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到Go环境，请先安装Go
    pause
    exit /b 1
)

REM 检查MongoDB是否运行
echo 检查MongoDB连接...
mongosh --eval "db.runCommand('ping')" >nul 2>&1
if errorlevel 1 (
    echo 警告: MongoDB可能未运行，请确保MongoDB服务已启动
    echo.
)

REM 安装依赖
echo 安装依赖包...
go mod tidy

REM 编译程序
echo 编译程序...
go build -o collyDemo.exe

REM 运行程序
echo 启动程序...
echo 按 Ctrl+C 停止程序
echo.
collyDemo.exe

pause 