-- 修复教务管理员角色的循环引用问题
-- 将 parent_id 从 9003（自己）改回 888（超级管理员）

-- 查看当前问题数据
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
WHERE authority_id = 9003;

-- 修复：将教务管理员的 parent_id 改回 888
UPDATE sys_authorities 
SET parent_id = 888 
WHERE authority_id = 9003;

-- 验证修复结果
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
WHERE authority_id = 9003;

-- 查看完整的角色树结构
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
ORDER BY parent_id, authority_id;
