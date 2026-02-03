-- Fix mojibake by reinterpreting latin1 bytes as utf8mb4

START TRANSACTION;

UPDATE sys_authorities
SET
  authority_name = CONVERT(CAST(CONVERT(authority_name USING latin1) AS BINARY) USING utf8mb4),
  default_router = CONVERT(CAST(CONVERT(default_router USING latin1) AS BINARY) USING utf8mb4)
WHERE authority_id IS NOT NULL;

UPDATE sys_users
SET
  username = CONVERT(CAST(CONVERT(username USING latin1) AS BINARY) USING utf8mb4),
  nick_name = CONVERT(CAST(CONVERT(nick_name USING latin1) AS BINARY) USING utf8mb4),
  header_img = CONVERT(CAST(CONVERT(header_img USING latin1) AS BINARY) USING utf8mb4),
  phone = CONVERT(CAST(CONVERT(phone USING latin1) AS BINARY) USING utf8mb4),
  email = CONVERT(CAST(CONVERT(email USING latin1) AS BINARY) USING utf8mb4)
WHERE id IS NOT NULL;

UPDATE edu_organization
SET
  name = CONVERT(CAST(CONVERT(name USING latin1) AS BINARY) USING utf8mb4),
  brief = CONVERT(CAST(CONVERT(brief USING latin1) AS BINARY) USING utf8mb4),
  details = CONVERT(CAST(CONVERT(details USING latin1) AS BINARY) USING utf8mb4),
  contact_info = CONVERT(CAST(CONVERT(contact_info USING latin1) AS BINARY) USING utf8mb4),
  address = CONVERT(CAST(CONVERT(address USING latin1) AS BINARY) USING utf8mb4)
WHERE id IS NOT NULL;

UPDATE edu_course
SET
  course_name = CONVERT(CAST(CONVERT(course_name USING latin1) AS BINARY) USING utf8mb4),
  description = CONVERT(CAST(CONVERT(description USING latin1) AS BINARY) USING utf8mb4),
  image_url = CONVERT(CAST(CONVERT(image_url USING latin1) AS BINARY) USING utf8mb4)
WHERE id IS NOT NULL;

UPDATE edu_class_session
SET
  reason = CONVERT(CAST(CONVERT(reason USING latin1) AS BINARY) USING utf8mb4),
  course_name = CONVERT(CAST(CONVERT(course_name USING latin1) AS BINARY) USING utf8mb4),
  user_name = CONVERT(CAST(CONVERT(user_name USING latin1) AS BINARY) USING utf8mb4)
WHERE id IS NOT NULL;

COMMIT;
