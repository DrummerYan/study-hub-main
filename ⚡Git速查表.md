# ⚡ Git 速查表

> 日常使用只需要记住这几个命令！

---

## 🔥 最常用（90%的时候用这些）

```bash
# 查看状态
git status

# 查看改了什么
git diff

# 添加所有修改
git add .

# 提交
git commit -m "说明"

# 查看历史
git log --oneline
```

---

## 📋 完整工作流

```bash
# 1. 改代码...

# 2. 查看改了啥
git status
git diff

# 3. 添加并提交
git add .
git commit -m "feat: 添加XXX功能"

# 4. 推送到 GitHub（如果有）
git push
```

---

## 🔍 查看类命令

```bash
git status              # 查看当前状态
git diff                # 查看工作区修改
git diff --staged       # 查看暂存区修改
git log                 # 查看提交历史
git log --oneline       # 简洁版历史
git log --stat          # 带统计的历史
git show 提交ID         # 查看某次提交详情
```

---

## 💾 提交类命令

```bash
git add .               # 添加所有修改
git add 文件路径        # 添加某个文件
git commit -m "说明"    # 提交
git commit --amend      # 修改上次提交
```

---

## 🔙 撤销类命令

```bash
# 丢弃工作区修改
git checkout -- 文件路径
git checkout -- .

# 取消暂存
git reset HEAD 文件路径

# 回退版本（保留修改）
git reset --soft HEAD~1

# 回退版本（丢弃修改）⚠️
git reset --hard HEAD~1

# 撤销某次提交（安全）
git revert 提交ID
```

---

## 🌐 远程类命令

```bash
# 添加远程仓库
git remote add origin 仓库地址

# 查看远程仓库
git remote -v

# 推送
git push
git push -u origin main     # 第一次推送

# 拉取
git pull

# 克隆
git clone 仓库地址
```

---

## 🌿 分支类命令

```bash
# 查看分支
git branch

# 创建分支
git branch 分支名

# 切换分支
git checkout 分支名

# 创建并切换
git checkout -b 分支名

# 合并分支
git merge 分支名

# 删除分支
git branch -d 分支名
```

---

## 🎯 一句话记忆

```
改代码 → status看状态 → add添加 → commit提交 → push推送
```

---

## ⚠️ 危险命令（慎用）

```bash
git reset --hard        # 会丢失未提交的修改
git push --force        # 会覆盖远程历史
git clean -fd           # 会删除未跟踪的文件
```

---

## 💡 提交信息模板

```bash
feat: 添加新功能
fix: 修复bug
docs: 更新文档
style: 代码格式
refactor: 重构
test: 测试
chore: 杂项
```

---

**📌 把这个文件收藏起来，忘了就来看！**
