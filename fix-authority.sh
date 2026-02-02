#!/bin/bash

echo "===================================="
echo "修复教务管理员角色的循环引用问题"
echo "===================================="

# 进入 Docker MySQL 容器并执行修复
docker exec -i study-hub-mysql mysql -uroot -pYandrum@2024 gva << EOF

-- 查看当前问题
SELECT '当前问题数据：' AS info;
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
WHERE authority_id = 9003;

-- 执行修复
UPDATE sys_authorities 
SET parent_id = 888 
WHERE authority_id = 9003;

-- 验证修复结果
SELECT '修复后的数据：' AS info;
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
WHERE authority_id = 9003;

-- 显示完整角色树
SELECT '完整角色树：' AS info;
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
ORDER BY parent_id, authority_id;

EOF

echo ""
echo "===================================="
echo "修复完成！请刷新浏览器页面。"
echo "===================================="
