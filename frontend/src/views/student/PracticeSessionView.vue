<template>
  <div class="practice-session">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-loading-directive
        :loading="loading"
        text="正在加载练习数据..."
        background="rgba(0, 0, 0, 0.7)"
      />
    </div>

    <!-- 练习内容 -->
    <div v-else>
      <!-- 练习头部信息 -->
      <div class="practice-header">
        <div class="practice-info">
          <h1 class="practice-title">{{ practiceTitle }}</h1>
          <div class="practice-meta">
            <span class="meta-item">
              <el-icon><Document /></el-icon>
              {{ practiceSubject }}
            </span>
            <span class="meta-item">
              <el-icon><Star /></el-icon>
              难度：{{ getDifficultyText(practiceDifficulty) }}
            </span>
            <span class="meta-item">
              <el-icon><Timer /></el-icon>
              练习模式
            </span>
          </div>
        </div>

        <div class="practice-actions">
          <el-button @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <el-button type="primary" @click="submitPractice">
            <el-icon><Check /></el-icon>
            提交练习
          </el-button>
        </div>
      </div>

      <!-- 题目导航与答题进度合并 -->
      <div class="question-nav">
        <div class="nav-header">
          <div class="nav-title-section">
            <h3>题目导航</h3>
            <div class="progress-info">
              <span class="progress-text">
                已答题：{{ answeredCount }}/{{ questions.length }}
              </span>
              <span class="progress-percentage">
                {{ Math.round((answeredCount / questions.length) * 100) }}%
              </span>
            </div>
          </div>
          <el-button
            type="primary"
            size="small"
            @click="showQuestionNav = !showQuestionNav"
          >
            {{ showQuestionNav ? "收起" : "展开" }}
          </el-button>
        </div>

        <!-- 进度条 -->
        <div class="nav-progress">
          <el-progress
            :percentage="Math.round((answeredCount / questions.length) * 100)"
            color="var(--dopamine-blue)"
            :stroke-width="6"
          />
        </div>

        <div v-show="showQuestionNav" class="nav-content">
          <!-- 左右布局容器 -->
          <div class="nav-layout">
            <!-- 左侧信息面板 -->
            <div class="nav-info-panel">
              <!-- 题型统计 -->
              <div class="type-stats">
                <h4 class="stats-title">题型分布</h4>
                <div class="type-list">
                  <div v-for="(count, type) in questionTypeStats" :key="type" class="type-item">
                    <span class="type-label">{{ getTypeText(type) }}：</span>
                    <span class="type-count">{{ count }}题</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 右侧题目导航 -->
            <div class="nav-grid-panel">
              <!-- 状态说明 -->
              <div class="nav-legend">
                <div class="legend-item">
                  <div class="legend-dot answered"></div>
                  <span>已答题</span>
                </div>
                <div class="legend-item">
                  <div class="legend-dot current"></div>
                  <span>当前题</span>
                </div>
                <div class="legend-item">
                  <div class="legend-dot unanswered"></div>
                  <span>未答题</span>
                </div>
              </div>
              <div 
                ref="navGridRef"
                class="nav-grid"
                @scroll="onNavScroll"
              >
                <div 
                  class="nav-grid-content"
                  :style="{ height: `${totalHeight}px` }"
                >
                  <div
                    v-for="item in visibleItems"
                    :key="item.question.id"
                    class="nav-item"
                    :class="{
                      current: currentQuestionIndex === item.index,
                      answered: answers[item.question.id] !== undefined,
                      unanswered: answers[item.question.id] === undefined,
                    }"
                    :style="{
                      position: 'absolute',
                      top: `${item.top}px`,
                      left: `${item.left}px`
                    }"
                    @click="goToQuestion(item.index)"
                  >
                    {{ item.index + 1 }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && questions.length === 0" class="empty-state">
        <el-empty description="暂无练习题目" />
      </div>

      <!-- 答题区域 -->
      <div v-else-if="currentQuestion" class="answer-area">
        <div class="question-container">
          <!-- 题目信息 -->
          <div class="question-header">
            <div class="question-number">第 {{ currentQuestionIndex + 1 }} 题</div>
            <div class="question-meta">
              <el-tag :type="getTypeTagType(currentQuestion.type)">
                {{ getTypeText(currentQuestion.type) }}
              </el-tag>
              <span class="question-score">{{ currentQuestion.score }} 分</span>
            </div>
          </div>

          <!-- 题目内容 -->
          <div class="question-content">
            <div class="question-text">
              {{ currentQuestion.content }}
            </div>

            <!-- 选择题选项 -->
            <div v-if="isChoiceQuestion" class="question-options">
              <div
                v-for="option in currentQuestion.options"
                :key="option.id"
                class="option-item"
                :class="{ selected: isOptionSelected(option.id) }"
                @click="currentQuestion.type === 'single' ? selectOption(option.id) : toggleOption(option.id)"
              >
                <div class="option-radio">
                  <el-radio
                    v-if="currentQuestion.type === 'single'"
                    :model-value="answers[currentQuestion.id]"
                    :label="option.id"
                    @change="selectOption(option.id)"
                  >
                  </el-radio>
                  <el-checkbox
                    v-else
                    :model-value="isOptionSelected(option.id)"
                    @change="toggleOption(option.id)"
                  >
                  </el-checkbox>
                </div>
                <div class="option-content">
                  <div class="option-prefix">{{ option.id }}.</div>
                  <div class="option-text">
                    {{ option.text.replace(/^[A-Z]\. /, "") }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 判断题 -->
            <div v-if="currentQuestion.type === 'judge'" class="judge-options">
              <div
                class="judge-option"
                :class="{ selected: answers[currentQuestion.id] === true }"
                @click="selectJudge(true)"
              >
                <el-radio
                  :model-value="answers[currentQuestion.id]"
                  :label="true"
                  @change="selectJudge(true)"
                >
                  正确
                </el-radio>
              </div>
              <div
                class="judge-option"
                :class="{ selected: answers[currentQuestion.id] === false }"
                @click="selectJudge(false)"
              >
                <el-radio
                  :model-value="answers[currentQuestion.id]"
                  :label="false"
                  @change="selectJudge(false)"
                >
                  错误
                </el-radio>
              </div>
            </div>

            <!-- 填空题 -->
            <div v-if="currentQuestion.type === 'fill'" class="fill-inputs">
              <div
                v-for="(blank, index) in currentQuestion.blanks || ['']"
                :key="index"
                class="fill-item"
              >
                <label class="fill-label">第 {{ index + 1 }} 空：</label>
                <el-input
                  v-model="fillAnswers[index]"
                  placeholder="请输入答案"
                  @input="updateFillAnswer(index)"
                />
              </div>
            </div>
          </div>

          <!-- 底部按钮区域 -->
          <div class="bottom-actions">
            <!-- 左侧：上一题按钮 -->
            <el-button 
              :disabled="currentQuestionIndex === 0"
              @click="previousQuestion"
              class="nav-button"
            >
              <el-icon><ArrowLeft /></el-icon>
              上一题
            </el-button>
            
            <!-- 右侧：确认答案按钮或下一题按钮 -->
            <el-button
              v-if="hasAnswered && (!answerFeedback.showFeedback || answerFeedback.questionId !== currentQuestion.id)"
              type="primary"
              @click="confirmAnswer"
              :loading="submittingAnswer"
              class="confirm-button"
            >
              <el-icon><Check /></el-icon>
              确认答案
            </el-button>
            
            <!-- 答题后显示下一题按钮 -->
            <el-button
              v-if="answerFeedback.showFeedback && answerFeedback.questionId === currentQuestion.id && currentQuestionIndex < questions.length - 1"
              type="primary"
              @click="nextQuestion"
              class="confirm-button"
            >
              下一题
              <el-icon><ArrowRight /></el-icon>
            </el-button>
            
            <!-- 最后一题答题后显示提交练习按钮 -->
            <el-button 
              v-if="answerFeedback.showFeedback && answerFeedback.questionId === currentQuestion.id && currentQuestionIndex >= questions.length - 1"
              type="success" 
              @click="submitPractice" 
              :loading="loading"
              class="confirm-button"
            >
              <el-icon><Check /></el-icon>
              提交练习
            </el-button>
          </div>

          <!-- 答题反馈 -->
          <div
            v-if="
              answerFeedback.showFeedback &&
              answerFeedback.questionId === currentQuestion.id
            "
            class="answer-feedback"
          >
            <div
              class="feedback-card"
              :class="{
                correct: answerFeedback.isCorrect,
                incorrect: !answerFeedback.isCorrect,
              }"
            >
              <div class="feedback-header">
                <div class="feedback-result">
                  <el-icon class="feedback-icon">
                    <Check v-if="answerFeedback.isCorrect" />
                    <Close v-else />
                  </el-icon>
                  <span class="feedback-text">
                    {{ answerFeedback.isCorrect ? "回答正确" : "回答错误" }}
                  </span>
                </div>
                <div class="feedback-actions">
                  <div class="feedback-score">得分：{{ answerFeedback.score }} 分</div>
                  <el-button
                    type="text"
                    size="small"
                    @click="answerFeedback.showFeedback = false"
                    class="close-feedback-btn"
                  >
                    <el-icon><Close /></el-icon>
                  </el-button>
                </div>
              </div>
              <div v-if="answerFeedback.explanation" class="feedback-explanation">
                <div class="explanation-title">解析：</div>
                <div class="explanation-content">{{ answerFeedback.explanation }}</div>
              </div>
            </div>
          </div>


        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Document,
  Star,
  Timer,
  ArrowLeft,
  Check,
  ArrowRight,
  Close,
} from "@element-plus/icons-vue";
import { startPractice, submitPracticeAnswer, completePractice } from "@/api/practice";

