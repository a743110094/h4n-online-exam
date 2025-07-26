-- 插入更多推荐练习测试数据
-- 每个练习至少包含5个题目

-- 插入更多推荐练习
INSERT INTO practice_recommendations (title, description, subject_id, difficulty, question_count, estimated_time, rating, knowledge_point, question_types, is_active, created_at, updated_at) VALUES
-- 数学相关练习
('函数基础综合练习', '掌握二次函数、三角函数等基础函数知识', 1, 2, 8, 15, 4.7, '函数基础', '["single_choice", "judge", "essay"]', true, NOW(), NOW()),
('数学解题技巧训练', '提升数学解题能力和思维逻辑', 1, 3, 10, 25, 4.8, '解题技巧', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),
('高等数学入门', '导数、积分等高等数学基础概念', 1, 4, 6, 30, 4.6, '高等数学', '["single_choice", "essay"]', true, NOW(), NOW()),

-- 语文相关练习
('古典文学赏析', '深入理解古诗词和文言文的内涵', 2, 3, 7, 20, 4.9, '古典文学', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),
('现代汉语基础', '掌握现代汉语语法和修辞手法', 2, 2, 9, 18, 4.5, '现代汉语', '["single_choice", "multiple_choice", "judge"]', true, NOW(), NOW()),
('文学常识大全', '中外文学作家作品及文学流派', 2, 2, 12, 22, 4.4, '文学常识', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),

-- 英语相关练习
('英语语法精讲', '系统学习英语语法规则和用法', 3, 2, 8, 16, 4.7, '英语语法', '["single_choice", "multiple_choice", "judge"]', true, NOW(), NOW()),
('英语词汇扩展', '提升英语词汇量和理解能力', 3, 1, 15, 20, 4.6, '英语词汇', '["single_choice", "essay"]', true, NOW(), NOW()),
('英语阅读理解', '提高英语阅读理解和翻译能力', 3, 3, 10, 25, 4.8, '阅读理解', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),

-- 物理相关练习
('经典力学专题', '牛顿定律、运动学等力学基础', 4, 3, 9, 28, 4.9, '经典力学', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),
('电磁学基础', '电场、磁场、电磁感应等概念', 4, 4, 7, 35, 4.7, '电磁学', '["single_choice", "judge", "essay"]', true, NOW(), NOW()),
('物理实验方法', '物理实验设计和数据分析', 4, 3, 11, 30, 4.5, '物理实验', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),

-- 化学相关练习
('化学元素周期律', '元素周期表规律和元素性质', 5, 2, 8, 15, 4.8, '元素周期律', '["single_choice", "judge"]', true, NOW(), NOW()),
('有机化学入门', '有机化合物结构和反应机理', 5, 4, 6, 40, 4.6, '有机化学', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),
('化学实验安全', '化学实验操作规范和安全知识', 5, 1, 10, 12, 4.9, '实验安全', '["single_choice", "judge"]', true, NOW(), NOW()),

-- 跨学科综合练习
('理科综合训练', '数学、物理、化学综合应用', 1, 4, 15, 45, 4.7, '理科综合', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),
('文科综合练习', '语文、历史、地理知识综合', 2, 3, 12, 35, 4.6, '文科综合', '["single_choice", "multiple_choice", "essay"]', true, NOW(), NOW()),
('基础学科入门', '各学科基础知识点梳理', 1, 1, 20, 25, 4.5, '基础知识', '["single_choice", "judge"]', true, NOW(), NOW()),

-- 专项能力训练
('逻辑思维训练', '提升逻辑推理和分析能力', 1, 3, 8, 20, 4.8, '逻辑思维', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),
('记忆力提升练习', '通过练习提高记忆效率', 2, 2, 10, 15, 4.4, '记忆训练', '["single_choice", "judge"]', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 插入更多练习记录示例
INSERT INTO practice_records (user_id, subject_id, title, description, question_ids, total_count, correct_count, wrong_count, score, duration, difficulty, practice_type, is_completed, created_at, updated_at) VALUES
(3, 1, '函数综合练习', '二次函数和三角函数综合训练', '[1,2,3,4,5]', 5, 4, 1, 80, 600, 2, 'sequence', true, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
(3, 2, '古诗词专项', '古典诗词理解和赏析', '[6,7,8,9,10]', 5, 5, 0, 100, 480, 3, 'random', true, NOW() - INTERVAL '5 hours', NOW() - INTERVAL '5 hours'),
(3, 3, '英语语法强化', '英语语法规则和应用', '[11,12,13,14,15]', 5, 3, 2, 60, 720, 2, 'sequence', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(3, 4, '物理力学基础', '牛顿定律和运动学', '[16,17,18,19,20]', 5, 4, 1, 80, 900, 3, 'random', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(3, 5, '化学基础概念', '元素周期表和化学键', '[1,6,11,16,20]', 5, 5, 0, 100, 360, 1, 'sequence', true, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
(3, 1, '数学思维训练', '逻辑推理和问题解决', '[2,4,6,8,10]', 5, 3, 2, 60, 840, 4, 'random', true, NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),
(3, 2, '文学常识测试', '中外文学作家作品', '[7,9,11,13,15]', 5, 4, 1, 80, 540, 2, 'sequence', true, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
(3, 3, '英语词汇扩展', '常用词汇和短语', '[12,14,16,18,20]', 5, 5, 0, 100, 420, 1, 'random', true, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days')
ON CONFLICT DO NOTHING;