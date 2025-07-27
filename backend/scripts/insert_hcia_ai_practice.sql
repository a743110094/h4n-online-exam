-- 添加华为AI认证科目和推荐练习

-- 插入华为AI认证科目
INSERT INTO subjects (id, name, description, created_at, updated_at) VALUES
(11, '华为AI认证', 'HCIA-AI华为认证人工智能工程师', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = NOW();

-- 插入华为AI认证推荐练习
INSERT INTO practice_recommendations (title, description, subject_id, difficulty, question_count, estimated_time, rating, knowledge_point, question_types, is_active, created_at, updated_at) VALUES
('华为AI基础认证练习', '掌握人工智能基础概念、机器学习算法和华为AI技术栈', 11, 2, 20, 45, 4.8, '华为AI认证基础', '["single_choice", "short_answer"]', true, NOW(), NOW()),
('MindSpore框架实战', '深入学习华为MindSpore深度学习框架的使用和开发', 11, 3, 15, 35, 4.7, 'MindSpore框架', '["single_choice", "short_answer"]', true, NOW(), NOW()),
('Ascend AI处理器应用', '了解华为Ascend AI芯片架构和应用开发', 11, 4, 12, 40, 4.9, 'Ascend处理器', '["single_choice", "short_answer"]', true, NOW(), NOW()),
('华为云AI服务实践', '掌握华为云ModelArts、语音识别、图像识别等AI服务', 11, 2, 18, 30, 4.6, '华为云AI服务', '["single_choice", "short_answer"]', true, NOW(), NOW()),
('AI算法与数据处理', '学习机器学习算法、深度学习和数据预处理技术', 11, 3, 16, 50, 4.8, 'AI算法基础', '["single_choice", "short_answer"]', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 更新序列
SELECT setval('subjects_id_seq', (SELECT MAX(id) FROM subjects));
SELECT setval('practice_recommendations_id_seq', (SELECT MAX(id) FROM practice_recommendations));