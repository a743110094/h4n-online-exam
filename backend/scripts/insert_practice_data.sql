-- 插入示例推荐练习数据

-- 插入推荐练习
INSERT INTO practice_recommendations (title, description, subject_id, difficulty, question_count, estimated_time, rating, knowledge_point, question_types, is_active, created_at, updated_at) VALUES
('二叉树遍历专项练习', '深入理解二叉树的前序、中序、后序遍历', 1, 3, 15, 20, 4.8, '二叉树遍历', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),
('排序算法基础', '掌握冒泡、选择、插入等基础排序算法', 1, 2, 20, 25, 4.6, '排序算法', '["single_choice", "true_false"]', true, NOW(), NOW()),
('进程调度算法', '理解FCFS、SJF、RR等调度算法', 2, 4, 12, 30, 4.9, '进程调度', '["single_choice", "short_answer"]', true, NOW(), NOW()),
('数据库索引优化', '掌握数据库索引的创建和优化技巧', 3, 3, 18, 35, 4.7, '数据库索引', '["single_choice", "multiple_choice"]', true, NOW(), NOW()),
('网络协议基础', '理解TCP/IP、HTTP等网络协议', 4, 2, 25, 40, 4.5, '网络协议', '["single_choice", "true_false"]', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 插入一些示例练习记录（为学生用户）
INSERT INTO practice_records (user_id, subject_id, title, description, question_ids, total_count, correct_count, wrong_count, score, duration, difficulty, practice_type, is_completed, created_at, updated_at) VALUES
(3, 1, '栈和队列基础', '练习栈和队列的基本操作', '[1,2,3,4,5,6,7,8,9,10]', 10, 8, 2, 80, 900, 2, 'sequence', true, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
(3, 1, '递归算法练习', '练习递归算法的实现', '[11,12,13,14,15,16,17,18]', 8, 7, 1, 87, 720, 3, 'random', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(3, 2, '内存管理', '操作系统内存管理相关题目', '[19,20,21,22,23,24,25,26,27,28,29,30,31,32,33]', 15, 13, 2, 86, 1080, 3, 'sequence', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;