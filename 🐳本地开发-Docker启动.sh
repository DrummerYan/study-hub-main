#!/bin/bash

echo "🐳 StudyHub 本地开发 - Docker 启动工具"
echo "========================================"
echo ""
echo "说明：此脚本只启动 MySQL 和 Redis 容器"
echo "      前后端请直接在本地运行"
echo ""

# 配置文件路径
COMPOSE_FILE="deploy/docker-compose/docker-compose-local.yaml"
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
echo "1) 启动 MySQL + Redis"
echo "2) 停止所有容器"
echo "3) 重启所有容器"
echo "4) 查看容器状态"
echo "5) 查看 MySQL 日志"
echo "6) 查看 Redis 日志"
echo "7) 进入 MySQL 终端"
echo "8) 清理所有容器和数据（危险！）"
echo "9) 退出"
echo ""

read -p "请输入选项 (1-9): " choice

case $choice in
    1)
        echo ""
        echo "🚀 正在启动 MySQL 和 Redis..."
        docker compose -f "$COMPOSE_FILE" up -d
        
        if [ $? -eq 0 ]; then
            echo ""
            echo "✅ 服务启动成功！"
            echo ""
            echo "📋 运行中的容器："
            docker ps --filter "name=gva-"
            echo ""
            echo "📝 连接信息："
            echo "MySQL:"
            echo "  - Host: 127.0.0.1"
            echo "  - Port: 13306"
            echo "  - User: root"
            echo "  - Password: Aa@6447985"
            echo "  - Database: qmPlus (或稍后创建 studyhub)"
            echo ""
            echo "Redis:"
            echo "  - Host: 127.0.0.1"
            echo "  - Port: 16379"
            echo ""
            echo "💡 下一步："
            echo "1. 运行重置数据库脚本: ./🔄Docker方式重置数据库.sh"
            echo "2. 启动后端: cd server && go run main.go"
            echo "3. 启动前端: cd web && npm run serve"
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
        echo "📋 容器状态："
        docker ps --filter "name=gva-"
        echo ""
        echo "💾 数据卷："
        docker volume ls | grep study-hub
        ;;
    
    5)
        echo ""
        echo "📝 MySQL 日志 (按 Ctrl+C 退出)："
        echo ""
        docker logs -f gva-mysql
        ;;
    
    6)
        echo ""
        echo "📝 Redis 日志 (按 Ctrl+C 退出)："
        echo ""
        docker logs -f gva-redis
        ;;
    
    7)
        echo ""
        echo "🔧 进入 MySQL 终端..."
        echo "提示：输入密码 Aa@6447985"
        echo ""
        docker exec -it gva-mysql mysql -u root -pAa@6447985
        ;;
    
    8)
        echo ""
        echo "⚠️  警告：此操作将删除所有容器和数据！"
        read -p "确定要继续吗？(输入 YES 继续): " confirm
        
        if [ "$confirm" = "YES" ]; then
            echo ""
            echo "🗑️  正在清理..."
            docker compose -f "$COMPOSE_FILE" down -v
            echo ""
            echo "✅ 清理完成"
        else
            echo "❌ 操作已取消"
        fi
        ;;
    
    9)
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