const router = useRouter();
const route = useRoute();

// 练习信息
const practiceId = ref<number>(0);
const practiceTitle = ref("");
const practiceSubject = ref("");
const practiceDifficulty = ref("");
const practiceType = ref("");
const loading = ref(false);

// 题目数据
const questions = ref<any[]>([]);

// 答题状态
const currentQuestionIndex = ref(0);
const answers = reactive<Record<number, any>>({});
const fillAnswers = ref<string[]>([]);
const showQuestionNav = ref(true);

// 答题反馈状态
const answerFeedback = ref<{
  questionId: number | null;
  isCorrect: boolean | null;
  explanation: string;
  score: number;
  showFeedback: boolean;
}>({
  questionId: null,
  isCorrect: null,
  explanation: "",
  score: 0,
  showFeedback: false,
});

// 答案提交状态
const submittingAnswer = ref(false);

// 进度保存相关
const localOperationCount = ref(0);
const progressSaveTimer = ref<NodeJS.Timeout | null>(null);
const lastSaveTime = ref(Date.now());

// 计算属性
const currentQuestion = computed(() => questions.value[currentQuestionIndex.value]);
const answeredCount = computed(() => Object.keys(answers).length);
const isChoiceQuestion = computed(
  () =>
    currentQuestion.value && ["single", "multiple"].includes(currentQuestion.value.type)
);

