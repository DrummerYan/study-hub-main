-- Fix role menu associations and casbin policies for roles 9001/9002/9003
-- Run inside your target database (e.g., `USE gva;`)

START TRANSACTION;

-- Ensure core roles exist
INSERT INTO sys_authorities (authority_id, authority_name, parent_id, default_router, created_at, updated_at)
VALUES
  (9001, '教师', 888, 'dashboard', NOW(), NOW()),
  (9002, '学员', 888, 'dashboard', NOW(), NOW()),
  (9003, '教务管理员', 888, 'dashboard', NOW(), NOW())
ON DUPLICATE KEY UPDATE
  authority_name = VALUES(authority_name),
  parent_id = VALUES(parent_id),
  default_router = VALUES(default_router),
  updated_at = VALUES(updated_at);

-- Reset menu associations for these roles
DELETE FROM sys_authority_menus
WHERE sys_authority_authority_id IN ('9001', '9002', '9003');

-- Teacher (9001)
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT '9001', id
FROM sys_base_menus
WHERE path IN ('dashboard', 'person', 'eduOrganization', 'eduCourse', 'eduEnrollment', 'eduClassSession');

-- Student (9002)
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT '9002', id
FROM sys_base_menus
WHERE path IN ('dashboard', 'person', 'eduCourse', 'eduClassSession');

-- Manager (9003)
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT '9003', id
FROM sys_base_menus
WHERE path IN ('dashboard', 'person', 'eduOrganization', 'eduCourse', 'eduEnrollment', 'eduClassSession');

-- Reset casbin policies for these roles
DELETE FROM casbin_rule
WHERE v0 IN ('9001', '9002', '9003');

-- Teacher (9001) policies
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
  ('p', '9001', '/menu/getMenu', 'POST'),
  ('p', '9001', '/user/setUserAuthority', 'POST'),
  ('p', '9001', '/jwt/jsonInBlacklist', 'POST'),
  ('p', '9001', '/user/getUserInfo', 'GET'),
  ('p', '9001', '/user/getUserList', 'POST'),
  ('p', '9001', '/eduOrganization/createEduOrganization', 'POST'),
  ('p', '9001', '/eduOrganization/deleteEduOrganization', 'DELETE'),
  ('p', '9001', '/eduOrganization/deleteEduOrganizationByIds', 'DELETE'),
  ('p', '9001', '/eduOrganization/updateEduOrganization', 'PUT'),
  ('p', '9001', '/eduOrganization/findEduOrganization', 'GET'),
  ('p', '9001', '/eduOrganization/getEduOrganizationList', 'GET'),
  ('p', '9001', '/eduCourse/createEduCourse', 'POST'),
  ('p', '9001', '/eduCourse/deleteEduCourse', 'DELETE'),
  ('p', '9001', '/eduCourse/deleteEduCourseByIds', 'DELETE'),
  ('p', '9001', '/eduCourse/updateEduCourse', 'PUT'),
  ('p', '9001', '/eduCourse/findEduCourse', 'GET'),
  ('p', '9001', '/eduCourse/getEduCourseList', 'GET'),
  ('p', '9001', '/eduEnrollment/createEduEnrollment', 'POST'),
  ('p', '9001', '/eduEnrollment/deleteEduEnrollment', 'DELETE'),
  ('p', '9001', '/eduEnrollment/deleteEduEnrollmentByIds', 'DELETE'),
  ('p', '9001', '/eduEnrollment/updateEduEnrollment', 'PUT'),
  ('p', '9001', '/eduEnrollment/findEduEnrollment', 'GET'),
  ('p', '9001', '/eduEnrollment/getEduEnrollmentList', 'GET'),
  ('p', '9001', '/eduEnrollment/consumptionClass', 'POST'),
  ('p', '9001', '/eduEnrollment/addSession', 'POST'),
  ('p', '9001', '/eduClassSession/getStudentsWithLessThanFiveSessions', 'GET'),
  ('p', '9001', '/eduClassSession/createEduClassSession', 'POST'),
  ('p', '9001', '/eduClassSession/deleteEduClassSession', 'DELETE'),
  ('p', '9001', '/eduClassSession/deleteEduClassSessionByIds', 'DELETE'),
  ('p', '9001', '/eduClassSession/updateEduClassSession', 'PUT'),
  ('p', '9001', '/eduClassSession/findEduClassSession', 'GET'),
  ('p', '9001', '/eduClassSession/getEduClassSessionList', 'GET');

