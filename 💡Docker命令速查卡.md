# 💡 Docker 命令速查卡

## 🚀 最常用（记住这几个就够了！）

### **启动服务**
```bash
cd ~/Desktop/study-hub-main
docker compose -f deploy/docker-compose/docker-compose.yaml up -d
```

**或者用快捷脚本：**
```bash
cd ~/Desktop/study-hub-main
./🐳Docker快速启动.sh
# 选择选项 1
```

---

### **停止服务**
```bash
docker compose -f deploy/docker-compose/docker-compose.yaml down
```

---

### **查看状态**
```bash
docker ps
```

---

### **重置数据库**
```bash
./🔄Docker方式重置数据库.sh
```

---

## 📋 完整命令参考

### **在项目根目录（study-hub-main）**

所有命令都需要加 `-f deploy/docker-compose/docker-compose.yaml`：

```bash
# 启动
docker compose -f deploy/docker-compose/docker-compose.yaml up -d

# 停止
docker compose -f deploy/docker-compose/docker-compose.yaml down

# 重启
docker compose -f deploy/docker-compose/docker-compose.yaml restart

# 查看日志
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f

# 查看 MySQL 日志
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f mysql
```

---

### **在配置文件目录（deploy/docker-compose）**

不需要 `-f` 参数：

```bash
cd ~/Desktop/study-hub-main/deploy/docker-compose

# 启动
docker compose up -d

# 停止
docker compose down

# 重启
docker compose restart

# 查看日志
docker compose logs -f
```

---

## 🎯 推荐做法

### **方法 A：使用快捷脚本（最简单）**

```bash
# 启动/停止/重启
./🐳Docker快速启动.sh

# 重置数据库
./🔄Docker方式重置数据库.sh
```

---

### **方法 B：创建别名**

在 `~/.zshrc` 或 `~/.bashrc` 添加：

```bash
alias dcup='cd ~/Desktop/study-hub-main && docker compose -f deploy/docker-compose/docker-compose.yaml up -d'
alias dcdown='cd ~/Desktop/study-hub-main && docker compose -f deploy/docker-compose/docker-compose.yaml down'
alias dcps='docker ps --filter "name=gva-"'
alias dclog='docker compose -f ~/Desktop/study-hub-main/deploy/docker-compose/docker-compose.yaml logs -f'
```

然后：
```bash
source ~/.zshrc  # 或 source ~/.bashrc
```

以后就可以直接：
```bash
dcup      # 启动
dcdown    # 停止
dcps      # 查看状态
dclog     # 查看日志
```

---

### **方法 C：创建软链接**

```bash
cd ~/Desktop/study-hub-main
ln -s deploy/docker-compose/docker-compose.yaml docker-compose.yaml
```

以后就可以：
```bash
cd ~/Desktop/study-hub-main
docker compose up -d
docker compose down
```

---

## 🔍 MySQL 操作

### **进入 MySQL**
```bash
docker exec -it gva-mysql mysql -u root -pAa@6447985
```

### **查看数据库**
```bash
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "SHOW DATABASES;"
```

### **查看表**
```bash
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "USE studyhub; SHOW TABLES;"
```

### **查看菜单数据**
```bash
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "USE studyhub; SELECT id, name, path, title FROM sys_base_menus;"
```

---

## 🎨 常用组合

### **完整重启流程**
```bash
# 1. 停止服务
docker compose -f deploy/docker-compose/docker-compose.yaml down

# 2. 重置数据库
./🔄Docker方式重置数据库.sh

# 3. 启动服务
docker compose -f deploy/docker-compose/docker-compose.yaml up -d

# 4. 启动后端
cd server
go run main.go

# 5. 浏览器访问初始化页面
# http://localhost:8080/#/init
```

---

### **快速检查**
```bash
# 查看所有容器
docker ps -a

# 查看 gva 相关容器
docker ps --filter "name=gva-"

# 查看容器占用资源
docker stats

# 查看容器详细信息
docker inspect gva-mysql
```

---

### **日志查看**
```bash
# 实时查看所有日志
docker compose -f deploy/docker-compose/docker-compose.yaml logs -f

# 只看 MySQL
docker logs -f gva-mysql

# 只看最后 100 行
docker logs --tail 100 gva-mysql

# 查看错误日志
docker logs gva-mysql 2>&1 | grep -i error
```

---

## 🐛 故障排查

### **容器无法启动**
```bash
# 查看详细日志
docker logs gva-mysql

# 查看容器状态
docker ps -a

# 重新构建并启动
docker compose -f deploy/docker-compose/docker-compose.yaml up -d --force-recreate
```

---

### **端口冲突**
```bash
# 查看端口占用
lsof -i:3306
lsof -i:13306

# 停止冲突进程
kill <PID>
```

---

### **清理所有数据**
```bash
# ⚠️ 警告：会删除所有数据！
docker compose -f deploy/docker-compose/docker-compose.yaml down -v

# 重新启动
docker compose -f deploy/docker-compose/docker-compose.yaml up -d
```

---

## 📦 容器信息

| 服务 | 容器名 | 镜像 | 端口映射 | 说明 |
|-----|--------|------|---------|------|
| MySQL | gva-mysql | mysql:8.0.21 | 13306:3306 | 数据库 |
| Redis | gva-redis | redis:6.0.6 | 16379:6379 | 缓存 |

---

## 💡 小贴士

1. **记住一个命令就够了：**
   ```bash
   ./🐳Docker快速启动.sh
   ```

2. **如果嫌路径太长，创建软链接：**
   ```bash
   ln -s deploy/docker-compose/docker-compose.yaml docker-compose.yaml
   ```

3. **使用 Tab 补全：**
   输入 `docker comp` 按 Tab 自动补全

4. **查看帮助：**
   ```bash
   docker compose --help
   docker compose up --help
   ```

---

**保存这个文件，需要的时候翻一翻！** 📖