// 题型统计
const questionTypeStats = computed(() => {
  const stats: Record<string, number> = {};
  questions.value.forEach(question => {
    const type = question.type;
    stats[type] = (stats[type] || 0) + 1;
  });
  return stats;
});

// 判断当前题目是否已答题
const hasAnswered = computed(() => {
  if (!currentQuestion.value) return false;
  const questionId = currentQuestion.value.id;
  const answer = answers[questionId];
  
  if (currentQuestion.value.type === 'single' || currentQuestion.value.type === 'judge') {
    return answer !== undefined && answer !== null;
  } else if (currentQuestion.value.type === 'multiple') {
    return Array.isArray(answer) && answer.length > 0;
  } else if (currentQuestion.value.type === 'fill') {
    return fillAnswers.value.some(ans => ans && ans.trim());
  }
  
  return false;
});

// 虚拟滚动相关
const navGridRef = ref<HTMLElement>();
const itemHeight = 40;
const itemWidth = 40;
const gap = 8;
const containerHeight = 200;
const scrollTop = ref(0);
const containerWidth = ref(0);

// 计算每行显示的项目数
const itemsPerRow = computed(() => {
  if (containerWidth.value === 0) return 10;
  return Math.floor((containerWidth.value + gap) / (itemWidth + gap));
});

// 计算总行数
const totalRows = computed(() => {
  return Math.ceil(questions.value.length / itemsPerRow.value);
});

