-- ============================================
-- 诊断教师用户消失问题
-- ============================================

-- 1. 检查角色表是否完整
SELECT '========== 1. 检查角色表 ==========' AS step;
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
ORDER BY parent_id, authority_id;

-- 2. 检查是否存在教师角色 (9001)
SELECT '========== 2. 检查教师角色9001是否存在 ==========' AS step;
SELECT authority_id, authority_name, parent_id 
FROM sys_authorities 
WHERE authority_id = 9001;

-- 3. 查找所有用户及其角色关联
SELECT '========== 3. 所有用户及其角色 ==========' AS step;
SELECT 
    u.id AS user_id,
    u.user_name,
    u.nick_name,
    u.authority_id AS default_authority_id,
    GROUP_CONCAT(DISTINCT ua.sys_authority_authority_id) AS all_authority_ids,
    GROUP_CONCAT(DISTINCT a.authority_name) AS all_authority_names
FROM sys_users u
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
LEFT JOIN sys_authorities a ON ua.sys_authority_authority_id = a.authority_id
GROUP BY u.id, u.user_name, u.nick_name, u.authority_id
ORDER BY u.id;

-- 4. 查找没有任何角色关联的用户
SELECT '========== 4. 没有角色关联的用户（孤儿用户）==========' AS step;
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.phone,
    u.authority_id AS default_authority_id
FROM sys_users u
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
WHERE ua.sys_authority_authority_id IS NULL
ORDER BY u.id;

-- 5. 检查sys_user_authority表中的数据
SELECT '========== 5. 用户-角色关联表 ==========' AS step;
SELECT 
    ua.*,
    u.nick_name,
    a.authority_name
FROM sys_user_authority ua
LEFT JOIN sys_users u ON ua.sys_user_id = u.id
LEFT JOIN sys_authorities a ON ua.sys_authority_authority_id = a.authority_id
ORDER BY ua.sys_user_id;

-- 6. 统计每个角色的用户数量
SELECT '========== 6. 每个角色的用户数量 ==========' AS step;
SELECT 
    a.authority_id,
    a.authority_name,
    a.parent_id,
    COUNT(DISTINCT ua.sys_user_id) AS user_count
FROM sys_authorities a
LEFT JOIN sys_user_authority ua ON a.authority_id = ua.sys_authority_authority_id
GROUP BY a.authority_id, a.authority_name, a.parent_id
ORDER BY a.authority_id;

-- 7. 查找应该是教师但缺少9001角色的用户
SELECT '========== 7. 疑似教师但缺少9001角色的用户 ==========' AS step;
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.phone,
    u.authority_id AS default_authority_id,
    GROUP_CONCAT(ua.sys_authority_authority_id) AS current_roles
FROM sys_users u
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
WHERE u.authority_id != 888 
  AND u.authority_id != 8881
  AND u.authority_id != 9002  -- 不是学员
GROUP BY u.id, u.user_name, u.nick_name, u.phone, u.authority_id
HAVING FIND_IN_SET('9001', current_roles) = 0 OR current_roles IS NULL;

-- ============================================
-- 修复建议（根据上面的诊断结果执行）
-- ============================================

-- 修复方案1：如果角色9001被误删，需要重新创建
-- INSERT INTO sys_authorities (authority_id, authority_name, parent_id, default_router, created_at, updated_at)
-- VALUES (9001, '教师', 888, 'dashboard', NOW(), NOW());

-- 修复方案2：如果用户失去了9001角色关联，需要重新建立
-- 示例：为用户ID=X添加教师角色
-- INSERT INTO sys_user_authority (sys_user_id, sys_authority_authority_id)
-- SELECT u.id, 9001
-- FROM sys_users u
-- WHERE u.id = X  -- 替换为具体的用户ID
-- AND NOT EXISTS (
--     SELECT 1 FROM sys_user_authority ua 
--     WHERE ua.sys_user_id = u.id AND ua.sys_authority_authority_id = 9001
-- );

-- 修复方案3：批量修复所有应该是教师但缺少9001角色的用户
-- 警告：执行前请仔细检查诊断结果！
-- INSERT INTO sys_user_authority (sys_user_id, sys_authority_authority_id)
-- SELECT DISTINCT u.id, 9001
-- FROM sys_users u
-- LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id AND ua.sys_authority_authority_id = 9001
-- WHERE u.authority_id = 9001  -- sys_users表中默认角色是9001
--   AND ua.sys_authority_authority_id IS NULL;  -- 但在关联表中没有9001角色
