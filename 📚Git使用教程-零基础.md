# 📚 Git 使用教程 - 零基础入门

> **给完全不懂 Git 的你**  
> 10分钟学会日常使用，保证看得懂、用得上！

---

## 🎯 Git 是什么？

想象一下：
- 你在写论文，每次修改都另存为 `论文v1.doc`、`论文v2.doc`、`论文最终版.doc`、`论文真正的最终版.doc`...
- **Git 就是帮你自动管理这些版本的工具**，不用自己起名字，还能随时回到任何一个历史版本！

---

## 📖 目录

1. [基础概念（5个核心词）](#基础概念)
2. [日常使用（3个命令搞定）](#日常使用)
3. [查看历史（看你改了啥）](#查看历史)
4. [撤销修改（后悔药）](#撤销修改)
5. [备份到云端（GitHub）](#备份到云端)
6. [常见问题](#常见问题)

---

## 1️⃣ 基础概念（5个核心词）

### 📦 **仓库（Repository）**
就是你的项目文件夹，里面有个隐藏的 `.git` 文件夹记录所有历史版本。

```
study-hub-main/          ← 这就是你的仓库
├── .git/                ← Git 的历史记录（隐藏文件夹）
├── server/              ← 你的代码文件
├── web/
└── ...
```

### 📝 **工作区（Working Directory）**
就是你正在编辑的文件，你看到的、能修改的文件。

### 📋 **暂存区（Staging Area）**
准备提交的文件"候车室"，用 `git add` 把文件放进去。

### 💾 **提交（Commit）**
给当前版本拍个快照，并写个说明"我改了啥"。

### 🌿 **分支（Branch）**
平行宇宙！可以在不同分支上尝试不同的改动，互不影响。

---

## 2️⃣ 日常使用（3个命令搞定）

### 🔄 **完整工作流程**

```bash
# 1️⃣ 修改代码
# （在 VSCode 或其他编辑器里改代码）

# 2️⃣ 查看改了啥
git status              # 看哪些文件被修改了
git diff                # 看具体改了哪些行

# 3️⃣ 添加到暂存区
git add .               # 添加所有修改的文件
# 或者
git add 文件路径         # 只添加某个文件

# 4️⃣ 提交到本地仓库
git commit -m "说明：我改了什么"

# ✅ 完成！你的修改已经被记录了
```

### 📝 **举例说明**

假设你修改了 `eduEnrollment.vue` 文件，添加了一个按钮：

```bash
# 步骤1：查看状态
$ git status
On branch main
Changes not staged for commit:
  modified:   web/src/view/eduEnrollment/eduEnrollment.vue

# 步骤2：添加到暂存区
$ git add web/src/view/eduEnrollment/eduEnrollment.vue
# 或者添加所有修改
$ git add .

# 步骤3：提交
$ git commit -m "feat: 在报名管理页面添加导出按钮"

# 步骤4：查看提交历史
$ git log --oneline
a1b2c3d feat: 在报名管理页面添加导出按钮  ← 新的提交
47ea599 feat: StudyHub 培训班课时管理系统 - 完整版本
71c100c feat: 初始提交 - 基于 StudyHub 开源项目...
```

---

## 3️⃣ 查看历史（看你改了啥）

### 📊 **查看提交历史**

```bash
# 简洁版（一行一个提交）
git log --oneline

# 详细版（包含作者、时间）
git log

# 图形化版本（看分支走向）
git log --oneline --graph --all

# 查看最近 5 次提交
git log --oneline -5
```

### 🔍 **查看具体改了什么**

```bash
# 查看某次提交的详细内容
git show 71c100c        # 71c100c 是提交的 ID（前7位）

# 查看某个文件的修改历史
git log -- web/src/view/eduEnrollment/eduEnrollment.vue

# 查看两次提交之间的差异
git diff 71c100c 47ea599

# 查看当前未提交的修改
git diff                # 工作区 vs 暂存区
git diff --staged       # 暂存区 vs 最后一次提交
```

### 📈 **查看文件修改统计**

```bash
# 查看每次提交修改了多少行
git log --stat

# 查看某次提交的统计
git show --stat 71c100c

# 查看哪些文件被修改过
git log --name-only
```

---

## 4️⃣ 撤销修改（后悔药）

### 🔙 **场景1：还没 add（工作区修改）**

```bash
# 丢弃某个文件的修改（回到最后一次提交的状态）
git checkout -- 文件路径

# 丢弃所有修改
git checkout -- .

# 举例
git checkout -- web/src/view/eduEnrollment/eduEnrollment.vue
```

### 🔙 **场景2：已经 add 但还没 commit（暂存区）**

```bash
# 从暂存区移除（但保留工作区的修改）
git reset HEAD 文件路径

# 举例
git reset HEAD web/src/view/eduEnrollment/eduEnrollment.vue

# 然后如果想丢弃修改
git checkout -- web/src/view/eduEnrollment/eduEnrollment.vue
```

### 🔙 **场景3：已经 commit 了（本地仓库）**

```bash
# 方法1：回退到上一个版本（保留修改）
git reset --soft HEAD~1     # 回退1个版本，修改保留在暂存区

# 方法2：回退到上一个版本（丢弃修改）
git reset --hard HEAD~1     # 回退1个版本，修改完全丢弃

# 方法3：回退到指定版本
git reset --hard 71c100c    # 回退到指定的提交

# 方法4：创建一个新提交来撤销（推荐）
git revert HEAD             # 撤销最后一次提交，生成新提交
```

### ⚠️ **重要提示**

```bash
# ⚠️ 危险操作（会丢失数据）
git reset --hard            # 慎用！会丢失所有未提交的修改

# ✅ 安全操作（可以恢复）
git reset --soft            # 安全！修改还在
git revert                  # 安全！创建新提交
```

---

## 5️⃣ 备份到云端（GitHub）

### 🌐 **为什么要推送到 GitHub？**

- ✅ 防止电脑坏了、硬盘挂了
- ✅ 多台电脑同步代码
- ✅ 团队协作
- ✅ 开源分享

### 📤 **如何推送到 GitHub？**

#### **第1步：在 GitHub 上创建仓库**

1. 打开 https://github.com
2. 登录（没有账号就注册一个）
3. 点击右上角 `+` → `New repository`
4. 填写：
   - Repository name: `study-hub`
   - 勾选 `Private`（私有仓库，别人看不到）
5. 点击 `Create repository`

#### **第2步：连接本地仓库和 GitHub**

```bash
# 进入你的项目目录
cd /Users/yandrummer/Desktop/study-hub-main

# 添加远程仓库（只需要做一次）
git remote add origin https://github.com/你的用户名/study-hub.git

# 查看远程仓库
git remote -v
```

#### **第3步：推送代码**

```bash
# 第一次推送（需要设置上游分支）
git push -u origin main

# 以后每次推送只需要
git push
```

#### **第4步：从 GitHub 拉取代码（其他电脑）**

```bash
# 克隆仓库到本地
git clone https://github.com/你的用户名/study-hub.git

# 进入目录
cd study-hub

# 拉取最新代码
git pull
```

---

## 6️⃣ 常见问题

### ❓ **问题1：提交信息写错了怎么办？**

```bash
# 修改最后一次提交的信息
git commit --amend -m "新的提交信息"
```

### ❓ **问题2：忘记某个文件没加进去**

```bash
# 添加遗漏的文件
git add 遗漏的文件

# 追加到上一次提交（不产生新提交）
git commit --amend --no-edit
```

### ❓ **问题3：查看 .git 文件夹（macOS）**

```bash
# 在 Finder 中显示隐藏文件
Command + Shift + .

# 在终端中查看
ls -la
```

### ❓ **问题4：误删了文件怎么恢复？**

```bash
# 从最后一次提交恢复
git checkout -- 文件路径

# 或者
git restore 文件路径
```

### ❓ **问题5：想看某个文件的历史版本**

```bash
# 查看文件的修改历史
git log -- 文件路径

# 恢复到某个版本
git checkout 71c100c -- 文件路径
```

### ❓ **问题6：Git 太难了，有图形界面吗？**

有！推荐工具：
- **VSCode 自带**：左侧第3个图标（Source Control）
- **GitHub Desktop**：https://desktop.github.com/
- **Sourcetree**：https://www.sourcetreeapp.com/

---

## 🎯 每天的工作流程（记住这个就够了）

```bash
# 早上开始工作
cd /Users/yandrummer/Desktop/study-hub-main
git status              # 看看昨天有没有未提交的

# 开始写代码...
# ...

# 中午/下午休息前
git status              # 看看改了啥
git add .               # 添加所有修改
git commit -m "feat: 实现了XXX功能"
git push                # 推送到 GitHub（如果配置了）

# 晚上下班前
git status              # 再检查一次
git add .
git commit -m "fix: 修复了XXX问题"
git push

# ✅ 安心下班！代码已经安全保存
```

---

## 📝 提交信息规范（可选）

写好提交信息，以后查找方便：

```bash
# 格式：类型: 简短描述

feat: 添加新功能          # 新功能
fix: 修复bug             # 修复问题
docs: 更新文档           # 只改文档
style: 代码格式调整       # 不影响功能的改动
refactor: 重构代码       # 既不是新功能也不是修复
test: 添加测试           # 测试相关
chore: 杂项             # 其他改动

# 举例
git commit -m "feat: 在报名管理页面添加导出Excel功能"
git commit -m "fix: 修复消课时剩余课时计算错误"
git commit -m "docs: 更新README使用说明"
```

---

## 🎓 进阶学习（可选）

当你熟悉基础命令后，可以学习：

### 🌿 **分支操作**

```bash
# 创建分支
git branch 分支名

# 切换分支
git checkout 分支名

# 创建并切换（合并上面两步）
git checkout -b 分支名

# 合并分支
git merge 分支名

# 删除分支
git branch -d 分支名
```

### 🔍 **高级查询**

```bash
# 搜索提交信息
git log --grep="关键词"

# 搜索代码内容
git log -S "代码片段"

# 查看某个人的提交
git log --author="yandrummer"

# 查看最近一周的提交
git log --since="1 week ago"
```

---

## 🚨 重要提醒

### ✅ **好习惯**

1. **经常提交**：一个小功能完成就提交，不要攒一天
2. **写清楚提交信息**：以后查找方便
3. **每天推送**：推送到 GitHub 备份
4. **提交前看一眼**：`git status` 和 `git diff` 确认没问题

### ❌ **坏习惯**

1. ❌ 不写提交信息或乱写：`git commit -m "update"`
2. ❌ 长时间不提交：改了一周才提交
3. ❌ 提交敏感信息：密码、密钥等
4. ❌ 直接 `git reset --hard`：容易丢失代码

---

## 🎁 实战练习

现在就试试！

```bash
# 1. 修改一个文件（比如 README.md）
# 2. 查看状态
git status

# 3. 查看具体改了什么
git diff

# 4. 添加并提交
git add README.md
git commit -m "docs: 更新项目说明"

# 5. 查看历史
git log --oneline

# ✅ 你已经会用 Git 了！
```

---

## 📚 参考资料

- 🌐 Git 官方文档（中文）：https://git-scm.com/book/zh/v2
- 📖 Git 可视化学习：https://learngitbranching.js.org/?locale=zh_CN
- 🎮 Git 游戏教程：https://ohmygit.org/

---

## 💬 遇到问题？

如果遇到任何问题：
1. 先用 `git status` 看看当前状态
2. Google 搜索错误信息
3. 问 AI 助手（比如我！）

---

**🎉 恭喜！你已经学会 Git 的日常使用了！**

记住：**Git 是用来帮你的，不要怕它。大胆尝试，99% 的操作都可以撤销！**