// 计算总高度
const totalHeight = computed(() => {
  return totalRows.value * (itemHeight + gap) - gap;
});

// 计算可见的开始行
const startRow = computed(() => {
  return Math.floor(scrollTop.value / (itemHeight + gap));
});

// 计算可见的结束行
const endRow = computed(() => {
  const visibleRows = Math.ceil(containerHeight / (itemHeight + gap));
  return Math.min(startRow.value + visibleRows + 1, totalRows.value);
});

// 计算可见的项目
const visibleItems = computed(() => {
  const items = [];
  for (let row = startRow.value; row < endRow.value; row++) {
    for (let col = 0; col < itemsPerRow.value; col++) {
      const index = row * itemsPerRow.value + col;
      if (index < questions.value.length) {
        items.push({
          question: questions.value[index],
          index,
          top: row * (itemHeight + gap),
          left: col * (itemWidth + gap)
        });
      }
    }
  }
  return items;
});

// 滚动事件处理
const onNavScroll = (event: Event) => {
  const target = event.target as HTMLElement;
  scrollTop.value = target.scrollTop;
};

// 更新容器宽度
const updateContainerWidth = () => {
  if (navGridRef.value) {
    containerWidth.value = navGridRef.value.clientWidth;
  }
};

// 初始化
onMounted(async () => {
  await loadPracticeData();
  startProgressSaveTimer();
  
  // 初始化虚拟滚动
  await nextTick();
  updateContainerWidth();
  
  // 监听窗口大小变化
  window.addEventListener('resize', updateContainerWidth);
});

// 组件卸载时清理定时器
onUnmounted(() => {
  if (progressSaveTimer.value) {
    clearInterval(progressSaveTimer.value);
  }
  
  // 清理事件监听器
  window.removeEventListener('resize', updateContainerWidth);
});



// 加载练习数据
const loadPracticeData = async () => {
  try {
    loading.value = true;
    const practiceIdParam = route.params.id as string;
    const practiceTypeParam = route.params.practiceType as string;

    // 从路由参数获取练习ID
    const subjectId = parseInt(practiceIdParam);
    practiceType.value = practiceTypeParam;

    // 调用后端API开始练习
    const difficultyMap: Record<string, number> = {
      easy: 1,
      medium: 2,
      hard: 3,
    };
    const difficultyStr = (route.query.difficulty as string) || "medium";
    const difficultyNum = difficultyMap[difficultyStr] || 2;

    const response = await startPractice({
      subject_id: subjectId,
      practice_type: practiceType.value,
      difficulty: difficultyNum,
      question_count: parseInt(route.query.question_count as string) || 10,
    });

    // 设置练习数据
    practiceId.value = response.practice_id;
    // 处理题目数据，解析options字段
    questions.value = (response.questions || []).map((question: any) => {
      const processedQuestion = { ...question };

      // 如果是选择题，需要解析options字段
      if (question.type === "single" || question.type === "multiple") {
        if (typeof question.options === "string") {
          try {
            const optionsArray = JSON.parse(question.options);
            // 将字符串数组转换为对象数组
            processedQuestion.options = optionsArray.map(
              (text: string, index: number) => ({
                id: String.fromCharCode(65 + index), // A, B, C, D...
                text: text,
              })
            );
          } catch (error) {
            console.error("解析选项失败:", error);
            processedQuestion.options = [];
          }
        } else if (Array.isArray(question.options)) {
          // 如果已经是数组，确保每个选项都有id和text属性
          processedQuestion.options = question.options.map(
            (option: any, index: number) => {
              if (typeof option === "string") {
                return {
                  id: String.fromCharCode(65 + index),
                  text: option,
                };
              }
              return {
                id: option.id || String.fromCharCode(65 + index),
                text: option.text || option,
              };
            }
          );
        }
      }

      return processedQuestion;
    });

    // 设置练习信息
    practiceTitle.value = (route.query.title as string) || "练习题目";
    practiceSubject.value = (route.query.subject as string) || "数据结构";
    practiceDifficulty.value = (route.query.difficulty as string) || "medium";

    // 初始化填空题答案数组
    initFillAnswers();

    // 从本地存储恢复进度
    loadProgressFromLocal();
  } catch (error) {
    console.error("加载练习数据失败:", error);
    ElMessage.error("加载练习数据失败，请重试");
  } finally {
    loading.value = false;
  }
};

