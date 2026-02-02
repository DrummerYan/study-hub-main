-- ============================================
-- 数据库层面的防护措施
-- 防止核心角色被误删除
-- ============================================

-- 1. 创建触发器：防止删除核心角色
DELIMITER //

DROP TRIGGER IF EXISTS before_delete_core_authority//

CREATE TRIGGER before_delete_core_authority
BEFORE DELETE ON sys_authorities
FOR EACH ROW
BEGIN
    -- 禁止删除核心角色 (888超级管理员, 9001教师, 9002学员, 9003教务管理员)
    IF OLD.authority_id IN (888, 9001, 9002, 9003) THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '❌ 禁止删除核心角色！这些角色是系统运行必需的。';
    END IF;
    
    -- 禁止删除有用户使用的角色（额外保护）
    IF EXISTS (
        SELECT 1 FROM sys_user_authority 
        WHERE sys_authority_authority_id = OLD.authority_id
        LIMIT 1
    ) THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '❌ 此角色有用户正在使用，禁止删除！请先将用户改为其他角色。';
    END IF;
    
    -- 禁止删除有子角色的角色
    IF EXISTS (
        SELECT 1 FROM sys_authorities 
        WHERE parent_id = OLD.authority_id
        LIMIT 1
    ) THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '❌ 此角色存在子角色，禁止删除！请先删除所有子角色。';
    END IF;
END//

DELIMITER ;

-- 2. 测试触发器
SELECT '========== 触发器创建成功 ==========' AS status;
SELECT '尝试删除核心角色将被阻止' AS message;

-- 测试示例（会失败，这是正确的）
-- DELETE FROM sys_authorities WHERE authority_id = 9001;

-- ============================================
-- 3. 创建审计日志表（可选）
-- ============================================

