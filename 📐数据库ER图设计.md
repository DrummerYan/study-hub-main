# 📐 StudyHub 数据库 ER 图设计

---

## 核心实体关系

```
┌─────────────────────┐
│  edu_organization   │  组织机构（培训机构）
│  ─────────────────  │
│  id                 │
│  name               │  机构名称
│  brief              │  简介
│  address            │  地址
└─────────────────────┘
         │
         │ 1:N
         ├──────────────────────────────────┐
         │                                  │
         ▼                                  ▼
┌─────────────────────┐          ┌─────────────────────┐
│    edu_teacher      │          │    edu_student      │
│  ─────────────────  │          │  ─────────────────  │
│  id                 │          │  id                 │
│  organization_id    │          │  organization_id    │
│  user_id           │          │  user_id           │
│  name               │          │  name               │
│  phone              │          │  phone              │
│  title              │          │  parent_name        │
│  status             │          │  parent_phone       │
└─────────────────────┘          └─────────────────────┘
         │                                  │
         │                                  │
         │ N:M                              │ 1:N
         │ (through                         │
         │  edu_teacher_course)             │
         │                                  │
         ▼                                  ▼
┌─────────────────────┐          ┌─────────────────────┐
│     edu_course      │          │   edu_enrollment    │
│  ─────────────────  │◄─────────│  ─────────────────  │
│  id                 │   N:1    │  id                 │
│  organization_id    │          │  user_id            │
│  course_name        │          │  student_id         │
│  description        │          │  course_id          │
│  image_url          │          │  total_sessions     │
└─────────────────────┘          │  remaining_sessions │
         │                       └─────────────────────┘
         │                                  │
         │ 1:N                              │ 1:N
         │                                  │
         ▼                                  ▼
┌─────────────────────┐          ┌─────────────────────┐
│   edu_schedule      │          │ edu_class_session   │
│  ─────────────────  │          │  ─────────────────  │
│  id                 │          │  id                 │
│  course_id          │          │  enrollment_id      │
│  teacher_id         │          │  action (add/sub)   │
│  classroom_id       │          │  num_sessions       │
│  class_date         │          │  reason             │
│  start_time         │          │  use_date           │
│  end_time           │          └─────────────────────┘
│  status             │
└─────────────────────┘
         │
         │ N:M (through
         │  edu_schedule_student)
         ▼
┌─────────────────────┐
│ edu_schedule_student│
│  ─────────────────  │
│  schedule_id        │
│  student_id         │
│  enrollment_id      │
│  status             │
└─────────────────────┘


┌─────────────────────┐
│   edu_classroom     │  教室
│  ─────────────────  │
│  id                 │
│  organization_id    │
│  name               │
│  capacity           │
│  status             │
└─────────────────────┘
         ▲
         │ N:1
         │
         │ (classroom_id in edu_schedule)
```

---

## 详细表关系说明

### 1. 组织机构 (edu_organization)
**定位：** 系统的顶级实体，所有数据都归属于某个机构

**关系：**
- 1:N → edu_teacher（一个机构有多个教师）
- 1:N → edu_student（一个机构有多个学员）
- 1:N → edu_course（一个机构有多个课程）
- 1:N → edu_classroom（一个机构有多个教室）

---

### 2. 教师 (edu_teacher)
**定位：** 授课的老师

**关系：**
- N:1 → edu_organization（属于某个机构）
- N:M → edu_course（一个教师可以教多门课，一门课可以有多个教师）
  - 通过中间表：`edu_teacher_course`
- 1:N → edu_schedule（一个教师有多个排课）

**关键字段：**
- `organization_id`: 所属机构
- `user_id`: 关联的系统账号（可登录管理后台）
- `status`: 在职状态（active/inactive/resigned）

---

### 3. 学员 (edu_student)
**定位：** 参加培训的学生

