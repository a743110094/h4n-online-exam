-- 在线考试系统数据库初始化脚本
-- PostgreSQL版本
-- 注意：数据库已在docker-compose.yml中创建

-- 创建管理员用户
-- 注意：这里的密码是 'admin123' 的bcrypt哈希值
-- 在实际部署时应该修改默认密码
INSERT INTO users (username, password, name, email, role, created_at, updated_at) 
VALUES (
    'admin',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    '系统管理员',
    'admin@example.com',
    'admin',
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- 创建示例科目
INSERT INTO subjects (name, description, created_at, updated_at) VALUES
('数学', '数学相关题目', NOW(), NOW()),
('语文', '语文相关题目', NOW(), NOW()),
('英语', '英语相关题目', NOW(), NOW()),
('物理', '物理相关题目', NOW(), NOW()),
('化学', '化学相关题目', NOW(), NOW())
ON CONFLICT (name) DO NOTHING;

-- 创建示例教师用户
INSERT INTO users (username, password, name, email, role, created_at, updated_at) VALUES
('teacher1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '张老师', 'teacher1@example.com', 'teacher', NOW(), NOW()),
('teacher2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '李老师', 'teacher2@example.com', 'teacher', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 创建示例学生用户
INSERT INTO users (username, password, name, email, role, created_at, updated_at) VALUES
('student1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '王同学', 'student1@example.com', 'student', NOW(), NOW()),
('student2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '刘同学', 'student2@example.com', 'student', NOW(), NOW()),
('student3', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '陈同学', 'student3@example.com', 'student', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- 创建示例练习推荐
INSERT INTO practice_recommendations (title, description, subject_id, difficulty, question_count, estimated_time, rating, created_at, updated_at) VALUES
('数学基础练习', '包含基础数学运算和概念的练习题目', 1, 1, 20, 30, 4.5, NOW(), NOW()),
('语文阅读理解', '提升语文阅读理解能力的专项练习', 2, 2, 15, 25, 4.3, NOW(), NOW()),
('英语语法专项', '英语语法知识点的集中练习', 3, 2, 18, 35, 4.6, NOW(), NOW()),
('物理力学基础', '物理力学相关概念和计算练习', 4, 3, 12, 40, 4.8, NOW(), NOW()),
('化学方程式', '化学方程式配平和反应类型练习', 5, 3, 16, 30, 4.4, NOW(), NOW())
ON CONFLICT DO NOTHING;