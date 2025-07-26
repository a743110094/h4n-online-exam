import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const practiceStartRate = new Rate('practice_start_errors');
const examStartRate = new Rate('exam_start_errors');

// 测试配置 - 专门针对核心业务接口的高强度测试
export const options = {
  stages: [
    { duration: '1m', target: 20 },  // 快速预热到20用户
    { duration: '3m', target: 50 },  // 3分钟内增加到50用户
    { duration: '15m', target: 50 }, // 保持50用户15分钟高强度测试
    { duration: '2m', target: 20 },  // 逐步降低
    { duration: '1m', target: 0 },   // 完全停止
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'],     // 95%请求响应时间<100ms
    http_req_failed: ['rate<0.05'],       // 总体错误率<5%
    errors: ['rate<0.05'],                // 自定义错误率<5%
    practice_start_errors: ['rate<0.02'], // 练习开始错误率<2%
    exam_start_errors: ['rate<0.02'],     // 考试开始错误率<2%
  },
};

// 基础配置
const BASE_URL = 'http://localhost:8080/api/v1';
const TENANT_ID = '100';

// 测试数据
const STUDENT_USERS = [
  { email: 'student1@test.com', password: 'student123' },
  { email: 'student2@test.com', password: 'student123' },
  { email: 'student3@test.com', password: 'student123' },
  { email: 'student@test.com', password: 'student123' },
];

const TEACHER_USER = { email: 'teacher@test.com', password: 'teacher123' };

// 获取随机学生用户
function getRandomStudent() {
  return STUDENT_USERS[Math.floor(Math.random() * STUDENT_USERS.length)];
}

// 登录函数
function login(user) {
  const loginPayload = {
    email: user.email,
    password: user.password
  };
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': TENANT_ID,
    },
  };
  
  const response = http.post(`${BASE_URL}/auth/login`, JSON.stringify(loginPayload), params);
  
  const success = check(response, {
    'login successful': (r) => r.status === 200,
  });
  
  if (!success) {
    errorRate.add(1);
    return null;
  }
  
  try {
    const body = JSON.parse(response.body);
    return body.data.token;
  } catch (e) {
    errorRate.add(1);
    return null;
  }
}

// 获取认证头
function getAuthHeaders(token) {
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    'X-Tenant-ID': TENANT_ID,
  };
}

