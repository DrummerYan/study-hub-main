#!/bin/bash

echo "🔄 StudyHub 数据库重置工具"
echo "================================"
echo ""

# 数据库配置
DB_USER="root"
DB_PASS="sYx199212240029"
DB_NAME="studyhub"

echo "⚠️  警告：此操作将删除所有数据！"
echo "数据库名：$DB_NAME"
echo ""
read -p "确定要继续吗？(输入 yes 继续): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ 操作已取消"
    exit 0
fi

echo ""
echo "🗑️  正在删除旧数据库..."

# 执行 MySQL 命令
mysql -u $DB_USER -p$DB_PASS <<EOF
DROP DATABASE IF EXISTS $DB_NAME;
CREATE DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
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
    echo "   - Username: $DB_USER"
    echo "   - Password: $DB_PASS"
    echo ""
else
    echo "❌ 数据库重置失败！"
    echo ""
    echo "可能的原因："
    echo "1. MySQL 服务未运行"
    echo "2. 密码不正确"
    echo "3. 权限不足"
    echo ""
    echo "请手动执行："
    echo "mysql -u root -p"
    echo "然后输入以下命令："
    echo "DROP DATABASE IF EXISTS $DB_NAME;"
    echo "CREATE DATABASE $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    exit 1
fi
