<template>
  <div class="question-preview">
    <!-- 左右布局容器 -->
    <div class="preview-layout">
      <!-- 左侧：题目信息 -->
      <div class="left-panel">
        <!-- 题目信息 -->
        <div class="question-header">
          <div class="question-meta">
            <div class="meta-item">
              <span class="meta-label">题目类型：</span>
              <el-tag :type="getTypeTagType(question.type)" class="type-tag">
                {{ getTypeText(question.type) }}
              </el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">难度：</span>
              <el-tag
                :type="getDifficultyTagType(question.difficulty)"
                class="difficulty-tag"
              >
                {{ getDifficultyText(question.difficulty) }}
              </el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">状态：</span>
              <el-tag :type="getStatusTagType(question.status)" class="status-tag">
                {{ getStatusText(question.status) }}
              </el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">分值：</span>
              <span class="score-info">{{ question.score }} 分</span>
            </div>
          </div>

          <div class="question-info">
            <div class="info-item">
              <span class="info-label">科目：</span>
              <span class="info-value">{{ question.subject?.name || "未分类" }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">知识点：</span>
              <span class="info-value">{{ question.knowledgePoint || "未分类" }}</span>
            </div>
          </div>
        </div>

        <!-- 题目内容 -->
        <div class="question-content">
          <div class="question-text">
            {{ question.content }}
          </div>

          <!-- 填空题答案 -->
          <div v-if="question.type === 'fill'" class="fill-answers">
            <div class="answer-section">
              <h4 class="answer-title">参考答案：</h4>
              <div class="answer-list">
                <div
                  v-for="(answer, index) in question.fillAnswers"
                  :key="index"
                  class="answer-item"
                >
                  <span class="answer-index">{{ index + 1 }}.</span>
                  <span class="answer-text">{{ answer }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 简答题答案 -->
          <div v-if="question.type === 'essay'" class="essay-answer">
            <div class="answer-section">
              <h4 class="answer-title">参考答案：</h4>
              <div class="answer-content">
                {{ question.essayAnswer }}
              </div>
            </div>

            <div v-if="question.scoringCriteria" class="scoring-section">
              <h4 class="answer-title">评分标准：</h4>
              <div class="scoring-content">
                {{ question.scoringCriteria }}
              </div>
            </div>
          </div>
        </div>

        <!-- 题目标签 -->
        <div v-if="question.tags && question.tags.length > 0" class="question-tags">
          <div class="tags-title">标签：</div>
          <div class="tags-list">
            <el-tag v-for="tag in question.tags" :key="tag" size="small" class="tag-item">
              {{ tag }}
            </el-tag>
          </div>
        </div>

        <!-- 题目统计信息 -->
        <div
          v-if="question.usageCount !== undefined || question.correctRate !== undefined"
          class="question-stats"
        >
          <div v-if="question.usageCount !== undefined" class="status-item">
            <span class="status-label">使用次数：</span>
            <span class="status-value">{{ question.usageCount }}</span>
          </div>

          <div v-if="question.correctRate !== undefined" class="status-item">
            <span class="status-label">正确率：</span>
            <span class="status-value" :class="getCorrectRateClass(question.correctRate)">
              {{ question.correctRate }}%
            </span>
          </div>
        </div>
      </div>

      <!-- 右侧：选项和解析 -->
      <div class="right-panel">
        <!-- 上半部分：选项 -->
        <div class="options-section">
          <!-- 选择题选项 -->
          <div v-if="isChoiceQuestion" class="question-options">
            <h4 class="section-title">选项</h4>
            <div
              v-for="(option, index) in question.options"
              :key="index"
              class="option-item"
              :class="{ correct: option.isCorrect, preview: true }"
            >
              <div class="option-prefix">{{ String.fromCharCode(65 + index) }}.</div>
              <div class="option-text">
                {{ option.text }}
              </div>
              <div v-if="option.isCorrect" class="correct-indicator">
                <el-icon color="#67C23A"><CircleCheck /></el-icon>
              </div>
            </div>
          </div>

          <!-- 判断题 -->
          <div v-if="question.type === 'judge'" class="judge-options">
            <h4 class="section-title">选项</h4>
            <div class="option-item" :class="{ correct: question.judgeAnswer === true }">
              <div class="option-prefix">A.</div>
              <div class="option-text">正确</div>
              <div v-if="question.judgeAnswer === true" class="correct-indicator">
                <el-icon color="#67C23A"><CircleCheck /></el-icon>
              </div>
            </div>
            <div class="option-item" :class="{ correct: question.judgeAnswer === false }">
              <div class="option-prefix">B.</div>
              <div class="option-text">错误</div>
              <div v-if="question.judgeAnswer === false" class="correct-indicator">
                <el-icon color="#67C23A"><CircleCheck /></el-icon>
              </div>
            </div>
          </div>
        </div>

        <!-- 下半部分：解析 -->
        <div class="explanation-section">
          <div v-if="question.explanation" class="question-explanation">
            <h4 class="explanation-title">
              <el-icon><InfoFilled /></el-icon>
              题目解析
            </h4>
            <div class="explanation-content">
              {{ question.explanation }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { CircleCheck, InfoFilled } from "@element-plus/icons-vue";

interface Props {
  question: any;
}

const props = defineProps<Props>();

// 是否为选择题
const isChoiceQuestion = computed(() => {
  return ["single", "multiple"].includes(props.question.type);
});

// 获取类型文本
const getTypeText = (type: string): string => {
  const textMap: Record<string, string> = {
    single: "单选题",
    multiple: "多选题",
    judge: "判断题",
    fill: "填空题",
    essay: "简答题",
  };
  return textMap[type] || "";
};

// 获取类型标签类型
const getTypeTagType = (type: string): string => {
  const typeMap: Record<string, string> = {
    single: "",
    multiple: "success",
    judge: "warning",
    fill: "danger",
    essay: "info",
  };
  return typeMap[type] || "";
};

// 获取难度文本
const getDifficultyText = (difficulty: string): string => {
  const textMap: Record<string, string> = {
    1: "简单",
    2: "中等",
    3: "困难",
  };
  return textMap[difficulty] || "";
};

// 获取难度标签类型
const getDifficultyTagType = (difficulty: string): string => {
  const typeMap: Record<string, string> = {
    1: "success",
    2: "warning",
    3: "danger",
  };
  return typeMap[difficulty] || "";
};

// 获取状态文本
const getStatusText = (status: string): string => {
  const textMap: Record<string, string> = {
    published: "已发布",
    draft: "草稿",
  };
  return textMap[status] || "";
};

// 获取状态标签类型
const getStatusTagType = (status: string): string => {
  const typeMap: Record<string, string> = {
    published: "success",
    draft: "warning",
  };
  return typeMap[status] || "";
};

// 获取正确率样式类
const getCorrectRateClass = (rate: number): string => {
  if (rate >= 90) return "rate-excellent";
  if (rate >= 80) return "rate-good";
  if (rate >= 70) return "rate-fair";
  return "rate-poor";
};
</script>

<style scoped>
.question-preview {
  /* 铺满整个弹出框，左右预留12px */
  width: calc(100% - 24px);
  height: 100%;
  margin: 0 12px;
  background: white;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.preview-layout {
  display: flex;
  height: 100%;
  gap: var(--spacing-lg);
}

.left-panel {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
  border-right: 1px solid var(--border-color-light);
}

.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.options-section {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
  border-bottom: 1px solid var(--border-color-light);
}

.explanation-section {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
  padding-bottom: var(--spacing-sm);
  border-bottom: 2px solid var(--dopamine-blue);
}

.question-status {
  display: flex;
  justify-content: flex-start;
  margin: var(--spacing-md) 0;
  padding: var(--spacing-sm) 0;
}

.status-tag {
  font-weight: 600;
  font-size: 12px;
}

.question-header {
  margin-bottom: var(--spacing-md);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.question-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-md);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.meta-label {
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
  white-space: nowrap;
}

.type-tag,
.difficulty-tag {
  font-weight: 600;
}

.score-info {
  font-size: 14px;
  font-weight: 600;
  color: var(--dopamine-orange);
  background: var(--dopamine-orange-light);
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}

.question-info {
  display: flex;
  gap: var(--spacing-lg);
}

.info-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.info-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.question-content {
  margin-bottom: var(--spacing-md);
}

.question-text {
  font-size: 16px;
  line-height: 1.6;
  color: var(--text-primary);
  margin-bottom: var(--spacing-lg);
  padding: var(--spacing-lg);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border-left: 4px solid var(--dopamine-blue);
}

.question-options,
.judge-options {
  margin-top: var(--spacing-lg);
}

.option-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: calc(var(--spacing-md) * 0.9); /* 缩小10%高度 */
  margin-bottom: var(--spacing-sm);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: white;
  transition: all 0.3s ease;
  min-height: calc(40px * 0.9); /* 缩小10%高度 */
}

.option-item.correct {
  background: rgba(103, 194, 58, 0.1);
  border-color: var(--dopamine-green);
}

.option-prefix {
  font-weight: 600;
  color: var(--text-primary);
  min-width: 24px;
}

.option-text {
  flex: 1;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
}

.correct-indicator {
  flex-shrink: 0;
}

.fill-answers,
.essay-answer {
  margin-top: var(--spacing-lg);
}

.answer-section,
.scoring-section {
  margin-bottom: var(--spacing-lg);
}

.answer-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-md) 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.answer-list {
  background: var(--bg-secondary);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
}

.answer-item {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.answer-item:last-child {
  margin-bottom: 0;
}

.answer-index {
  font-weight: 600;
  color: var(--dopamine-blue);
  min-width: 20px;
}

.answer-text {
  flex: 1;
  color: var(--text-primary);
}

.answer-content,
.scoring-content {
  background: var(--bg-secondary);
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  line-height: 1.6;
  color: var(--text-primary);
  white-space: pre-wrap;
}

.question-explanation {
  padding: var(--spacing-md);
  background: rgba(64, 158, 255, 0.05);
  border-radius: var(--radius-md);
  border: 1px solid rgba(64, 158, 255, 0.2);
  height: 100%;
}

.explanation-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--dopamine-blue);
  margin: 0 0 var(--spacing-md) 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.explanation-content {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  white-space: pre-wrap;
}

.question-tags {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  padding: var(--spacing-sm);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.tags-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.tags-list {
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.tag-item {
  background: var(--dopamine-blue-light);
  color: var(--dopamine-blue);
  border: 1px solid var(--dopamine-blue);
}

.question-stats {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  padding: var(--spacing-sm);
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  border-top: 1px solid var(--border-color);
}

.status-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.status-label {
  font-size: 14px;
  color: var(--text-secondary);
}

.status-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.rate-excellent {
  color: var(--dopamine-green) !important;
}

.rate-good {
  color: var(--dopamine-blue) !important;
}

.rate-fair {
  color: var(--dopamine-orange) !important;
}

.rate-poor {
  color: var(--dopamine-red) !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .question-info {
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .question-meta {
    flex-wrap: wrap;
  }

  .question-status {
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .question-tags {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-sm);
  }
}
</style>
