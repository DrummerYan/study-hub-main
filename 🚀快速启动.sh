#!/bin/bash

echo "🎉 StudyHub 阶段 1 - 快速启动脚本"
echo "=================================="
echo ""

# 检查是否在正确的目录
if [ ! -f "阶段1-完成报告.md" ]; then
    echo "❌ 错误：请在 study-hub-main 目录下运行此脚本"
    exit 1
fi

echo "📋 启动前检查..."
echo ""

# 检查后端目录
if [ ! -d "server" ]; then
    echo "❌ 错误：找不到 server 目录"
    exit 1
fi

# 检查前端目录
if [ ! -d "web" ]; then
    echo "❌ 错误：找不到 web 目录"
    exit 1
fi

echo "✅ 目录检查通过"
echo ""

# 询问用户要启动哪个服务
echo "请选择要启动的服务："
echo "1) 只启动后端"
echo "2) 只启动前端"
echo "3) 同时启动后端和前端（推荐）"
echo ""
read -p "请输入选项 (1/2/3): " choice

case $choice in
    1)
        echo ""
        echo "🚀 启动后端服务..."
        cd server
        go run main.go
        ;;
    2)
        echo ""
        echo "🚀 启动前端服务..."
        cd web
        
        # 检查是否需要安装依赖
        if [ ! -d "node_modules" ]; then
            echo "📦 首次运行，正在安装依赖..."
            npm install
        fi
        
        npm run serve
        ;;
    3)
        echo ""
        echo "🚀 同时启动后端和前端..."
        echo ""
        
        # 检查是否安装了 tmux
        if command -v tmux &> /dev/null; then
            echo "使用 tmux 启动服务..."
            
            # 创建新会话并启动后端
            tmux new-session -d -s studyhub "cd server && go run main.go"
            
            # 分割窗口并启动前端
            tmux split-window -h "cd web && npm run serve"
            
            # 附加到会话
            echo "✅ 服务已在 tmux 中启动"
            echo "💡 使用 Ctrl+B 然后 D 可以分离会话"
            echo "💡 使用 tmux attach -t studyhub 可以重新连接"
            tmux attach -t studyhub
        else
            echo "⚠️  未检测到 tmux，将使用普通方式启动"
            echo ""
            echo "请打开两个终端窗口："
            echo ""
            echo "终端 1 - 启动后端："
            echo "cd /Users/yandrummer/Desktop/study-hub-main/server && go run main.go"
            echo ""
            echo "终端 2 - 启动前端："
            echo "cd /Users/yandrummer/Desktop/study-hub-main/web && npm run serve"
            echo ""
            echo "按回车键退出..."
            read
        fi
        ;;
    *)
        echo "❌ 无效选项"
        exit 1
        ;;
esac