CREATE TABLE IF NOT EXISTS `sys_authority_audit_log` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `authority_id` int(11) NOT NULL COMMENT '角色ID',
    `authority_name` varchar(255) NOT NULL COMMENT '角色名称',
    `operation` varchar(50) NOT NULL COMMENT '操作类型：CREATE/UPDATE/DELETE/ATTEMPT_DELETE',
    `operator_id` int(11) DEFAULT NULL COMMENT '操作人ID',
    `operator_name` varchar(255) DEFAULT NULL COMMENT '操作人姓名',
    `details` text DEFAULT NULL COMMENT '操作详情',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    PRIMARY KEY (`id`),
    KEY `idx_authority_id` (`authority_id`),
    KEY `idx_operation` (`operation`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色操作审计日志';

-- 4. 创建审计触发器（可选）
DELIMITER //

DROP TRIGGER IF EXISTS after_delete_authority_audit//

CREATE TRIGGER after_delete_authority_audit
AFTER DELETE ON sys_authorities
FOR EACH ROW
BEGIN
    INSERT INTO sys_authority_audit_log 
    (authority_id, authority_name, operation, details, created_at)
    VALUES 
    (OLD.authority_id, OLD.authority_name, 'DELETE', 
     CONCAT('角色被删除。父角色ID: ', IFNULL(OLD.parent_id, 'NULL')), NOW());
END//

DROP TRIGGER IF EXISTS after_insert_authority_audit//

CREATE TRIGGER after_insert_authority_audit
AFTER INSERT ON sys_authorities
FOR EACH ROW
BEGIN
    INSERT INTO sys_authority_audit_log 
    (authority_id, authority_name, operation, details, created_at)
    VALUES 
    (NEW.authority_id, NEW.authority_name, 'CREATE', 
     CONCAT('新角色创建。父角色ID: ', IFNULL(NEW.parent_id, 'NULL')), NOW());
END//

DROP TRIGGER IF EXISTS after_update_authority_audit//

CREATE TRIGGER after_update_authority_audit
AFTER UPDATE ON sys_authorities
FOR EACH ROW
BEGIN
    INSERT INTO sys_authority_audit_log 
    (authority_id, authority_name, operation, details, created_at)
    VALUES 
    (NEW.authority_id, NEW.authority_name, 'UPDATE', 
     CONCAT('角色更新。旧名称: ', OLD.authority_name, 
            ', 新名称: ', NEW.authority_name,
            ', 旧父ID: ', IFNULL(OLD.parent_id, 'NULL'),
            ', 新父ID: ', IFNULL(NEW.parent_id, 'NULL')), NOW());
END//

DELIMITER ;

-- ============================================
-- 5. 创建数据一致性检查视图
-- ============================================

-- 查看孤儿用户（没有角色的用户）
CREATE OR REPLACE VIEW v_orphan_users AS
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.phone,
    u.email,
    u.authority_id AS default_authority_id,
    a.authority_name AS default_authority_name,
    '缺少角色关联' AS issue_type
FROM sys_users u
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
LEFT JOIN sys_authorities a ON u.authority_id = a.authority_id
WHERE ua.sys_authority_authority_id IS NULL
  AND u.deleted_at IS NULL
  AND u.authority_id NOT IN (888, 8881);

-- 查看角色不一致的用户（默认角色在sys_users中，但在sys_user_authority中没有）
CREATE OR REPLACE VIEW v_inconsistent_user_roles AS
SELECT 
    u.id,
    u.user_name,
    u.nick_name,
    u.authority_id AS default_authority_id,
    a.authority_name AS default_authority_name,
    GROUP_CONCAT(ua.sys_authority_authority_id) AS actual_roles,
    '默认角色与实际角色不匹配' AS issue_type
FROM sys_users u
LEFT JOIN sys_authorities a ON u.authority_id = a.authority_id
LEFT JOIN sys_user_authority ua ON u.id = ua.sys_user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.user_name, u.nick_name, u.authority_id, a.authority_name
HAVING FIND_IN_SET(u.authority_id, actual_roles) = 0 OR actual_roles IS NULL;

-- 查看使用已删除角色的用户
CREATE OR REPLACE VIEW v_users_with_deleted_roles AS
SELECT 
    ua.sys_user_id,
    u.user_name,
    u.nick_name,
    ua.sys_authority_authority_id AS deleted_role_id,
    '使用已删除的角色' AS issue_type
FROM sys_user_authority ua
LEFT JOIN sys_authorities a ON ua.sys_authority_authority_id = a.authority_id
LEFT JOIN sys_users u ON ua.sys_user_id = u.id
WHERE a.authority_id IS NULL;

-- ============================================
-- 6. 定期检查任务（手动执行或通过cron）
-- ============================================

-- 检查并报告所有数据不一致问题
SELECT '========== 数据一致性检查报告 ==========' AS report;

SELECT '1. 孤儿用户（没有角色关联）' AS check_type, COUNT(*) AS issue_count
FROM v_orphan_users;

SELECT * FROM v_orphan_users;

SELECT '2. 角色不一致的用户' AS check_type, COUNT(*) AS issue_count
FROM v_inconsistent_user_roles;

SELECT * FROM v_inconsistent_user_roles;

SELECT '3. 使用已删除角色的用户' AS check_type, COUNT(*) AS issue_count
FROM v_users_with_deleted_roles;

SELECT * FROM v_users_with_deleted_roles;

-- ============================================
-- 使用说明
-- ============================================

/*
1. 执行此脚本创建防护措施：
   mysql -u root -p your_database < 🛡️防护措施-数据库触发器.sql

2. 触发器会自动阻止：
   - 删除核心角色（888, 9001, 9002, 9003）
   - 删除有用户使用的角色
   - 删除有子角色的角色

3. 定期检查数据一致性：
   SELECT * FROM v_orphan_users;
   SELECT * FROM v_inconsistent_user_roles;
   SELECT * FROM v_users_with_deleted_roles;

4. 查看审计日志：
   SELECT * FROM sys_authority_audit_log 
   ORDER BY created_at DESC 
   LIMIT 50;

5. 如需临时禁用触发器（谨慎操作）：
   DROP TRIGGER before_delete_core_authority;
   
6. 重新启用触发器：
   重新执行本脚本即可
*/