// 获取难度文本
const getDifficultyText = (difficulty: string): string => {
  const textMap: Record<string, string> = {
    easy: "简单",
    medium: "中等",
    hard: "困难",
  };
  return textMap[difficulty] || "";
};

// 获取类型文本
const getTypeText = (type: string): string => {
  const textMap: Record<string, string> = {
    single: "单选题",
    multiple: "多选题",
    judge: "判断题",
    fill: "填空题",
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
  };
  return typeMap[type] || "";
};

// 初始化填空题答案
const initFillAnswers = () => {
  const currentQ = currentQuestion.value;
  if (currentQ && currentQ.type === "fill") {
    const blankCount = currentQ.blanks?.length || 1;
    fillAnswers.value = new Array(blankCount).fill("");
  }
};

// 开始进度保存定时器
const startProgressSaveTimer = () => {
  // 每60秒自动保存一次进度
  progressSaveTimer.value = setInterval(() => {
    saveProgressToBackend();
  }, 60000);
};

// 保存进度到后端
const saveProgressToBackend = async () => {
  try {
    const progressData = {
      current_question_index: currentQuestionIndex.value,
      answers: answers,
      timestamp: Date.now(),
    };

    // 保存到本地存储
    localStorage.setItem(
      `practice_progress_${practiceId.value}`,
      JSON.stringify(progressData)
    );

    // 这里可以添加后端API调用来保存进度
    // await savePracticeProgress(practiceId.value, progressData)

    lastSaveTime.value = Date.now();
    localOperationCount.value = 0;

    console.log("进度已保存");
  } catch (error) {
    console.error("保存进度失败:", error);
  }
};

// 从本地存储加载进度
const loadProgressFromLocal = () => {
  try {
    const savedProgress = localStorage.getItem(`practice_progress_${practiceId.value}`);
    if (savedProgress) {
      const progressData = JSON.parse(savedProgress);
      currentQuestionIndex.value = progressData.current_question_index || 0;
      Object.assign(answers, progressData.answers || {});
      initFillAnswers();
      console.log("已恢复练习进度");
    }
  } catch (error) {
    console.error("加载本地进度失败:", error);
  }
};

// 增加本地操作计数
const incrementLocalOperation = () => {
  localOperationCount.value++;
  // 如果本地操作超过2次，立即保存进度
  if (localOperationCount.value >= 2) {
    saveProgressToBackend();
  }
};

// 跳转到指定题目
const goToQuestion = (index: number) => {
  currentQuestionIndex.value = index;
  initFillAnswers();
  incrementLocalOperation();
  // 切换题目时隐藏反馈
  answerFeedback.value.showFeedback = false;
};

// 上一题
const previousQuestion = () => {
  if (currentQuestionIndex.value > 0) {
    currentQuestionIndex.value--;
    initFillAnswers();
    incrementLocalOperation();
    // 切换题目时隐藏反馈
    answerFeedback.value.showFeedback = false;
  }
};
 




// 下一题
 const nextQuestion = () => {
   if (currentQuestionIndex.value < questions.value.length - 1) {
     currentQuestionIndex.value++;
     initFillAnswers();
     incrementLocalOperation();
     // 重置反馈状态
     answerFeedback.value.showFeedback = false;
   }
 };

// 选择选项（单选）
const selectOption = (optionId: string) => {
  answers[currentQuestion.value.id] = optionId;
  // 重置反馈状态，允许重新提交
  if (answerFeedback.value.questionId === currentQuestion.value.id) {
    answerFeedback.value.showFeedback = false;
  }
  incrementLocalOperation();
};

// 切换选项（多选）
const toggleOption = (optionId: string) => {
  const questionId = currentQuestion.value.id;
  if (!answers[questionId]) {
    answers[questionId] = [];
  }

  const selectedOptions = answers[questionId] as string[];
  const index = selectedOptions.indexOf(optionId);

  if (index > -1) {
    selectedOptions.splice(index, 1);
  } else {
    selectedOptions.push(optionId);
  }

  // 重置反馈状态，允许重新提交
  if (answerFeedback.value.questionId === questionId) {
    answerFeedback.value.showFeedback = false;
  }
  incrementLocalOperation();
};

