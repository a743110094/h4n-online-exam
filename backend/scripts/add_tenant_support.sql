-- 添加多租户支持的数据库迁移脚本

-- 创建租户表
CREATE TABLE IF NOT EXISTS tenants (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认演示租户
INSERT INTO tenants (id, name, code, description, is_active) 
VALUES (100, '演示租户', 'demo', '默认演示租户', true)
ON CONFLICT (id) DO NOTHING;

-- 为所有现有表添加租户字段
ALTER TABLE users ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE subjects ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE questions ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE papers ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE exams ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE exam_records ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE answers ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE practice_records ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE practice_answers ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE practice_recommendations ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;
ALTER TABLE ai_chats ADD COLUMN IF NOT EXISTS tenant_id INTEGER NOT NULL DEFAULT 100;

-- 为租户字段创建索引
CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subjects_tenant_id ON subjects(tenant_id);
CREATE INDEX IF NOT EXISTS idx_questions_tenant_id ON questions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_papers_tenant_id ON papers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_exams_tenant_id ON exams(tenant_id);
CREATE INDEX IF NOT EXISTS idx_exam_records_tenant_id ON exam_records(tenant_id);
CREATE INDEX IF NOT EXISTS idx_answers_tenant_id ON answers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_practice_records_tenant_id ON practice_records(tenant_id);
CREATE INDEX IF NOT EXISTS idx_practice_answers_tenant_id ON practice_answers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_practice_recommendations_tenant_id ON practice_recommendations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_chats_tenant_id ON ai_chats(tenant_id);

-- 更新现有数据的租户ID为默认值100
UPDATE users SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE subjects SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE questions SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE papers SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE exams SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE exam_records SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE answers SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE practice_records SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE practice_answers SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE practice_recommendations SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;
UPDATE ai_chats SET tenant_id = 100 WHERE tenant_id IS NULL OR tenant_id = 0;