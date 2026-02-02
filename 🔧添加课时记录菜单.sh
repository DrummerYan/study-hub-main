#!/bin/bash

echo "🔧 添加课时记录菜单到数据库"
echo "============================="
echo ""

# 执行SQL插入菜单
docker exec -i gva-mysql mysql -u root -pAa@6447985 studyhub <<'EOF'
-- 添加课时记录菜单
INSERT INTO sys_base_menus (created_at, updated_at, deleted_at, menu_level, parent_id, path, name, hidden, component, sort, keep_alive, default_menu, title, icon, close_tab)
VALUES (
  NOW(),
  NOW(),
  NULL,
  0,
  '0',
  'eduClassSession',
  'eduClassSession',
  0,
  'view/eduClassSession/eduClassSession.vue',
  4,
  0,
  0,
  '课时记录',
  'tickets',
  0
);

-- 获取刚插入的菜单ID
SET @menu_id = LAST_INSERT_ID();

-- 为超级管理员（ID=888）授权这个菜单
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT '888', @menu_id
WHERE NOT EXISTS (
  SELECT 1 FROM sys_authority_menus 
  WHERE sys_authority_authority_id = '888' AND sys_base_menu_id = @menu_id
);

-- 验证插入结果
SELECT id, name, path, title FROM sys_base_menus WHERE name = 'eduClassSession';

EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 课时记录菜单添加成功！"
    echo ""
    echo "📋 下一步操作："
    echo "1. 刷新浏览器页面（Cmd+Shift+R 或 Ctrl+Shift+R）"
    echo "2. 点击左侧菜单，应该能看到"课时记录"了"
    echo "3. 再次点击"历史"按钮测试"
    echo ""
else
    echo ""
    echo "❌ 添加失败！"
    exit 1
fi
