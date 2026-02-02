#!/bin/bash

echo "🔄 StudyHub Docker 数据库重置工具"
echo "================================"
echo ""

# 容器配置
CONTAINER_NAME="gva-mysql"
DB_NAME="studyhub"

echo "⚠️  警告：此操作将删除并重建数据库！"
echo "容器名：$CONTAINER_NAME"
echo "数据库名：$DB_NAME"
echo ""
read -p "确定要继续吗？(输入 yes 继续): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ 操作已取消"
    exit 0
fi

echo ""
echo "📋 检查容器状态..."

# 检查容器是否运行
if ! docker ps | grep -q $CONTAINER_NAME; then
    echo "❌ 容器 $CONTAINER_NAME 未运行！"
    echo ""
    echo "请先启动容器："
    echo "cd ~/Desktop/study-hub-main"
    echo "docker compose -f deploy/docker-compose/docker-compose.yaml up -d"
    echo ""
    echo "或者："
    echo "cd ~/Desktop/study-hub-main/deploy/docker-compose"
    echo "docker compose up -d"
    exit 1
fi

echo "✅ 容器正在运行"
echo ""
echo "🗑️  正在重置数据库..."

# 通过 Docker exec 执行 MySQL 命令
docker exec -i $CONTAINER_NAME mysql -u root -p"Aa@6447985" <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
CREATE DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON $DB_NAME.* TO 'root'@'%';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO 'gva'@'%';
FLUSH PRIVILEGES;
EOF

# 检查是否成功
if [ $? -eq 0 ]; then
    echo "✅ 数据库重置成功！"
    echo ""
    echo "📋 下一步操作："
    echo "1. 重启后端服务："
    echo "   cd /Users/yandrummer/Desktop/study-hub-main/server"
    echo "   go run main.go"
    echo ""
    echo "2. 访问初始化页面："
    echo "   http://localhost:8080/#/init"
    echo ""
    echo "3. 填写数据库信息："
    echo "   - Host: 127.0.0.1"
    echo "   - Port: 3306"
    echo "   - Database: $DB_NAME"
    echo "   - Username: root"
    echo "   - Password: sYx199212240029"
    echo ""
else
    echo "❌ 数据库重置失败！"
    echo ""
    echo "请检查："
    echo "1. 容器是否正常运行：docker ps"
    echo "2. 手动执行命令测试："
    echo "   docker exec -it gva-mysql mysql -u root -pAa@6447985"
    exit 1
fi