// 判断选项是否被选中
const isOptionSelected = (optionId: string): boolean => {
  if (!currentQuestion.value) return false;
  const answer = answers[currentQuestion.value.id];
  if (currentQuestion.value.type === "single") {
    return answer === optionId;
  } else {
    return Array.isArray(answer) && answer.includes(optionId);
  }
};

// 选择判断题答案
const selectJudge = (value: boolean) => {
  answers[currentQuestion.value.id] = value;
  // 重置反馈状态，允许重新提交
  if (answerFeedback.value.questionId === currentQuestion.value.id) {
    answerFeedback.value.showFeedback = false;
  }
  incrementLocalOperation();
};

// 更新填空题答案
const updateFillAnswer = (index: number) => {
  const questionId = currentQuestion.value.id;
  answers[questionId] = [...fillAnswers.value];
  // 重置反馈状态，允许重新提交
  if (answerFeedback.value.questionId === questionId) {
    answerFeedback.value.showFeedback = false;
  }
  incrementLocalOperation();
};

// 确认答案
const confirmAnswer = async () => {
  if (!currentQuestion.value || !hasAnswered.value || submittingAnswer.value) {
    return;
  }
  
  submittingAnswer.value = true;
  
  try {
    const questionId = currentQuestion.value.id;
    let answer = '';
    
    if (currentQuestion.value.type === 'single' || currentQuestion.value.type === 'judge') {
      // 对于判断题，需要特殊处理boolean值，确保false也能正确转换为字符串
      const answerValue = answers[questionId];
      if (currentQuestion.value.type === 'judge' && typeof answerValue === 'boolean') {
        answer = String(answerValue);
      } else {
        answer = String(answerValue || '');
      }
    } else if (currentQuestion.value.type === 'multiple') {
      const selectedOptions = answers[questionId] || [];
      answer = Array.isArray(selectedOptions) ? selectedOptions.join(',') : '';
    } else if (currentQuestion.value.type === 'fill') {
      answer = fillAnswers.value.filter(ans => ans && ans.trim()).join('|');
    }
    
    // 构造正确的请求格式
    const requestData = {
      question_id: questionId,
      answer: answer,
      time_spent: Math.floor((Date.now() - lastSaveTime.value) / 1000)
    };
    
    const response = await submitPracticeAnswer(practiceId.value, requestData);
     
     // 更新反馈状态
     answerFeedback.value = {
       questionId,
       isCorrect: response.is_correct,
       explanation: response.explanation || "",
       score: response.score || 0,
       showFeedback: true,
     };
     
     // 保存进度
     incrementLocalOperation();
  } catch (error) {
    console.error("提交答案失败:", error);
    ElMessage.error("提交答案失败，请重试");
  } finally {
    submittingAnswer.value = false;
  }
};



// 返回练习页面
const goBack = () => {
  router.push({ name: "PracticeView" });
};