**关系：**
- N:1 → edu_organization（属于某个机构）
- 1:N → edu_enrollment（一个学员可以报名多门课程）
- N:M → edu_schedule（通过 edu_schedule_student）

**关键字段：**
- `organization_id`: 所属机构
- `user_id`: 关联的系统账号（可登录小程序）
- `parent_name`, `parent_phone`: 家长信息

---

### 4. 课程 (edu_course)
**定位：** 培训机构开设的课程

**关系：**
- N:1 → edu_organization（属于某个机构）
- N:M → edu_teacher（通过 edu_teacher_course）
- 1:N → edu_enrollment（一门课程有多个报名记录）
- 1:N → edu_schedule（一门课程有多个排课）

**关键字段：**
- `organization_id`: 所属机构
- `course_name`: 课程名称
- `description`: 课程描述

---

### 5. 报名记录 (edu_enrollment)
**定位：** 学员报名课程的记录（核心业务表）

**关系：**
- N:1 → edu_student（多个报名记录属于同一学员）
- N:1 → edu_course（多个报名记录属于同一课程）
- 1:N → edu_class_session（一个报名记录有多条课时记录）

**关键字段：**
- `student_id`: 学员 ID
- `course_id`: 课程 ID
- `total_sessions`: 总课时数
- `remaining_sessions`: 剩余课时数

**核心逻辑：**
- 每次消课：`remaining_sessions -= 1`
- 每次加课：`total_sessions += X`, `remaining_sessions += X`

---

### 6. 课时记录 (edu_class_session)
**定位：** 记录课时的增加和消耗历史

**关系：**
- N:1 → edu_enrollment（多条记录属于同一个报名记录）

**关键字段：**
- `enrollment_id`: 报名记录 ID
- `action`: 操作类型（add=加课, subtract=消课）
- `num_sessions`: 课时数量
- `reason`: 操作原因
- `use_date`: 操作日期

**核心逻辑：**
- 所有加课/消课操作都会生成记录
- 用于追溯历史和对账

---

### 7. 排课表 (edu_schedule)
**定位：** 课程的具体安排（哪天、哪个时间、哪个老师、哪个教室）

**关系：**
- N:1 → edu_course（多个排课属于同一课程）
- N:1 → edu_teacher（多个排课属于同一教师）
- N:1 → edu_classroom（多个排课使用同一教室）
- N:M → edu_student（通过 edu_schedule_student）

**关键字段：**
- `course_id`: 课程 ID
- `teacher_id`: 教师 ID
- `classroom_id`: 教室 ID
- `class_date`: 上课日期
- `start_time`: 开始时间
- `end_time`: 结束时间
- `status`: 状态（scheduled/ongoing/completed/cancelled）

**核心逻辑：**
- 创建时检测冲突（教师、教室、学员）
- 支持批量创建（周期性排课）

---

### 8. 排课-学员关联表 (edu_schedule_student)
**定位：** 记录哪些学员参加某次课

**关系：**
- N:1 → edu_schedule
- N:1 → edu_student
- N:1 → edu_enrollment

**关键字段：**
- `schedule_id`: 排课 ID
- `student_id`: 学员 ID
- `enrollment_id`: 报名记录 ID
- `status`: 状态（enrolled/attended/absent/leave）

**核心逻辑：**
- 教师点名后更新 `status`
- 未来可用于考勤统计

---

### 9. 教室 (edu_classroom)
**定位：** 上课的场地

**关系：**
- N:1 → edu_organization（属于某个机构）
- 1:N → edu_schedule（一个教室有多个排课）

**关键字段：**
- `organization_id`: 所属机构
- `name`: 教室名称
- `capacity`: 容纳人数
- `status`: 状态（available/unavailable/maintenance）

---

### 10. 教师-课程关联表 (edu_teacher_course)
**定位：** 记录教师能教哪些课程

**关系：**
- N:1 → edu_teacher
- N:1 → edu_course

