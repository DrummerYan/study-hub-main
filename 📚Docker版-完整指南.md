# 📚 Docker 版 MySQL 操作指南

## 🎯 快速操作

### **方法 1：一键脚本（推荐！）** ⭐⭐⭐⭐⭐

```bash
cd ~/Desktop/study-hub-main
./🔄Docker方式重置数据库.sh
```

输入 `yes` 确认，等待完成！

---

### **方法 2：手动 Docker 命令**

```bash
# 1. 确保容器正在运行
docker ps | grep gva-mysql

# 2. 进入 MySQL 容器
docker exec -it gva-mysql mysql -u root -pAa@6447985

# 3. 在 MySQL 中执行
DROP DATABASE IF EXISTS studyhub;
CREATE DATABASE studyhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON studyhub.* TO 'root'@'%';
FLUSH PRIVILEGES;
exit;
```

---

## 🔧 完整流程

### **第 1 步：启动 Docker 容器**

```bash
cd ~/Desktop/study-hub-main
docker compose -f deploy/docker-compose/docker-compose.yaml up -d
```

**或者：**
```bash
cd ~/Desktop/study-hub-main/deploy/docker-compose
docker compose up -d
```

**验证容器状态：**
```bash
docker ps
```

应该看到：
- `gva-mysql`
- `gva-redis`

---

### **第 2 步：重置数据库**

运行脚本：
```bash
./🔄Docker方式重置数据库.sh
```

或手动执行：
```bash
docker exec -i gva-mysql mysql -u root -pAa@6447985 <<EOF
DROP DATABASE IF EXISTS studyhub;
CREATE DATABASE studyhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF
```

---

### **第 3 步：重启后端**

```bash
cd /Users/yandrummer/Desktop/study-hub-main/server
go run main.go
```

---

### **第 4 步：访问初始化页面**

浏览器打开：
```
http://localhost:8080/#/init
```

填写信息：
- **Host**: 127.0.0.1
- **Port**: 3306
- **Database**: studyhub
- **Username**: root
- **Password**: sYx199212240029

---

## 💡 常用 Docker 命令

### **查看容器状态**
```bash
docker ps                    # 查看运行中的容器
docker ps -a                 # 查看所有容器
docker logs gva-mysql        # 查看 MySQL 日志
docker logs -f gva-mysql     # 实时查看日志
```

### **启动/停止容器**
```bash
# 在 study-hub-main 根目录下
docker compose -f deploy/docker-compose/docker-compose.yaml up -d         # 启动所有服务
docker compose -f deploy/docker-compose/docker-compose.yaml down          # 停止并删除容器
docker compose -f deploy/docker-compose/docker-compose.yaml stop          # 停止容器（不删除）
docker compose -f deploy/docker-compose/docker-compose.yaml start         # 启动已存在的容器
docker compose -f deploy/docker-compose/docker-compose.yaml restart       # 重启容器

# 或者进入配置文件目录
cd ~/Desktop/study-hub-main/deploy/docker-compose
docker compose up -d
docker compose down
docker compose stop
docker compose start
docker compose restart
```

### **进入 MySQL 容器**
```bash
# 方式 1：直接执行 MySQL 命令
docker exec -it gva-mysql mysql -u root -pAa@6447985

# 方式 2：进入容器 shell
docker exec -it gva-mysql bash
# 然后再执行 mysql -u root -p
```

### **查看数据库数据**
```bash
# 查看所有数据库
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "SHOW DATABASES;"

# 查看 studyhub 数据库的表
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "USE studyhub; SHOW TABLES;"

# 查看某个表的数据
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "USE studyhub; SELECT * FROM sys_base_menus LIMIT 10;"
```

---

## 🐛 常见问题

### **Q1: 容器未运行**
```bash
# 启动容器
cd ~/Desktop/study-hub-main
docker compose up -d

# 查看启动日志
docker logs gva-mysql
```

### **Q2: 端口冲突**
如果 13306 或 3306 端口被占用：
```bash
# 查看占用端口的进程
lsof -i:3306
lsof -i:13306

# 停止冲突的进程或修改 docker-compose.yaml
```

### **Q3: 密码错误**
Docker 容器中 MySQL 的 root 密码是：`Aa@6447985`

如果忘记密码，可以重建容器：
```bash
docker compose down -v  # 删除容器和卷
docker compose up -d    # 重新创建
```

### **Q4: 数据持久化**
Docker Compose 使用 volume 保存数据：
```bash
# 查看卷
docker volume ls

# 删除卷（会丢失所有数据！）
docker volume rm study-hub-main_mysql
```

---

## 🎯 Docker vs 本地 MySQL

### **Docker 方式（你当前使用的）**

**优点：**
- ✅ 环境隔离，不污染系统
- ✅ 一键启动/停止
- ✅ 版本管理方便
- ✅ 可以同时运行多个 MySQL 版本

**缺点：**
- ⚠️ 需要 Docker 环境
- ⚠️ 占用更多资源

### **本地 MySQL**

**优点：**
- ✅ 性能稍好
- ✅ 不需要 Docker

**缺点：**
- ⚠️ 可能与系统其他应用冲突
- ⚠️ 卸载不彻底

---

## 📋 你的配置摘要

### **Docker Compose 配置**
- **容器名**: gva-mysql
- **镜像**: mysql:8.0.21
- **端口**: 13306:3306（主机:容器）
- **Root 密码**: Aa@6447985
- **初始数据库**: qmPlus

### **应用配置 (config.yaml)**
- **Host**: 127.0.0.1
- **Port**: 3306
- **Database**: studyhub
- **Username**: root
- **Password**: sYx199212240029

### **注意事项**
- 应用连接的是容器内部的 3306 端口
- 容器内部的 root 密码是 Aa@6447985
- 但你的 config.yaml 配置的密码是 sYx199212240029
- **这可能会导致连接失败！**

---

## ⚠️ 重要：配置不匹配问题

你的 Docker MySQL root 密码和 config.yaml 中的密码不一致！

**解决方案 A：修改 Docker 配置**
编辑 `config.yaml`，改为 Docker 的密码：
```yaml
mysql:
  username: root
  password: Aa@6447985  # 改为 Docker 的密码
```

**解决方案 B：修改 MySQL 密码**
在 Docker 容器中修改密码：
```bash
docker exec -it gva-mysql mysql -u root -pAa@6447985 -e "ALTER USER 'root'@'%' IDENTIFIED BY 'sYx199212240029'; FLUSH PRIVILEGES;"
```

---

## 🚀 现在开始

1. **运行脚本重置数据库：**
   ```bash
   cd ~/Desktop/study-hub-main
   ./🔄Docker方式重置数据库.sh
   ```

2. **重启后端**
3. **访问初始化页面**
4. **开始测试！**

---

**祝你成功！** 🎉