// 提交练习
const submitPractice = async () => {
  try {
    await ElMessageBox.confirm("确定要提交练习吗？提交后将无法修改答案。", "确认提交", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    loading.value = true;

    // 最后保存一次进度
    await saveProgressToBackend();

    // 调用后端API完成练习
    const result = await completePractice(practiceId.value);

    // 清除本地存储的进度
    localStorage.removeItem(`practice_progress_${practiceId.value}`);

    // 清除定时器
    if (progressSaveTimer.value) {
      clearInterval(progressSaveTimer.value);
      progressSaveTimer.value = null;
    }

    ElMessage.success(
      `练习完成！得分：${result.score}/${
        result.total_questions * 5
      }，正确率：${result.accuracy.toFixed(1)}%`
    );

    // 跳转回练习页面
    router.push({ name: "PracticeView" });
  } catch (error) {
    if (error !== "cancel") {
      console.error("提交练习失败:", error);
      ElMessage.error("提交练习失败，请重试");
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.practice-session {
  padding: var(--spacing-md);
}

.loading-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  z-index: 1000;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

/* 底部按钮区域样式 */
.bottom-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding: 16px 0;
}

.nav-button {
  min-width: 100px;
}

.confirm-button {
  min-width: 120px;
}

/* 答题后导航按钮样式 */
.question-navigation {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}

.next-button {
  min-width: 100px;
}

.practice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--gradient-primary);
  border-radius: var(--radius-md);
  color: white;
}

.practice-title {
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 var(--spacing-xs) 0;
}

.practice-meta {
  display: flex;
  gap: var(--spacing-md);
  font-size: 13px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.practice-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.question-nav {
  margin-bottom: var(--spacing-md);
  padding: var(--spacing-md);
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

.nav-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.nav-title-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  flex: 1;
}

.nav-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
}

.nav-progress {
  margin-bottom: var(--spacing-sm);
}

.progress-info {
  display: flex;
  gap: var(--spacing-md);
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
}

.progress-text {
  color: var(--text-primary);
}

.progress-percentage {
  color: var(--dopamine-blue);
  font-weight: 600;
}

/* 左右布局容器 */
.nav-layout {
  display: flex;
  gap: var(--spacing-md);
  height: 220px;
}

/* 左侧信息面板 */
.nav-info-panel {
  width: 240px;
  flex-shrink: 0;
  background: var(--bg-color);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
  overflow: hidden;
}

/* 右侧题目导航面板 */
.nav-grid-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.nav-grid-panel .nav-legend {
  background: rgba(255, 193, 7, 0.05);
  border-radius: var(--radius-sm);
  padding: var(--spacing-xs) var(--spacing-sm);
  display: flex;
  gap: var(--spacing-md);
  justify-content: center;
  flex-wrap: wrap;
}

.grid-title {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 统计标题 */
.stats-title {
  margin: 0 0 var(--spacing-sm) 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: var(--spacing-xs);
}

/* 进度统计 */
.progress-stats {
  background: rgba(64, 158, 255, 0.05);
  border-radius: var(--radius-sm);
  padding: var(--spacing-sm);
}

.stats-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-xs);
  font-size: 13px;
}

.stats-item:last-child {
  margin-bottom: 0;
}

.stats-label {
  color: var(--text-secondary);
}

.stats-value {
  font-weight: 600;
  color: var(--text-primary);
}

.answered-count {
  color: var(--dopamine-green);
}

.unanswered-count {
  color: var(--text-secondary);
}

.progress-percent {
  color: var(--dopamine-blue);
}

/* 题型统计 */
.type-stats {
  background: rgba(103, 194, 58, 0.05);
  border-radius: var(--radius-sm);
  padding: var(--spacing-sm);
}

.type-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.type-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.type-label {
  color: var(--text-secondary);
}

.type-count {
  font-weight: 600;
  color: var(--dopamine-green);
}

/* 图例样式 */
.nav-legend {
  margin-bottom: var(--spacing-sm);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
  font-size: 13px;
}

.legend-item:last-child {
  margin-bottom: 0;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-dot.answered {
  background: var(--dopamine-green);
}

.legend-dot.current {
  background: var(--dopamine-blue);
}

.legend-dot.unanswered {
  background: var(--border-color);
}

/* 当前题目信息样式已移除 */

/* 题目导航网格 */
.nav-grid {
  height: 170px;
  overflow-y: auto;
  position: relative;
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-color);
  flex: 1;
}

.nav-grid-content {
  position: relative;
  width: 100%;
}

/* 自定义滚动条样式 */
.nav-grid::-webkit-scrollbar {
  width: 6px;
}

.nav-grid::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 3px;
}

.nav-grid::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
  transition: background 0.3s ease;
}

.nav-grid::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

.nav-item {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
}

.nav-item.current {
  background: var(--dopamine-blue);
  border-color: var(--dopamine-blue);
  color: white;
}

.nav-item.answered {
  background: var(--dopamine-green);
  border-color: var(--dopamine-green);
  color: white;
}

.nav-item.unanswered:hover {
  border-color: var(--dopamine-blue);
}

.answer-area {
  background: white;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
}

.question-container {
  padding: var(--spacing-lg);
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-md);
  padding-bottom: var(--spacing-sm);
  border-bottom: 1px solid var(--border-color);
}

.question-number {
  font-size: 16px;
  font-weight: 700;
  color: var(--dopamine-blue);
}

