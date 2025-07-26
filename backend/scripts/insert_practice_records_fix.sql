-- 修正练习记录数据，只使用存在的用户ID (1,2,3)

INSERT INTO practice_records (user_id, subject_id, title, description, question_ids, total_count, correct_count, wrong_count, score, duration, difficulty, practice_type, is_completed, created_at, updated_at) VALUES
(1, 1, '基础代数练习', '完成基础代数综合练习', '[1,2,6,7,12,13,17,21]', 8, 7, 1, 87, 45, 1, 'recommendation', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
(1, 1, '高等数学练习', '完成高等数学进阶练习', '[9,10,14,18,19,25]', 6, 5, 1, 83, 55, 3, 'recommendation', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(2, 2, '古诗文阅读', '完成古诗文阅读专项练习', '[26,27,30,33,35,41,45,50]', 8, 7, 1, 88, 50, 2, 'recommendation', true, NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours'),
(2, 2, '语言基础练习', '正在进行语言基础知识练习', '[26,28,37,42,43,46]', 6, 0, 0, 0, 0, 1, 'recommendation', false, NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
(3, 3, '英语语法练习', '完成英语语法基础练习', '[51,52,54,55,57,62,67,72]', 8, 7, 1, 88, 35, 1, 'recommendation', true, NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),
(3, 3, '英语阅读写作', '完成英语阅读写作练习', '[59,60,68,69,70,74,76]', 7, 5, 2, 71, 48, 2, 'recommendation', true, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
(1, 4, '力学基础练习', '完成力学基础训练', '[77,78,83,87,90,95,102]', 7, 6, 1, 86, 42, 1, 'recommendation', true, NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
(1, 4, '电磁学练习', '正在进行电磁学专项练习', '[79,82,84,91,97,98,101,102]', 8, 0, 0, 0, 0, 2, 'recommendation', false, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
(2, 5, '化学基础练习', '完成化学基础知识练习', '[103,104,106,109,116,118]', 6, 5, 1, 83, 38, 1, 'recommendation', true, NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days'),
(2, 5, '化学反应练习', '完成化学反应原理练习', '[105,108,113,114,117,119,123]', 7, 6, 1, 86, 52, 2, 'recommendation', true, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
(3, 1, '理科综合练习', '完成理科综合基础练习', '[1,21,77,78,103,104,109,115,118]', 9, 8, 1, 89, 65, 1, 'recommendation', true, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
(3, 2, '文科综合练习', '完成文科综合训练', '[26,28,34,42,51,52,62,72]', 8, 7, 1, 88, 58, 1, 'recommendation', true, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days');