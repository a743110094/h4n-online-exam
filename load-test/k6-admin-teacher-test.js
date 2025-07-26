import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const createQuestionRate = new Rate('create_question_errors');
const createExamRate = new Rate('create_exam_errors');

// 测试配置 - 针对教师和管理员接口
export const options = {
  stages: [
    { duration: '1m', target: 5 },   // 预热到5用户
    { duration: '3m', target: 15 },  // 增加到15用户
    { duration: '10m', target: 15 }, // 保持15用户10分钟
    { duration: '2m', target: 0 },   // 降到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<150'],     // 管理接口允许稍高的响应时间
    http_req_failed: ['rate<0.05'],       // 错误率<5%
    errors: ['rate<0.05'],
    create_question_errors: ['rate<0.02'],
    create_exam_errors: ['rate<0.02'],
  },
};

// 基础配置
const BASE_URL = 'http://localhost:8080/api/v1';
const TENANT_ID = '100';

// 测试用户
const ADMIN_USER = { email: 'admin@test.com', password: 'admin123' };
const TEACHER_USER = { email: 'teacher@test.com', password: 'teacher123' };

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
    'admin/teacher login successful': (r) => r.status === 200,
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

// 测试管理员仪表板
function testAdminDashboard(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/admin/dashboard`, params);
  
  const success = check(response, {
    'admin dashboard status 200': (r) => r.status === 200,
    'admin dashboard RT < 200ms': (r) => r.timings.duration < 200,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试用户管理接口
function testUserManagement(token) {
  const params = { headers: getAuthHeaders(token) };
  
  // 获取用户列表
  const response = http.get(`${BASE_URL}/admin/users?page=1&limit=20`, params);
  
  const success = check(response, {
    'get users status 200': (r) => r.status === 200,
    'get users RT < 150ms': (r) => r.timings.duration < 150,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试科目管理
function testSubjectManagement(token) {
  const params = { headers: getAuthHeaders(token) };
  
  // 获取科目列表
  const getResponse = http.get(`${BASE_URL}/subjects/all`, params);
  
  const success = check(getResponse, {
    'get all subjects status 200': (r) => r.status === 200,
    'get all subjects RT < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试题目管理（教师）
function testQuestionManagement(token) {
  const params = { headers: getAuthHeaders(token) };
  
  // 获取题目列表
  const getResponse = http.get(`${BASE_URL}/questions?page=1&limit=20`, params);
  
  const getSuccess = check(getResponse, {
    'get questions for teacher status 200': (r) => r.status === 200,
    'get questions for teacher RT < 120ms': (r) => r.timings.duration < 120,
  });
  
  if (!getSuccess) {
    errorRate.add(1);
  }
  
  // 创建题目（低频操作）
  if (Math.random() < 0.1) {
    const questionPayload = {
      subject_id: 1,
      type: 'single_choice',
      difficulty: 'medium',
      content: `测试题目 ${Date.now()}`,
      options: ['选项A', '选项B', '选项C', '选项D'],
      correct_answer: ['A'],
      explanation: '这是一道测试题目',
      score: 5
    };
    
    const createResponse = http.post(
      `${BASE_URL}/teacher/questions`, 
      JSON.stringify(questionPayload), 
      params
    );
    
    const createSuccess = check(createResponse, {
      'create question status 200': (r) => r.status === 200,
      'create question RT < 200ms': (r) => r.timings.duration < 200,
    });
    
    if (!createSuccess) {
      errorRate.add(1);
      createQuestionRate.add(1);
    }
  }
}

// 测试试卷管理
function testPaperManagement(token) {
  const params = { headers: getAuthHeaders(token) };
  
  // 获取试卷列表
  const getResponse = http.get(`${BASE_URL}/papers`, params);
  
  const getSuccess = check(getResponse, {
    'get papers for teacher status 200': (r) => r.status === 200,
    'get papers for teacher RT < 120ms': (r) => r.timings.duration < 120,
  });
  
  if (!getSuccess) {
    errorRate.add(1);
  }
  
  // 自动组卷（低频但重要操作）
  if (Math.random() < 0.05) {
    const autoPaperPayload = {
      name: `自动组卷测试 ${Date.now()}`,
      subject_id: 1,
      total_score: 100,
      question_config: {
        single_choice: { count: 10, score: 5 },
        multiple_choice: { count: 5, score: 10 }
      },
      difficulty_distribution: {
        easy: 0.3,
        medium: 0.5,
        hard: 0.2
      }
    };
    
    const autoResponse = http.post(
      `${BASE_URL}/teacher/papers/auto`, 
      JSON.stringify(autoPaperPayload), 
      params
    );
    
    const autoSuccess = check(autoResponse, {
      'auto create paper status 200': (r) => r.status === 200,
      'auto create paper RT < 300ms': (r) => r.timings.duration < 300,
    });
    
    if (!autoSuccess) {
      errorRate.add(1);
    }
  }
}

// 测试考试管理
function testExamManagement(token) {
  const params = { headers: getAuthHeaders(token) };
  
  // 获取考试列表
  const getResponse = http.get(`${BASE_URL}/exams`, params);
  
  const getSuccess = check(getResponse, {
    'get exams for teacher status 200': (r) => r.status === 200,
    'get exams for teacher RT < 120ms': (r) => r.timings.duration < 120,
  });
  
  if (!getSuccess) {
    errorRate.add(1);
  }
  
  // 创建考试（低频操作）
  if (Math.random() < 0.08) {
    const examPayload = {
      title: `压测考试 ${Date.now()}`,
      paper_id: 1,
      start_time: new Date(Date.now() + 3600000).toISOString(), // 1小时后
      end_time: new Date(Date.now() + 7200000).toISOString(),   // 2小时后
      duration: 60,
      description: '这是一个压力测试创建的考试'
    };
    
    const createResponse = http.post(
      `${BASE_URL}/teacher/exams`, 
      JSON.stringify(examPayload), 
      params
    );
    
    const createSuccess = check(createResponse, {
      'create exam status 200': (r) => r.status === 200,
      'create exam RT < 200ms': (r) => r.timings.duration < 200,
    });
    
    if (!createSuccess) {
      errorRate.add(1);
      createExamRate.add(1);
    }
  }
}

// 测试统计接口
function testTeacherStats(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/stats/teacher`, params);
  
  const success = check(response, {
    'teacher stats status 200': (r) => r.status === 200,
    'teacher stats RT < 150ms': (r) => r.timings.duration < 150,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试题目统计
function testQuestionStats(token) {
  const params = { headers: getAuthHeaders(token) };
  const response = http.get(`${BASE_URL}/questions/stats`, params);
  
  const success = check(response, {
    'question stats status 200': (r) => r.status === 200,
    'question stats RT < 120ms': (r) => r.timings.duration < 120,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 主测试函数
export default function () {
  // 随机选择管理员或教师身份
  const isAdmin = Math.random() < 0.3; // 30%概率是管理员
  const user = isAdmin ? ADMIN_USER : TEACHER_USER;
  const token = login(user);
  
  if (!token) {
    sleep(1);
    return;
  }
  
  // 管理员专用操作
  if (isAdmin) {
    testAdminDashboard(token);
    
    if (Math.random() < 0.7) {
      testUserManagement(token);
    }
    
    if (Math.random() < 0.5) {
      testSubjectManagement(token);
    }
  }
  
  // 教师操作（管理员也可以执行）
  testQuestionManagement(token);
  
  if (Math.random() < 0.8) {
    testPaperManagement(token);
  }
  
  if (Math.random() < 0.6) {
    testExamManagement(token);
  }
  
  if (Math.random() < 0.4) {
    testTeacherStats(token);
  }
  
  if (Math.random() < 0.3) {
    testQuestionStats(token);
  }
  
  // 管理操作通常需要更多思考时间
  sleep(Math.random() * 3 + 1); // 1-4秒
}

// 设置函数
export function setup() {
  console.log('Starting admin/teacher APIs load test...');
  console.log('Target: 15 VU for management operations');
  return {};
}

// 清理函数
export function teardown(data) {
  console.log('Admin/teacher APIs load test completed');
}