**关键字段：**
- `teacher_id`: 教师 ID
- `course_id`: 课程 ID
- `is_main_teacher`: 是否主讲教师

**核心逻辑：**
- 排课时只能选择该课程的授课教师
- 支持多个教师教同一门课

---

## 数据流示例

### 示例 1：学员报名流程
```
1. 机构管理员创建学员
   INSERT INTO edu_student (organization_id, name, phone, ...)

2. 学员报名课程
   INSERT INTO edu_enrollment (student_id, course_id, total_sessions=10, remaining_sessions=10)

3. 生成加课记录
   INSERT INTO edu_class_session (enrollment_id, action='add', num_sessions=10, reason='初次报名')
```

### 示例 2：排课流程
```
1. 创建排课
   INSERT INTO edu_schedule (course_id, teacher_id, classroom_id, class_date, start_time, end_time)

2. 添加学员到排课
   INSERT INTO edu_schedule_student (schedule_id, student_id, enrollment_id, status='enrolled')

3. 冲突检测
   SELECT * FROM edu_schedule 
   WHERE teacher_id = ? AND class_date = ? 
   AND (start_time < ? AND end_time > ?)
```

### 示例 3：消课流程
```
1. 教师上完课，执行消课
   UPDATE edu_enrollment SET remaining_sessions = remaining_sessions - 1 WHERE id = ?

2. 生成消课记录
   INSERT INTO edu_class_session (enrollment_id, action='subtract', num_sessions=1, reason='第1课')

3. 更新排课状态
   UPDATE edu_schedule SET status='completed' WHERE id = ?

4. 更新学员考勤
   UPDATE edu_schedule_student SET status='attended' WHERE schedule_id = ? AND student_id = ?
```

---

## 索引优化建议

### 高频查询索引
```sql
-- 按机构查询（几乎所有表）
CREATE INDEX idx_organization_id ON edu_teacher (organization_id);
CREATE INDEX idx_organization_id ON edu_student (organization_id);
CREATE INDEX idx_organization_id ON edu_course (organization_id);

-- 排课查询
CREATE INDEX idx_course_date ON edu_schedule (course_id, class_date);
CREATE INDEX idx_teacher_date ON edu_schedule (teacher_id, class_date);
CREATE INDEX idx_classroom_date ON edu_schedule (classroom_id, class_date);

-- 报名记录查询
CREATE INDEX idx_student_course ON edu_enrollment (student_id, course_id);

-- 课时记录查询
CREATE INDEX idx_enrollment_id ON edu_class_session (enrollment_id);
```

---

## 数据完整性约束

### 外键约束
```sql
-- 确保引用完整性
ALTER TABLE edu_teacher ADD FOREIGN KEY (organization_id) REFERENCES edu_organization(id);
ALTER TABLE edu_student ADD FOREIGN KEY (organization_id) REFERENCES edu_organization(id);
ALTER TABLE edu_course ADD FOREIGN KEY (organization_id) REFERENCES edu_organization(id);
ALTER TABLE edu_enrollment ADD FOREIGN KEY (student_id) REFERENCES edu_student(id);
ALTER TABLE edu_enrollment ADD FOREIGN KEY (course_id) REFERENCES edu_course(id);
```

### 检查约束
```sql
-- 确保剩余课时不为负
ALTER TABLE edu_enrollment ADD CHECK (remaining_sessions >= 0);

-- 确保结束时间大于开始时间
ALTER TABLE edu_schedule ADD CHECK (end_time > start_time);

-- 确保容量大于0
ALTER TABLE edu_classroom ADD CHECK (capacity > 0);
```

---

## 分表分库策略（未来扩展）

### 当数据量达到千万级时
```
按机构 ID 分表
- edu_enrollment_1
- edu_enrollment_2
- ...

按日期分表
- edu_schedule_202601
- edu_schedule_202602
- ...
```

---

**数据库设计完成！下一步：创建表并开始开发！** 🎉
