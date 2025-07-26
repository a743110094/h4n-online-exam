-- 插入更多推荐练习数据，确保每个练习至少包含5个题目
-- 基于现有的123个题目创建多样化的练习

INSERT INTO practice_recommendations (title, description, subject_id, difficulty, question_count, estimated_time, rating, knowledge_point, question_types, is_active, created_at, updated_at) VALUES
-- 数学练习
('基础代数综合练习', '涵盖不等式、对数、函数等基础代数知识点的综合练习', 1, 1, 8, 45, 4.2, '基础代数', '["single", "multiple", "judge"]', true, NOW(), NOW()),
('高等数学进阶', '包含积分、极限、复数等高等数学内容的进阶练习', 1, 3, 6, 60, 4.5, '高等数学', '["single", "essay"]', true, NOW(), NOW()),
('几何与三角函数', '立体几何、三角函数、向量运算的综合训练', 1, 2, 7, 50, 4.3, '几何三角', '["single", "multiple", "essay"]', true, NOW(), NOW()),
('数学建模与应用', '概率统计、数学建模、实际应用问题练习', 1, 2, 5, 40, 4.1, '应用数学', '["single", "essay"]', true, NOW(), NOW()),
('数学竞赛训练', '排列组合、数学归纳法等竞赛级别题目', 1, 3, 6, 75, 4.6, '竞赛数学', '["single", "multiple", "essay"]', true, NOW(), NOW()),

-- 语文练习
('古诗文阅读专项', '古诗词理解、文言文翻译、诗歌赏析综合练习', 2, 2, 8, 55, 4.4, '古诗文', '["single", "multiple", "essay"]', true, NOW(), NOW()),
('现代文阅读理解', '现代文阅读、语言文字运用、阅读理解技巧', 2, 2, 7, 45, 4.2, '现代文阅读', '["single", "judge", "essay"]', true, NOW(), NOW()),
('语言基础知识', '成语理解、字音字形、语法知识基础练习', 2, 1, 6, 35, 4.0, '语言基础', '["single", "judge"]', true, NOW(), NOW()),
('写作技巧训练', '修辞手法、写作手法、语言运用综合训练', 2, 2, 5, 40, 4.3, '写作技巧', '["multiple", "essay"]', true, NOW(), NOW()),
('文学常识与体裁', '文学常识、文学体裁、标点符号等知识点', 2, 1, 7, 30, 3.9, '文学常识', '["single", "judge"]', true, NOW(), NOW()),

-- 英语练习
('英语语法基础', '时态、语态、从句等基础语法知识练习', 3, 1, 8, 40, 4.1, '基础语法', '["single", "multiple", "judge"]', true, NOW(), NOW()),
('英语词汇与搭配', '词汇辨析、固定搭配、介词用法专项练习', 3, 2, 6, 35, 4.0, '词汇搭配', '["single", "multiple"]', true, NOW(), NOW()),
('英语阅读写作', '阅读理解、写作练习、语法改错综合训练', 3, 2, 7, 50, 4.3, '读写综合', '["essay", "single"]', true, NOW(), NOW()),
('英语交际与应用', '疑问句、代词、连词等交际英语练习', 3, 1, 5, 30, 3.8, '交际英语', '["single", "judge"]', true, NOW(), NOW()),
('英语进阶语法', '句子成分、句型转换、习语理解等进阶内容', 3, 3, 6, 45, 4.4, '进阶语法', '["multiple", "essay", "single"]', true, NOW(), NOW()),

-- 物理练习
('力学基础训练', '牛顿定律、运动学、力学计算基础练习', 4, 1, 7, 50, 4.2, '基础力学', '["single", "judge", "essay"]', true, NOW(), NOW()),
('电磁学专项', '电路分析、电磁感应、库仑定律等电磁学内容', 4, 2, 8, 60, 4.4, '电磁学', '["single", "multiple", "essay"]', true, NOW(), NOW()),
('光学与波动', '光学现象、波动光学、简谐运动综合练习', 4, 2, 6, 45, 4.1, '光学波动', '["single", "multiple", "judge"]', true, NOW(), NOW()),
('现代物理学', '原子物理、量子力学、相对论等现代物理内容', 4, 3, 5, 55, 4.5, '现代物理', '["single", "judge", "essay"]', true, NOW(), NOW()),
('物理实验与应用', '流体力学、热力学、能量转换等应用物理', 4, 2, 7, 50, 4.0, '应用物理', '["single", "multiple", "essay"]', true, NOW(), NOW()),

