# 📚 MySQL 终端操作 - 新手指南

## 🎯 你只需要记住这几个命令

### **方法 1：一键脚本（最简单！）** ⭐⭐⭐⭐⭐

我已经帮你做好了脚本，只需要运行：

```bash
cd /Users/yandrummer/Desktop/study-hub-main
./🔄一键重置数据库.sh
```

**步骤：**
1. 运行脚本
2. 输入 `yes` 确认
3. 等待完成
4. 按照提示重启后端和访问初始化页面

**就这么简单！** 🎉

---

### **方法 2：手动操作（学习版）**

如果你想学习 MySQL 命令，这是基础教程：

#### **1️⃣ 连接 MySQL**
```bash
mysql -u root -p
```
- `-u root`：使用 root 用户
- `-p`：需要密码（回车后输入）

**输入密码：** `sYx199212240029`

**成功标志：** 看到 `mysql>` 提示符

---

#### **2️⃣ 查看所有数据库（可选）**
```sql
SHOW DATABASES;
```

你会看到类似这样的列表：
```
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| studyhub           |
| sys                |
+--------------------+
```

---

#### **3️⃣ 删除数据库**
```sql
DROP DATABASE IF EXISTS studyhub;
```

**解释：**
- `DROP DATABASE`：删除数据库
- `IF EXISTS`：如果存在才删除（避免报错）
- `studyhub`：数据库名

**成功提示：** `Query OK, X rows affected`

---

#### **4️⃣ 创建新数据库**
```sql
CREATE DATABASE studyhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

**解释：**
- `CREATE DATABASE studyhub`：创建名为 studyhub 的数据库
- `CHARACTER SET utf8mb4`：使用 UTF-8 字符集（支持中文和 emoji）
- `COLLATE utf8mb4_unicode_ci`：排序规则（不区分大小写）

**成功提示：** `Query OK, 1 row affected`

---

#### **5️⃣ 验证创建成功（可选）**
```sql
SHOW DATABASES;
```

应该能看到新的 `studyhub` 数据库

---

#### **6️⃣ 退出 MySQL**
```sql
exit;
```

或者按 `Ctrl+D`

---

## 💡 **常用 MySQL 命令速查表**

### **数据库操作**
```sql
-- 查看所有数据库
SHOW DATABASES;

-- 创建数据库
CREATE DATABASE 数据库名;

-- 删除数据库
DROP DATABASE 数据库名;

-- 使用某个数据库
USE 数据库名;
```

### **表操作**
```sql
-- 查看当前数据库的所有表
SHOW TABLES;

-- 查看表结构
DESCRIBE 表名;
-- 或者
DESC 表名;

-- 查看表的创建语句
SHOW CREATE TABLE 表名;
```

### **数据查询**
```sql
-- 查看表中所有数据
SELECT * FROM 表名;

-- 查看前 10 条数据
SELECT * FROM 表名 LIMIT 10;

-- 统计数据条数
SELECT COUNT(*) FROM 表名;
```

### **其他有用命令**
```sql
-- 查看 MySQL 版本
SELECT VERSION();

-- 查看当前用户
SELECT USER();

-- 查看当前数据库
SELECT DATABASE();
```

---

## 🎯 **针对你的项目 - 完整流程**

### **重置数据库并初始化**

```bash
# 1. 连接 MySQL
mysql -u root -p
# 输入密码：sYx199212240029

# 2. 在 MySQL 中执行（复制粘贴这三行）
DROP DATABASE IF EXISTS studyhub;
CREATE DATABASE studyhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;

# 3. 重启后端
cd /Users/yandrummer/Desktop/study-hub-main/server
go run main.go

# 4. 浏览器访问
# http://localhost:8080/#/init
```

---

## ❓ **常见问题**

### **Q1: 报错 "Access denied"**
**原因：** 密码错误或用户没权限

**解决：**
```bash
# 确认密码是否正确
mysql -u root -psYx199212240029
# 注意：-p 后面直接跟密码，中间没有空格
```

### **Q2: 报错 "Can't connect to MySQL server"**
**原因：** MySQL 服务未运行

**解决：**
```bash
# Mac 上启动 MySQL
brew services start mysql
# 或者
sudo /usr/local/mysql/support-files/mysql.server start
```

### **Q3: 输入命令后没反应**
**原因：** 忘记加分号 `;`

**解决：** 
MySQL 命令必须以 `;` 结尾，如果忘记了，再输入一个 `;` 然后回车

### **Q4: 看到 `->` 提示符**
**原因：** 命令没结束（可能是引号没关闭）

**解决：**
- 如果是误操作，输入 `'` 或 `"` 然后 `;` 结束
- 或者按 `Ctrl+C` 取消当前命令

---

## 🎉 **总结**

### **最简单的方法（推荐）：**
```bash
./🔄一键重置数据库.sh
```

### **学习版（手动操作）：**
```bash
mysql -u root -p
# 输入密码
DROP DATABASE IF EXISTS studyhub;
CREATE DATABASE studyhub CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;
```

### **下一步：**
1. 重启后端：`cd server && go run main.go`
2. 访问：`http://localhost:8080/#/init`
3. 初始化数据库
4. 登录查看左侧菜单，应该有"报名管理"了！

---

**选择你喜欢的方法，立即开始吧！** 🚀
