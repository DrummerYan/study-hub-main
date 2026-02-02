#!/bin/bash

echo "🔧 StudyHub 快速问题排查工具"
echo "================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查后端端口
echo "1️⃣  检查后端服务（8888 端口）..."
if lsof -Pi :8888 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    PID=$(lsof -ti:8888)
    echo -e "${GREEN}✅ 后端服务正在运行${NC}"
    echo "   进程 ID: $PID"
else
    echo -e "${RED}❌ 后端服务未运行${NC}"
    echo "   解决方法："
    echo "   cd /Users/yandrummer/Desktop/study-hub-main/server"
    echo "   go run main.go"
fi
echo ""

# 检查前端端口
echo "2️⃣  检查前端服务（8080 端口）..."
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    PID=$(lsof -ti:8080)
    echo -e "${GREEN}✅ 前端服务正在运行${NC}"
    echo "   进程 ID: $PID"
else
    echo -e "${RED}❌ 前端服务未运行${NC}"
    echo "   解决方法："
    echo "   cd /Users/yandrummer/Desktop/study-hub-main/web"
    echo "   npm run serve"
fi
echo ""

# 测试后端健康检查
echo "3️⃣  测试后端健康检查..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8888/health 2>/dev/null)
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ 后端健康检查通过${NC}"
    echo "   HTTP 状态码: $HTTP_CODE"
else
    echo -e "${RED}❌ 后端健康检查失败${NC}"
    echo "   HTTP 状态码: $HTTP_CODE"
    echo "   可能原因："
    echo "   - 后端服务未启动"
    echo "   - 端口被占用"
    echo "   - 数据库连接失败"
fi
echo ""

# 测试前端访问
echo "4️⃣  测试前端页面..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080 2>/dev/null)
if [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ 前端页面可访问${NC}"
    echo "   HTTP 状态码: $HTTP_CODE"
    echo "   浏览器访问: http://localhost:8080"
else
    echo -e "${RED}❌ 前端页面无法访问${NC}"
    echo "   HTTP 状态码: $HTTP_CODE"
    echo "   可能原因："
    echo "   - 前端服务未启动"
    echo "   - npm run serve 失败"
fi
echo ""

# 检查关键文件
echo "5️⃣  检查关键文件..."
FILES=(
    "/Users/yandrummer/Desktop/study-hub-main/server/main.go"
    "/Users/yandrummer/Desktop/study-hub-main/web/src/main.js"
    "/Users/yandrummer/Desktop/study-hub-main/web/src/view/eduEnrollment/eduEnrollment.vue"
    "/Users/yandrummer/Desktop/study-hub-main/web/src/api/eduEnrollment.js"
)

ALL_EXISTS=true
for FILE in "${FILES[@]}"; do
    if [ -f "$FILE" ]; then
        echo -e "${GREEN}✅${NC} $(basename $FILE)"
    else
        echo -e "${RED}❌${NC} $(basename $FILE) 不存在"
        ALL_EXISTS=false
    fi
done
echo ""

# 总结
echo "================================"
echo "📊 排查结果总结"
echo "================================"

if lsof -Pi :8888 -sTCP:LISTEN -t >/dev/null 2>&1 && \
   lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1 && \
   [ "$HTTP_CODE" = "200" ] && \
   [ "$ALL_EXISTS" = true ]; then
    echo -e "${GREEN}🎉 所有检查通过！系统运行正常！${NC}"
    echo ""
    echo "你现在可以："
    echo "1. 打开浏览器访问: http://localhost:8080"
    echo "2. 按照 🧪测试指南-跟着做就行.md 进行测试"
else
    echo -e "${YELLOW}⚠️  发现问题，请按照上述提示解决${NC}"
    echo ""
    echo "常见解决方法："
    echo "1. 停止占用端口的进程: kill \$(lsof -ti:8888)"
    echo "2. 重新启动后端: cd server && go run main.go"
    echo "3. 重新启动前端: cd web && npm run serve"
fi
echo ""

# 提供快速修复选项
echo "需要快速修复吗？"
echo "1) 停止所有服务"
echo "2) 重启后端服务"
echo "3) 重启前端服务"
echo "4) 退出"
echo ""
read -p "请选择 (1-4): " choice

case $choice in
    1)
        echo ""
        echo "正在停止所有服务..."
        if lsof -ti:8888 >/dev/null 2>&1; then
            kill $(lsof -ti:8888) 2>/dev/null && echo "✅ 后端服务已停止"
        fi
        if lsof -ti:8080 >/dev/null 2>&1; then
            kill $(lsof -ti:8080) 2>/dev/null && echo "✅ 前端服务已停止"
        fi
        echo "✅ 所有服务已停止"
        ;;
    2)
        echo ""
        echo "正在重启后端服务..."
        if lsof -ti:8888 >/dev/null 2>&1; then
            kill $(lsof -ti:8888) 2>/dev/null
            sleep 2
        fi
        cd /Users/yandrummer/Desktop/study-hub-main/server
        echo "后端服务启动中..."
        go run main.go
        ;;
    3)
        echo ""
        echo "正在重启前端服务..."
        if lsof -ti:8080 >/dev/null 2>&1; then
            kill $(lsof -ti:8080) 2>/dev/null
            sleep 2
        fi
        cd /Users/yandrummer/Desktop/study-hub-main/web
        echo "前端服务启动中..."
        npm run serve
        ;;
    4)
        echo "退出"
        ;;
    *)
        echo "无效选项"
        ;;
esac