-- 化学练习
('化学基础知识', '原子结构、化学键、元素周期表基础知识', 5, 1, 6, 40, 4.0, '化学基础', '["single", "multiple", "judge"]', true, NOW(), NOW()),
('化学反应原理', '氧化还原、化学平衡、反应类型综合练习', 5, 2, 7, 50, 4.3, '反应原理', '["single", "judge", "essay"]', true, NOW(), NOW()),
('有机化学入门', '有机化合物、化学式、物质分类基础练习', 5, 1, 5, 35, 3.9, '有机化学', '["single", "multiple"]', true, NOW(), NOW()),
('化学实验安全', '实验安全、实验设计、化学用语规范练习', 5, 2, 6, 45, 4.2, '实验安全', '["multiple", "judge", "essay"]', true, NOW(), NOW()),
('化学计算专项', '化学计算、方程式书写、化学与生活应用', 5, 2, 5, 40, 4.1, '化学计算', '["essay", "single"]', true, NOW(), NOW()),

-- 跨学科综合练习
('理科综合基础', '数学、物理、化学基础知识综合练习', 1, 1, 9, 60, 4.0, '理科综合', '["single", "judge"]', true, NOW(), NOW()),
('文科综合训练', '语文、英语基础知识与应用综合练习', 2, 1, 8, 50, 3.9, '文科综合', '["single", "multiple"]', true, NOW(), NOW()),
('逻辑思维训练', '数学逻辑、物理推理、化学分析综合训练', 1, 2, 7, 55, 4.3, '逻辑思维', '["single", "essay"]', true, NOW(), NOW()),
('应用能力测试', '各学科实际应用问题综合测试', 1, 2, 10, 75, 4.2, '应用能力', '["essay", "single"]', true, NOW(), NOW()),
('竞赛预备训练', '各学科竞赛级别题目综合训练', 1, 3, 8, 90, 4.6, '竞赛预备', '["single", "multiple", "essay"]', true, NOW(), NOW());

-- 插入对应的练习记录数据
INSERT INTO practice_records (user_id, subject_id, title, description, question_ids, total_count, correct_count, wrong_count, score, duration, difficulty, practice_type, is_completed, created_at, updated_at) VALUES
(1, 1, '基础代数练习', '完成基础代数综合练习', '[1,2,6,7,12,13,17,21]', 8, 7, 1, 87, 45, 1, 'recommendation', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(1, 1, '高等数学练习', '完成高等数学进阶练习', '[9,10,14,18,19,25]', 6, 5, 1, 83, 55, 3, 'recommendation', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(2, 2, '古诗文阅读', '完成古诗文阅读专项练习', '[26,27,30,33,35,41,45,50]', 8, 7, 1, 88, 50, 2, 'recommendation', true, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
(2, 2, '语言基础练习', '正在进行语言基础知识练习', '[26,28,37,42,43,46]', 6, 0, 0, 0, 0, 1, 'recommendation', false, NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
(3, 3, '英语语法练习', '完成英语语法基础练习', '[51,52,54,55,57,62,67,72]', 8, 7, 1, 88, 35, 1, 'recommendation', true, NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),
(3, 3, '英语阅读写作', '完成英语阅读写作练习', '[59,60,68,69,70,74,76]', 7, 5, 2, 71, 48, 2, 'recommendation', true, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
(4, 4, '力学基础练习', '完成力学基础训练', '[77,78,83,87,90,95,102]', 7, 6, 1, 86, 42, 1, 'recommendation', true, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
(4, 4, '电磁学练习', '正在进行电磁学专项练习', '[79,82,84,91,97,98,101,102]', 8, 0, 0, 0, 0, 2, 'recommendation', false, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
(5, 5, '化学基础练习', '完成化学基础知识练习', '[103,104,106,109,116,118]', 6, 5, 1, 83, 38, 1, 'recommendation', true, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
(5, 5, '化学反应练习', '完成化学反应原理练习', '[105,108,113,114,117,119,123]', 7, 6, 1, 86, 52, 2, 'recommendation', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(1, 1, '理科综合练习', '完成理科综合基础练习', '[1,21,77,78,103,104,109,115,118]', 9, 8, 1, 89, 65, 1, 'recommendation', true, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
(2, 2, '文科综合练习', '完成文科综合训练', '[26,28,34,42,51,52,62,72]', 8, 7, 1, 88, 58, 1, 'recommendation', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days');