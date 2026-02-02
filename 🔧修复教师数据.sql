-- ============================================
-- 修复教师用户数据脚本
-- 执行前请先运行 🔍诊断教师消失问题.sql 了解具体情况
-- ============================================

-- 步骤1：确保教师角色(9001)存在
-- 如果角色被误删，先恢复角色
INSERT IGNORE INTO sys_authorities (authority_id, authority_name, parent_id, default_router, created_at, updated_at, deleted_at)
VALUES (9001, '教师', 888, 'dashboard', NOW(), NOW(), NULL);

-- 步骤2：确保学员角色(9002)存在
INSERT IGNORE INTO sys_authorities (authority_id, authority_name, parent_id, default_router, created_at, updated_at, deleted_at)
VALUES (9002, '学员', 888, 'dashboard', NOW(), NOW(), NULL);

-- 步骤3：确保教务管理员角色(9003)存在且parent_id正确
INSERT INTO sys_authorities (authority_id, authority_name, parent_id, default_router, created_at, updated_at, deleted_at)
VALUES (9003, '教务管理员', 888, 'dashboard', NOW(), NOW(), NULL)
ON DUPLICATE KEY UPDATE 
    parent_id = 888,
    updated_at = NOW();

-- ============================================
-- 步骤4：修复用户角色关联
-- ============================================

-- 4.1 为默认角色是9001但在sys_user_authority中缺失的用户添加教师角色
SELECT '========== 准备为缺失教师角色的用户添加关联 ==========' AS step;
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.authority_id AS default_authority,
    '将添加角色9001(教师)' AS action
FROM sys_users u
WHERE u.authority_id = 9001  -- 默认角色是教师
  AND NOT EXISTS (
      SELECT 1 FROM sys_user_authority ua 
      WHERE ua.sys_user_id = u.id 
      AND ua.sys_authority_authority_id = 9001
  )
  AND u.deleted_at IS NULL;  -- 排除已删除的用户

-- 执行添加
INSERT INTO sys_user_authority (sys_user_id, sys_authority_authority_id)
SELECT u.id, 9001
FROM sys_users u
WHERE u.authority_id = 9001
  AND NOT EXISTS (
      SELECT 1 FROM sys_user_authority ua 
      WHERE ua.sys_user_id = u.id 
      AND ua.sys_authority_authority_id = 9001
  )
  AND u.deleted_at IS NULL;

SELECT CONCAT('已为 ', ROW_COUNT(), ' 个用户添加教师角色') AS result;

-- 4.2 为默认角色是9002但在sys_user_authority中缺失的用户添加学员角色
SELECT '========== 准备为缺失学员角色的用户添加关联 ==========' AS step;
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.authority_id AS default_authority,
    '将添加角色9002(学员)' AS action
FROM sys_users u
WHERE u.authority_id = 9002
  AND NOT EXISTS (
      SELECT 1 FROM sys_user_authority ua 
      WHERE ua.sys_user_id = u.id 
      AND ua.sys_authority_authority_id = 9002
  )
  AND u.deleted_at IS NULL;

-- 执行添加
INSERT INTO sys_user_authority (sys_user_id, sys_authority_authority_id)
SELECT u.id, 9002
FROM sys_users u
WHERE u.authority_id = 9002
  AND NOT EXISTS (
      SELECT 1 FROM sys_user_authority ua 
      WHERE ua.sys_user_id = u.id 
      AND ua.sys_authority_authority_id = 9002
  )
  AND u.deleted_at IS NULL;

SELECT CONCAT('已为 ', ROW_COUNT(), ' 个用户添加学员角色') AS result;

-- 4.3 为默认角色是9003但在sys_user_authority中缺失的用户添加教务管理员角色
INSERT INTO sys_user_authority (sys_user_id, sys_authority_authority_id)
SELECT u.id, 9003
FROM sys_users u
WHERE u.authority_id = 9003
  AND NOT EXISTS (
      SELECT 1 FROM sys_user_authority ua 
      WHERE ua.sys_user_id = u.id 
      AND ua.sys_authority_authority_id = 9003
  )
  AND u.deleted_at IS NULL;

SELECT CONCAT('已为 ', ROW_COUNT(), ' 个用户添加教务管理员角色') AS result;

-- ============================================
-- 步骤5：验证修复结果
-- ============================================

SELECT '========== 修复后的用户角色统计 ==========' AS step;
SELECT 
    a.authority_id,
    a.authority_name,
    COUNT(DISTINCT ua.sys_user_id) AS user_count
FROM sys_authorities a
LEFT JOIN sys_user_authority ua ON a.authority_id = ua.sys_authority_authority_id
WHERE a.authority_id IN (888, 9001, 9002, 9003)
GROUP BY a.authority_id, a.authority_name
ORDER BY a.authority_id;

SELECT '========== 仍然缺少角色关联的用户（需要手动处理）==========' AS step;
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.phone,
    u.authority_id AS default_authority_id,
    a.authority_name AS default_authority_name
FROM sys_users u
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
LEFT JOIN sys_authorities a ON u.authority_id = a.authority_id
WHERE ua.sys_authority_authority_id IS NULL
  AND u.deleted_at IS NULL
  AND u.authority_id NOT IN (888, 8881)  -- 排除管理员
ORDER BY u.id;

-- ============================================
-- 完成
-- ============================================
SELECT '========== 修复完成 ==========' AS step;
SELECT '请刷新前端页面查看结果' AS message;
