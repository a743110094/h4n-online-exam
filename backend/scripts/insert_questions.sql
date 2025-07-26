-- 插入20个示例题目数据
-- 注意：这里假设科目ID 1-5 分别对应数学、语文、英语、物理、化学
-- 创建者ID使用2（teacher1）

INSERT INTO questions (subject_id, type, title, content, options, answer, explanation, difficulty, score, created_by, status, knowledge_point, usage_count, correct_rate, created_at, updated_at) VALUES
-- 数学题目 (科目ID: 1)
(1, 'single', '二次函数的顶点坐标', '函数 f(x) = x² - 4x + 3 的顶点坐标是？', '["A. (2, -1)", "B. (2, 1)", "C. (-2, -1)", "D. (-2, 1)"]', 'A', '二次函数 f(x) = ax² + bx + c 的顶点坐标为 (-b/2a, f(-b/2a))。对于 f(x) = x² - 4x + 3，顶点横坐标为 4/2 = 2，纵坐标为 f(2) = 4 - 8 + 3 = -1。', 2, 5, 2, 'published', '二次函数', 15, 0.85, NOW(), NOW()),
(1, 'single', '三角函数值', 'sin(π/6) 的值是？', '["A. 1/2", "B. √3/2", "C. √2/2", "D. 1"]', 'A', 'sin(π/6) = sin(30°) = 1/2，这是基本的三角函数值。', 1, 3, 2, 'published', '三角函数', 22, 0.92, NOW(), NOW()),
(1, 'multiple', '一元二次方程的解法', '解一元二次方程可以使用哪些方法？', '["A. 因式分解法", "B. 配方法", "C. 公式法", "D. 图像法"]', 'A,B,C', '解一元二次方程的常用方法包括因式分解法、配方法和公式法。图像法可以用来观察解的情况，但不是直接的解法。', 2, 4, 2, 'published', '一元二次方程', 18, 0.78, NOW(), NOW()),
(1, 'judge', '函数单调性', '函数 f(x) = x³ 在整个定义域上单调递增。', '[]', 'true', '函数 f(x) = x³ 的导数 f\'(x) = 3x² ≥ 0 对所有实数成立，且仅在 x = 0 时等于 0，因此函数在整个定义域上单调递增。', 2, 3, 2, 'published', '函数单调性', 12, 0.88, NOW(), NOW()),
(1, 'essay', '导数计算', '求函数 f(x) = 3x² + 2x - 1 的导数。', '[]', 'f\'(x) = 6x + 2', '根据导数的基本公式，常数的导数为0，x^n的导数为nx^(n-1)，所以 f\'(x) = 3×2x + 2×1 - 0 = 6x + 2。', 1, 4, 2, 'published', '导数', 25, 0.76, NOW(), NOW()),

-- 语文题目 (科目ID: 2)
(2, 'single', '古诗词理解', '"春蚕到死丝方尽，蜡炬成灰泪始干"出自哪位诗人的作品？', '["A. 李白", "B. 杜甫", "C. 李商隐", "D. 白居易"]', 'C', '这句诗出自李商隐的《无题》，表达了深挚的爱情和无私的奉献精神。', 2, 3, 2, 'published', '古诗词', 20, 0.82, NOW(), NOW()),
(2, 'single', '文言文翻译', '"学而时习之，不亦说乎"中的"说"字的含义是？', '["A. 说话", "B. 高兴", "C. 说明", "D. 劝说"]', 'B', '在这句话中，"说"通"悦"，意思是高兴、愉快。整句话的意思是：学了知识并且按时复习，不是很高兴的事吗？', 2, 3, 2, 'published', '文言文', 16, 0.75, NOW(), NOW()),
(2, 'multiple', '修辞手法', '下列句子中使用了哪些修辞手法？"桂花飘香十里，如同天上的仙女撒下的花瓣。"', '["A. 比喻", "B. 夸张", "C. 拟人", "D. 排比"]', 'A,B', '这句话使用了比喻（把桂花比作仙女撒下的花瓣）和夸张（飘香十里）两种修辞手法。', 2, 4, 2, 'published', '修辞手法', 14, 0.71, NOW(), NOW()),
(2, 'judge', '汉字知识', '"美"字是象形字。', '[]', 'false', '"美"字是会意字，不是象形字。它由"羊"和"大"组成，古人认为大羊为美。', 1, 2, 2, 'published', '汉字知识', 28, 0.89, NOW(), NOW()),
(2, 'essay', '诗歌赏析', '请简要分析"落红不是无情物，化作春泥更护花"的含义。', '[]', '表达了诗人虽然离开官场，但仍然关心国家前途和人民命运的思想感情，体现了无私奉献的精神。', '这句诗运用比喻的修辞手法，以落花为喻，表达了诗人的爱国情怀和奉献精神。', 3, 5, 2, 'published', '诗歌赏析', 11, 0.68, NOW(), NOW()),

