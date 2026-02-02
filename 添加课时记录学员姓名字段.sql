-- 为课时记录表添加学员姓名字段
-- 日期: 2025-02-01
-- 说明: 在课时记录中添加学员姓名，方便查询和显示

-- 1. 添加学员姓名字段
ALTER TABLE `edu_class_session` 
ADD COLUMN `user_name` VARCHAR(255) NULL COMMENT '学员姓名' AFTER `course_name`;

-- 2. 为已存在的记录填充学员姓名（通过关联查询）
UPDATE `edu_class_session` ecs
INNER JOIN `edu_enrollment` ee ON ecs.enrollment_id = ee.id
INNER JOIN `sys_users` su ON ee.user_id = su.id
SET ecs.user_name = su.nick_name
WHERE ecs.user_name IS NULL OR ecs.user_name = '';

-- 3. 验证更新结果
SELECT 
    COUNT(*) as total_records,
    COUNT(user_name) as has_user_name,
    COUNT(*) - COUNT(user_name) as missing_user_name
FROM `edu_class_session`;