-- Student (9002) policies
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
  ('p', '9002', '/menu/getMenu', 'POST'),
  ('p', '9002', '/user/setUserAuthority', 'POST'),
  ('p', '9002', '/jwt/jsonInBlacklist', 'POST'),
  ('p', '9002', '/user/getUserInfo', 'GET'),
  ('p', '9002', '/eduCourse/findEduCourse', 'GET'),
  ('p', '9002', '/eduCourse/getEduCourseList', 'GET'),
  ('p', '9002', '/eduClassSession/findEduClassSession', 'GET'),
  ('p', '9002', '/eduClassSession/getEduClassSessionList', 'GET');

-- Manager (9003) policies
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
  ('p', '9003', '/menu/getMenu', 'POST'),
  ('p', '9003', '/user/setUserAuthority', 'POST'),
  ('p', '9003', '/jwt/jsonInBlacklist', 'POST'),
  ('p', '9003', '/user/getUserInfo', 'GET'),
  ('p', '9003', '/user/getUserList', 'POST'),
  ('p', '9003', '/eduOrganization/createEduOrganization', 'POST'),
  ('p', '9003', '/eduOrganization/deleteEduOrganization', 'DELETE'),
  ('p', '9003', '/eduOrganization/deleteEduOrganizationByIds', 'DELETE'),
  ('p', '9003', '/eduOrganization/updateEduOrganization', 'PUT'),
  ('p', '9003', '/eduOrganization/findEduOrganization', 'GET'),
  ('p', '9003', '/eduOrganization/getEduOrganizationList', 'GET'),
  ('p', '9003', '/eduCourse/createEduCourse', 'POST'),
  ('p', '9003', '/eduCourse/deleteEduCourse', 'DELETE'),
  ('p', '9003', '/eduCourse/deleteEduCourseByIds', 'DELETE'),
  ('p', '9003', '/eduCourse/updateEduCourse', 'PUT'),
  ('p', '9003', '/eduCourse/findEduCourse', 'GET'),
  ('p', '9003', '/eduCourse/getEduCourseList', 'GET'),
  ('p', '9003', '/eduEnrollment/createEduEnrollment', 'POST'),
  ('p', '9003', '/eduEnrollment/deleteEduEnrollment', 'DELETE'),
  ('p', '9003', '/eduEnrollment/deleteEduEnrollmentByIds', 'DELETE'),
  ('p', '9003', '/eduEnrollment/updateEduEnrollment', 'PUT'),
  ('p', '9003', '/eduEnrollment/findEduEnrollment', 'GET'),
  ('p', '9003', '/eduEnrollment/getEduEnrollmentList', 'GET'),
  ('p', '9003', '/eduEnrollment/consumptionClass', 'POST'),
  ('p', '9003', '/eduEnrollment/addSession', 'POST'),
  ('p', '9003', '/eduClassSession/getStudentsWithLessThanFiveSessions', 'GET'),
  ('p', '9003', '/eduClassSession/createEduClassSession', 'POST'),
  ('p', '9003', '/eduClassSession/deleteEduClassSession', 'DELETE'),
  ('p', '9003', '/eduClassSession/deleteEduClassSessionByIds', 'DELETE'),
  ('p', '9003', '/eduClassSession/updateEduClassSession', 'PUT'),
  ('p', '9003', '/eduClassSession/findEduClassSession', 'GET'),
  ('p', '9003', '/eduClassSession/getEduClassSessionList', 'GET');

COMMIT;