.question-meta {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.question-score {
  font-weight: 600;
  color: var(--dopamine-orange);
}

.question-content {
  margin-bottom: var(--spacing-lg);
}

.question-text {
  font-size: 15px;
  line-height: 1.5;
  margin-bottom: var(--spacing-md);
  color: var(--text-primary);
}

.question-options {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.option-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.3s ease;
}

.option-item:hover {
  border-color: var(--dopamine-blue);
  background: rgba(64, 158, 255, 0.05);
}

.option-item.selected {
  border-color: var(--dopamine-blue);
  background: rgba(64, 158, 255, 0.1);
}

.option-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  flex: 1;
}

.option-prefix {
  font-weight: 600;
  color: var(--dopamine-blue);
}

.option-text {
  flex: 1;
}

.judge-options {
  display: flex;
  gap: var(--spacing-md);
}

.judge-option {
  flex: 1;
  padding: var(--spacing-md);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.judge-option:hover {
  border-color: var(--dopamine-blue);
  background: rgba(64, 158, 255, 0.05);
}

.judge-option.selected {
  border-color: var(--dopamine-blue);
  background: rgba(64, 158, 255, 0.1);
}

.fill-inputs {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.fill-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.fill-label {
  min-width: 70px;
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.question-navigation {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
}

/* 答题反馈样式 */
.answer-feedback {
  margin: var(--spacing-md) 0;
}

.feedback-card {
  padding: var(--spacing-md);
  border-radius: var(--radius-md);
  border: 2px solid;
  animation: fadeInUp 0.3s ease-out;
}

.feedback-card.correct {
  border-color: var(--dopamine-green);
  background: rgba(103, 194, 58, 0.1);
}

.feedback-card.incorrect {
  border-color: #f56565;
  background: rgba(245, 101, 101, 0.1);
}

.feedback-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.feedback-result {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.feedback-icon {
  font-size: 18px;
}

.feedback-card.correct .feedback-icon {
  color: var(--dopamine-green);
}

.feedback-card.incorrect .feedback-icon {
  color: #f56565;
}

.feedback-text {
  font-weight: 600;
  font-size: 15px;
}

.feedback-card.correct .feedback-text {
  color: var(--dopamine-green);
}

.feedback-card.incorrect .feedback-text {
  color: #f56565;
}

.feedback-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.feedback-score {
  font-weight: 600;
  color: var(--dopamine-orange);
}

.close-feedback-btn {
  padding: 4px;
  min-height: auto;
  color: var(--text-secondary);
  transition: color 0.3s ease;
}

.close-feedback-btn:hover {
  color: var(--text-primary);
}

.feedback-explanation {
  border-top: 1px solid var(--border-color);
  padding-top: var(--spacing-sm);
}

.explanation-title {
  font-weight: 600;
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.explanation-content {
  line-height: 1.5;
  color: var(--text-secondary);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .nav-layout {
    height: 170px;
  }
  
  .nav-info-panel {
    width: 200px;
  }
  
  .nav-grid {
    height: 120px;
  }
}

@media (max-width: 768px) {
  .nav-layout {
    flex-direction: column;
    height: auto;
    gap: var(--spacing-sm);
  }
  
  .nav-info-panel {
    width: 100%;
    height: auto;
    max-height: 120px;
  }
  
  .nav-grid {
    height: 70px;
  }
  
  .nav-item {
    width: 35px;
    height: 35px;
    font-size: 12px;
  }
  
  .practice-header {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  
  .practice-meta {
    flex-wrap: wrap;
  }
  
  .question-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-xs);
  }
  
  .stats-title {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .nav-info-panel {
    padding: var(--spacing-sm);
  }
  
  .nav-grid {
    height: 40px;
    padding: 6px;
  }
  
  .nav-item {
    width: 30px;
    height: 30px;
    font-size: 11px;
  }
  
  .practice-session {
    padding: var(--spacing-sm);
  }
  
  .question-container {
    padding: var(--spacing-md);
  }
  
  .bottom-actions {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  
  .nav-button,
  .confirm-button {
    width: 100%;
  }
  
  .stats-item,
  .type-item,
  .legend-item {
    font-size: 12px;
  }
}
</style>
