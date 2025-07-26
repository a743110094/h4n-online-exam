import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');

// 测试配置
export const options = {
  stages: [
    { duration: '2m', target: 10 }, // 预热阶段：2分钟内逐步增加到10个用户
    { duration: '5m', target: 50 }, // 压力阶段：5分钟内增加到50个用户
    { duration: '10m', target: 50 }, // 稳定阶段：保持50个用户10分钟
    { duration: '2m', target: 0 },  // 冷却阶段：2分钟内减少到0个用户
  ],
  thresholds: {
    http_req_duration: ['p(95)<100'], // 95%的请求响应时间小于100ms
    http_req_failed: ['rate<0.1'],    // 错误率小于10%
    errors: ['rate<0.1'],             // 自定义错误率小于10%
  },
};

// 基础配置
const BASE_URL = 'http://localhost:8080/api/v1';
const TENANT_ID = '100';

// 测试用户数据
const TEST_USERS = {
  admin: {
    email: 'admin@test.com',
    password: 'admin123',
    role: 'admin'
  },
  teacher: {
    email: 'teacher@test.com',
    password: 'teacher123',
    role: 'teacher'
  },
  student: {
    email: 'student@test.com',
    password: 'student123',
    role: 'student'
  }
};

// 获取随机用户
function getRandomUser() {
  const users = Object.values(TEST_USERS);
  return users[Math.floor(Math.random() * users.length)];
}

// 登录并获取token
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
    'login status is 200': (r) => r.status === 200,
    'login response has token': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.token;
      } catch (e) {
        return false;
      }
    },
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

// 测试健康检查接口
function testHealthCheck() {
  const response = http.get(`${BASE_URL}/health`);
  
  const success = check(response, {
    'health check status is 200': (r) => r.status === 200,
    'health check response time < 50ms': (r) => r.timings.duration < 50,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试科目接口
function testSubjects(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/subjects`, params);
  
  const success = check(response, {
    'get subjects status is 200': (r) => r.status === 200,
    'get subjects response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试题目接口
function testQuestions(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/questions?page=1&limit=10`, params);
  
  const success = check(response, {
    'get questions status is 200': (r) => r.status === 200,
    'get questions response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试试卷接口
function testPapers(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/papers`, params);
  
  const success = check(response, {
    'get papers status is 200': (r) => r.status === 200,
    'get papers response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试考试接口
function testExams(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/exams`, params);
  
  const success = check(response, {
    'get exams status is 200': (r) => r.status === 200,
    'get exams response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试学生考试列表
function testStudentExams(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/exams/student`, params);
  
  const success = check(response, {
    'get student exams status is 200': (r) => r.status === 200,
    'get student exams response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试练习推荐接口
function testPracticeRecommendations(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/practice/recommendations`, params);
  
  const success = check(response, {
    'get practice recommendations status is 200': (r) => r.status === 200,
    'get practice recommendations response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试练习历史
function testPracticeHistory(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/practice/history`, params);
  
  const success = check(response, {
    'get practice history status is 200': (r) => r.status === 200,
    'get practice history response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试用户资料接口
function testUserProfile(token) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  const response = http.get(`${BASE_URL}/user/profile`, params);
  
  const success = check(response, {
    'get user profile status is 200': (r) => r.status === 200,
    'get user profile response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 测试统计接口
function testStats(token, userRole) {
  const params = {
    headers: getAuthHeaders(token),
  };
  
  let endpoint;
  if (userRole === 'student') {
    endpoint = `${BASE_URL}/stats/student`;
  } else if (userRole === 'teacher') {
    endpoint = `${BASE_URL}/stats/teacher`;
  } else {
    return; // admin没有专门的stats接口
  }
  
  const response = http.get(endpoint, params);
  
  const success = check(response, {
    'get stats status is 200': (r) => r.status === 200,
    'get stats response time < 100ms': (r) => r.timings.duration < 100,
  });
  
  if (!success) {
    errorRate.add(1);
  }
}

// 主测试函数
export default function () {
  // 1. 健康检查（无需认证）
  testHealthCheck();
  
  // 2. 随机选择用户进行登录
  const user = getRandomUser();
  const token = login(user);
  
  if (!token) {
    console.log('Login failed, skipping authenticated tests');
    sleep(1);
    return;
  }
  
  // 3. 测试需要认证的接口
  testUserProfile(token);
  testSubjects(token);
  testQuestions(token);
  testPapers(token);
  testExams(token);
  
  // 4. 根据用户角色测试特定接口
  if (user.role === 'student') {
    testStudentExams(token);
    testPracticeRecommendations(token);
    testPracticeHistory(token);
  }
  
  testStats(token, user.role);
  
  // 5. 随机休眠1-3秒，模拟真实用户行为
  sleep(Math.random() * 2 + 1);
}

// 测试结束后的清理函数
export function teardown(data) {
  console.log('Load test completed');
}