@echo off
chcp 65001 >nul
echo === 在线考试系统后端部署脚本 (Windows) ===
echo.

REM 检查是否存在可执行文件
if not exist "exam-backend-windows-amd64.exe" (
    echo 错误: 找不到 exam-backend-windows-amd64.exe 文件
    echo 请先运行编译命令: go build -o exam-backend-windows-amd64.exe main.go
    pause
    exit /b 1
)

REM 检查是否存在环境配置文件
if not exist ".env" (
    echo 警告: 找不到 .env 文件，将使用 .env.example 作为模板
    if exist ".env.example" (
        copy ".env.example" ".env" >nul
        echo 已复制 .env.example 到 .env，请根据实际情况修改配置
    ) else (
        echo 错误: 找不到 .env.example 文件
        pause
        exit /b 1
    )
)

REM 创建部署包
echo 正在创建部署包...
set DEPLOY_DIR=exam-backend-deploy-windows
if