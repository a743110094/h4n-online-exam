-- 插入科目数据
INSERT INTO subjects (id, name, description, created_at, updated_at) VALUES
(1, '数学', '数学科目', NOW(), NOW()),
(2, '语文', '语文科目', NOW(), NOW()),
(3, '英语', '英语科目', NOW(), NOW()),
(4, '物理', '物理科目', NOW(), NOW()),
(5, '化学', '化学科目', NOW(), NOW()),
(6, '生物', '生物科目', NOW(), NOW()),
(7, '历史', '历史科目', NOW(), NOW()),
(8, '地理', '地理科目', NOW(), NOW()),
(9, '政治', '政治科目', NOW(), NOW()),
(10, '计算机', '计算机科目', NOW(), NOW());

-- 重置序列
SELECT setval('subjects_id_seq', (SELECT MAX(id) FROM subjects));