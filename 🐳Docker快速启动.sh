#!/bin/bash

echo "🐳 StudyHub Docker 快速启动工具"
echo "================================"
echo ""

# 配置文件路径
COMPOSE_FILE="deploy/docker-compose/docker-compose.yaml"
PROJECT_ROOT="$HOME/Desktop/study-hub-main"

# 进入项目目录
cd "$PROJECT_ROOT" || exit 1

echo "📍 当前目录: $(pwd)"
echo ""

# 检查配置文件是否存在
if [ ! -f "$COMPOSE_FILE" ]; then
    echo "❌ 找不到配置文件: $COMPOSE_FILE"
    exit 1
fi

echo "✅ 找到配置文件: $COMPOSE_FILE"
echo ""

# 显示菜单
echo "请选择操作："
echo "1) 启动所有服务 (docker compose up -d)"
echo "2) 停止所有服务 (docker compose down)"
echo "3) 重启所有服务 (docker compose restart)"
echo "4) 查看服务状态 (docker ps)"
echo "5) 查看 MySQL 日志"
echo "6) 退出"
echo ""

read -p "请输入选项 (1-6): " choice

case $choice in
    1)
        echo ""
        echo "🚀 正在启动服务..."
        docker compose -f "$COMPOSE_FILE" up -d
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ 服务启动成功！"
            echo ""
            echo "📋 运行中的容器："
            docker ps --filter "name=gva-"
            echo ""
            echo "📝 提示："
            echo "- MySQL 容器: gva-mysql (端口 13306:3306)"
            echo "- Redis 容器: gva-redis (端口 16379:6379)"
        else
            echo ""
            echo "❌ 启动失败！请检查错误信息"
        fi
        ;;
    
    2)
        echo ""
        echo "🛑 正在停止服务..."
        docker compose -f "$COMPOSE_FILE" down
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ 服务已停止"
        else
            echo ""
            echo "❌ 停止失败！请检查错误信息"
        fi
        ;;
    
    3)
        echo ""
        echo "🔄 正在重启服务..."
        docker compose -f "$COMPOSE_FILE" restart
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ 服务重启成功！"
        else
            echo ""
            echo "❌ 重启失败！请检查错误信息"
        fi
        ;;
    
    4)
        echo ""
        echo "📋 服务状态："
        docker ps --filter "name=gva-"
        ;;
    
    5)
        echo ""
        echo "📝 MySQL 日志 (按 Ctrl+C 退出)："
        echo ""
        docker logs -f gva-mysql
        ;;
    
    6)
        echo ""
        echo "👋 再见！"
        exit 0
        ;;
    
    *)
        echo ""
        echo "❌ 无效的选项！"
        exit 1
        ;;
esac

echo ""
echo "完成！"