-- 英语题目 (科目ID: 3)
(3, 'single', '语法选择', 'She _____ to the library every day.', '["A. go", "B. goes", "C. going", "D. gone"]', 'B', '主语she是第三人称单数，在一般现在时中，动词要用第三人称单数形式goes。', 1, 3, 2, 'published', '语法', 30, 0.93, NOW(), NOW()),
(3, 'single', '词汇理解', 'What does "ambitious" mean?', '["A. lazy", "B. having strong desire for success", "C. tired", "D. angry"]', 'B', 'Ambitious means having a strong desire for success or achievement.', 2, 3, 2, 'published', '词汇', 24, 0.87, NOW(), NOW()),
(3, 'multiple', '时态用法', 'Which sentences use the present perfect tense correctly?', '["A. I have finished my homework.", "B. She has been to Paris twice.", "C. They have ate dinner.", "D. We have lived here for five years."]', 'A,B,D', 'Present perfect tense is formed with have/has + past participle. "ate" should be "eaten" in sentence C.', 3, 4, 2, 'published', '时态', 19, 0.74, NOW(), NOW()),
(3, 'judge', '语法规则', 'In English, adjectives always come before nouns.', '[]', 'false', 'While adjectives usually come before nouns in English, there are exceptions, such as in certain phrases like "something special" or "nothing important".', 2, 3, 2, 'published', '语法规则', 17, 0.81, NOW(), NOW()),
(3, 'essay', '翻译练习', 'Translate: "Knowledge is power."', '[]', '知识就是力量。', 'This is a famous quote by Francis Bacon, emphasizing the importance of knowledge.', 1, 3, 2, 'published', '翻译', 13, 0.85, NOW(), NOW()),

-- 物理题目 (科目ID: 4)
(4, 'single', '力学基础', '一个物体在水平面上做匀速直线运动，它受到的合力是？', '["A. 向前的力", "B. 向后的力", "C. 零", "D. 向上的力"]', 'C', '根据牛顿第一定律，物体做匀速直线运动时，受到的合力为零。', 2, 4, 2, 'published', '力学', 26, 0.88, NOW(), NOW()),
(4, 'single', '电学知识', '欧姆定律的表达式是？', '["A. U = IR", "B. P = UI", "C. W = UIt", "D. Q = It"]', 'A', '欧姆定律表述为：在同一电路中，通过某段导体两端的电压成正比，跟这段导体的电阻成反比，即U = IR。', 1, 3, 2, 'published', '电学', 32, 0.91, NOW(), NOW()),
(4, 'multiple', '能量形式', '下列哪些是机械能的形式？', '["A. 动能", "B. 势能", "C. 内能", "D. 弹性势能"]', 'A,B,D', '机械能包括动能和势能，其中势能又分为重力势能和弹性势能。内能不属于机械能。', 2, 4, 2, 'published', '能量', 21, 0.79, NOW(), NOW()),
(4, 'judge', '光学现象', '光在真空中的传播速度是3×10⁸ m/s。', '[]', 'true', '光在真空中的传播速度确实是3×10⁸ m/s，这是一个重要的物理常数。', 1, 2, 2, 'published', '光学', 35, 0.94, NOW(), NOW()),
(4, 'essay', '运动学计算', '一个物体从静止开始做匀加速直线运动，加速度为2 m/s²，求它在第3秒末的速度。', '[]', 'v = at = 2 × 3 = 6 m/s', '根据匀加速直线运动的速度公式 v = v₀ + at，其中v₀ = 0，a = 2 m/s²，t = 3s。', 2, 5, 2, 'published', '运动学', 18, 0.72, NOW(), NOW()),

-- 化学题目 (科目ID: 5)
(5, 'single', '元素周期表', '氢元素在元素周期表中的原子序数是？', '["A. 1", "B. 2", "C. 3", "D. 4"]', 'A', '氢元素是最简单的元素，只有一个质子，因此原子序数为1。', 1, 2, 2, 'published', '元素周期表', 29, 0.96, NOW(), NOW());