// 测试练习推荐接口（高频接口）
function testPracticeRecommendations(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/practice/recommendations`, params);
  
  const success = check(response, {
    'practice recommendations status 200': (r) => r.status === 200,
    'practice recommendations RT < 80ms': (r) => r.timings.duration < 80,
  });
  
  if (!success) {
    errorRate.add(1);
  }
  
  return response;
}

// 测试开始练习接口（核心业务）
function testStartPractice(token) {
  const practicePayload = {
    subject_id: 1,
    question_count: 5,
    difficulty: 'medium'
  };
  
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.post(`${BASE_URL}/practice/start`, JSON.stringify(practicePayload), params);
  
  const success = check(response, {
    'start practice status 200': (r) => r.status === 200,
    'start practice RT < 150ms': (r) => r.timings.duration < 150,
    'start practice has data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.practice_id;
      } catch (e) {
        return false;
      }
    },
  });
  
  if (!success) {
    errorRate.add(1);
    practiceStartRate.add(1);
    return null;
  }
  
  try {
    const body = JSON.parse(response.body);
    return body.data.practice_id;
  } catch (e) {
    errorRate.add(1);
    return null;
  }
}

// 测试提交练习答案（高频接口）
function testSubmitPracticeAnswer(token, practiceId) {
  if (!practiceId) return;
  
  const answerPayload = {
    question_id: Math.floor(Math.random() * 100) + 1,
    answer: ['A'], // 模拟选择题答案
    time_spent: Math.floor(Math.random() * 30) + 10 // 10-40秒
  };
  
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.post(
    `${BASE_URL}/practice/${practiceId}/answer`, 
    JSON.stringify(answerPayload), 
    params
  );
  
  const success = check(response, {
    'submit practice answer status 200': (r) => r.status === 200,
    'submit practice answer RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试练习历史接口
function testPracticeHistory(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/practice/history?page=1&limit=10`, params);
  
  const success = check(response, {
    'practice history status 200': (r) => r.status === 200,
    'practice history RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试练习统计接口
function testPracticeStats(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/practice/stats`, params);
  
  const success = check(response, {
    'practice stats status 200': (r) => r.status === 200,
    'practice stats RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试学生考试列表
function testStudentExams(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/exams/student`, params);
  
  const success = check(response, {
    'student exams status 200': (r) => r.status === 200,
    'student exams RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
  
  return response;
}

// 测试开始考试接口
function testStartExam(token) {
  // 先获取可用考试
  const examsResponse = testStudentExams(token);
  
  try {
    const examsBody = JSON.parse(examsResponse.body);
    if (examsBody.data && examsBody.data.length > 0) {
      const examId = examsBody.data[0].id;
      
      const params = { headers: getAuthHeaders(token) };
      const response = http.post(`${BASE_URL}/exams/${examId}/start`, '{}', params);
      
      const success = check(response, {
        'start exam status 200 or 400': (r) => r.status === 200 || r.status === 400, // 400可能是已经开始
        'start exam RT < 150ms': (r) => r.timings.duration < 150,
      });
      
      if (!success) {
        errorRate.add(1);
        examStartRate.add(1);
      }
    }
  } catch (e) {
    // 忽略解析错误
  }
}

// 测试题目查询接口
function testQuestions(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/questions?page=1&limit=20&subject_id=1`, params);
  
  const success = check(response, {
    'questions status 200': (r) => r.status === 200,
    'questions RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试科目接口
function testSubjects(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/subjects`, params);
  
  const success = check(response, {
    'subjects status 200': (r) => r.status === 200,
    'subjects RT < 50ms': (r) => r.timings.duration < 50,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试用户资料接口
function testUserProfile(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/user/profile`, params);
  
  const success = check(response, {
    'user profile status 200': (r) => r.status === 200,
    'user profile RT < 80ms': (r) => r.timings.duration < 80,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 主测试函数 - 模拟真实学生使用场景
export default function () {
  // 随机选择学生用户登录
  const student = getRandomStudent();
  const token = login(student);
  
  if (!token) {
    sleep(1);
    return;
  }
  
  // 模拟学生的典型使用流程
  
  // 1. 查看个人资料（低频）
  if (Math.random() < 0.3) {
    testUserProfile(token);
  }
  
  // 2. 查看科目列表（中频）
  if (Math.random() < 0.6) {
    testSubjects(token);
  }
  
  // 3. 查看练习推荐（高频）
  testPracticeRecommendations(token);
  
  // 4. 开始练习（核心业务流程）
  if (Math.random() < 0.7) {
    const practiceId = testStartPractice(token);
    
    if (practiceId) {
      // 模拟答题过程
      const answerCount = Math.floor(Math.random() * 3) + 1; // 1-3题
      for (let i = 0; i < answerCount; i++) {
        testSubmitPracticeAnswer(token, practiceId);
        sleep(0.5); // 答题间隔
      }
    }
  }
  
  // 5. 查看练习历史（中频）
  if (Math.random() < 0.5) {
    testPracticeHistory(token);
  }
  
  // 6. 查看练习统计（中频）
  if (Math.random() < 0.4) {
    testPracticeStats(token);
  }
  
  // 7. 查看题目（中频）
  if (Math.random() < 0.5) {
    testQuestions(token);
  }
  
  // 8. 考试相关操作（低频但重要）
  if (Math.random() < 0.2) {
    testStartExam(token);
  }
  
  // 模拟用户思考时间
  sleep(Math.random() * 2 + 0.5); // 0.5-2.5秒
}

// 设置函数 - 测试开始前的准备
export function setup() {
  console.log('Starting core APIs load test...');
  console.log('Target: 50 VU, 95% RT < 100ms');
  return {};
}

// 清理函数 - 测试结束后的清理
export function teardown(data) {
  console.log('Core APIs load test completed');
}