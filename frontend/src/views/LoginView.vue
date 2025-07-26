<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="floating-circle circle-1"></div>
      <div class="floating-circle circle-2"></div>
      <div class="floating-circle circle-3"></div>
      <div class="floating-circle circle-4"></div>
      <div class="floating-circle circle-5"></div>
      <div class="floating-circle circle-6"></div>
    </div>

    <!-- 登录内容 -->
    <div class="login-content">
      <!-- 登录卡片 -->
      <div class="login-card">
        <div class="login-card-left">
          <!-- 登录头部 -->
          <div class="login-header">
            <h1 class="login-title">欢迎回来</h1>
            <p class="login-subtitle">登录您的账户继续学习之旅</p>
          </div>

          <!-- 登录表单 -->
          <form @submit.prevent="handleLogin" class="login-form">
            <div class="form-group">
              <input
                id="username"
                v-model="loginForm.username"
                type="text"
                class="form-input"
                placeholder="用户名"
                required
              />
            </div>

            <div class="form-group">
              <input
                id="password"
                v-model="loginForm.password"
                type="password"
                class="form-input"
                placeholder="密码"
                required
              />
            </div>

            <button type="submit" class="login-button" :disabled="authStore.isLoading">
              <span v-if="!authStore.isLoading">登录</span>
              <span v-else>登录中...</span>
            </button>
          </form>

          <!-- 注册链接 -->
          <div class="register-link">
            <span>还没有账户？</span>
            <button @click="showRegister = true" class="link-button">
              立即注册
            </button>
          </div>
        </div>

        <div class="login-card-right">
          <!-- 分割线 -->
          <div class="divider">
            <span class="divider-text">快速登录</span>
          </div>

          <!-- 快速登录 -->
          <div class="quick-login">
            <div class="quick-login-buttons">
              <button @click="quickLogin('admin')" class="quick-button admin">
                <span class="quick-button-icon">👨‍💼</span>
                <span class="quick-button-text">管理员</span>
              </button>
              <button @click="quickLogin('teacher')" class="quick-button teacher">
                <span class="quick-button-icon">👨‍🏫</span>
                <span class="quick-button-text">教师</span>
              </button>
              <button @click="quickLogin('student')" class="quick-button student">
                <span class="quick-button-icon">👨‍🎓</span>
                <span class="quick-button-text">学生</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
    
    <!-- 注册对话框 -->
    <div v-if="showRegister" class="dialog-overlay" @click="showRegister = false">
      <div class="dialog-content" @click.stop>
        <div class="dialog-header">
          <h2>用户注册</h2>
          <button @click="showRegister = false" class="close-button">×</button>
        </div>
        <form @submit.prevent="handleRegister" class="register-form">
          <div class="form-group">
            <input
              id="reg-username"
              v-model="registerForm.username"
              type="text"
              class="form-input"
              placeholder="用户名"
              required
            />
          </div>
          <div class="form-group">
            <input
              id="reg-email"
              v-model="registerForm.email"
              type="email"
              class="form-input"
              placeholder="邮箱"
              required
            />
          </div>
          <div class="form-group">
            <input
              id="reg-realName"
              v-model="registerForm.realName"
              type="text"
              class="form-input"
              placeholder="真实姓名"
              required
            />
          </div>
          <div class="form-group">
            <input
              id="reg-phone"
              v-model="registerForm.phone"
              type="tel"
              class="form-input"
              placeholder="手机号（可选）"
            />
          </div>
          <div class="form-group">
            <input
              id="reg-password"
              v-model="registerForm.password"
              type="password"
              class="form-input"
              placeholder="密码"
              required
            />
          </div>
          <div class="form-group">
            <input
              id="reg-confirm-password"
              v-model="registerForm.confirmPassword"
              type="password"
              class="form-input"
              placeholder="确认密码"
              required
            />
          </div>
          <div class="dialog-footer">
            <button type="button" @click="showRegister = false" class="cancel-button">取消</button>
            <button type="submit" class="register-button" :disabled="authStore.isLoading">
              <span v-if="!authStore.isLoading">注册</span>
              <span v-else>注册中...</span>
            </button>
          </div>
        </form>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'sonner'

const router = useRouter()
const authStore = useAuthStore()

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 注册相关
const showRegister = ref(false)
const registerForm = reactive({
  username: '',
  email: '',
  realName: '',
  phone: '',
  password: '',
  confirmPassword: ''
})

// 登录处理
const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    toast.error('请输入用户名和密码')
    return
  }
  
  try {
    const result = await authStore.login({
      username: loginForm.username,
      password: loginForm.password,
      remember: false
    })
    
    if (result.success) {
      toast.success('登录成功')
      // 根据用户角色跳转到对应页面
      const userRole = authStore.user?.role
      if (userRole === 'admin') {
        router.push('/admin')
      } else if (userRole === 'teacher') {
        router.push('/teacher')
      } else if (userRole === 'student') {
        router.push('/student')
      } else {
        router.push('/')
      }
    } else {
      toast.error(result.message || '登录失败')
    }
  } catch (error) {
    console.error('登录失败:', error)
    toast.error('登录失败，请稍后重试')
  }
}

// 快速登录
const quickLogin = async (role: 'admin' | 'teacher' | 'student') => {
  const credentials = {
    admin: { username: 'admin', password: 'password' },
    teacher: { username: 'teacher1', password: 'password' },
    student: { username: 'student1', password: 'password' }
  }
  
  const { username, password } = credentials[role]
  
  try {
    const result = await authStore.login({ username, password, remember: false })
    
    if (result.success) {
      toast.success(`${role === 'admin' ? '管理员' : role === 'teacher' ? '教师' : '学生'}登录成功`)
      router.push(`/${role}`)
    } else {
      toast.error(result.message || '登录失败')
    }
  } catch (error) {
    console.error('快速登录失败:', error)
    toast.error('登录失败，请稍后重试')
  }
}

// 注册处理
const handleRegister = async () => {
  if (!registerForm.username || !registerForm.email || !registerForm.password) {
    toast.error('请填写必填字段')
    return
  }
  
  if (registerForm.password !== registerForm.confirmPassword) {
    toast.error('两次输入密码不一致')
    return
  }
  
  try {
    const result = await authStore.register({
      username: registerForm.username,
      email: registerForm.email,
      realName: registerForm.realName,
      phone: registerForm.phone,
      password: registerForm.password,
      confirmPassword: registerForm.confirmPassword
    })
    
    if (result.success) {
      toast.success('注册成功，请登录')
      showRegister.value = false
      // 清空注册表单
      Object.assign(registerForm, {
        username: '',
        email: '',
        realName: '',
        phone: '',
        password: '',
        confirmPassword: ''
      })
    } else {
      toast.error(result.message || '注册失败')
    }
  } catch (error) {
    console.error('注册失败:', error)
    toast.error('注册失败，请稍后重试')
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.floating-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.circle-1 {
  width: 80px;
  height: 80px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.circle-2 {
  width: 120px;
  height: 120px;
  top: 20%;
  right: 15%;
  animation-delay: 1s;
}

.circle-3 {
  width: 60px;
  height: 60px;
  bottom: 30%;
  left: 20%;
  animation-delay: 2s;
}

.circle-4 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  right: 10%;
  animation-delay: 3s;
}

.circle-5 {
  width: 40px;
  height: 40px;
  top: 50%;
  left: 5%;
  animation-delay: 4s;
}

.circle-6 {
  width: 90px;
  height: 90px;
  top: 70%;
  right: 25%;
  animation-delay: 5s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.7;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 1;
  }
}

/* 登录内容 */
.login-content {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 640px;
  padding: 0 20px;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  width: 100%;
  max-width: 640px;
  height: 360px;
  aspect-ratio: 16/9;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 40px;
}

.login-card-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
}

.login-card-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
  padding-left: 20px;
  border-left: 1px solid rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-content {
    max-width: 420px;
  }
  
  .login-card {
    flex-direction: column;
    height: auto;
    aspect-ratio: auto;
    max-width: 420px;
    gap: 20px;
  }
  
  .login-card-left {
    height: auto;
  }
  
  .login-card-right {
    height: auto;
    padding-left: 0;
    border-left: none;
    border-top: 1px solid rgba(0, 0, 0, 0.1);
    padding-top: 20px;
  }
  
  .login-header {
    text-align: center;
    margin-bottom: 24px;
  }
  
  .register-link {
    text-align: center;
  }
  
  .quick-login-buttons {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
  }
  
  .quick-button {
    flex-direction: column;
    justify-content: center;
    height: 80px;
    padding: 8px;
    gap: 4px;
  }
  
  .quick-button-icon {
    font-size: 24px;
    margin-bottom: 4px;
  }
}

.login-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 35px 70px rgba(0, 0, 0, 0.2);
}

/* 登录头部 */
.login-header {
  text-align: left;
  margin-bottom: 24px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 6px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-subtitle {
  color: #7f8c8d;
  font-size: 14px;
  margin: 0;
  font-weight: 400;
}

/* 登录表单 */
.login-form {
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-input {
  width: 100%;
  height: 44px;
  padding: 0 16px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  font-size: 14px;
  background: rgba(255, 255, 255, 0.8);
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
  background: rgba(255, 255, 255, 1);
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
  color: #a0a0a0;
  font-weight: 400;
}

.login-button {
  width: 100%;
  height: 44px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 6px;
}

.login-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.login-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 分割线 */
.divider {
  position: relative;
  text-align: center;
  margin: 0 0 20px 0;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: rgba(0, 0, 0, 0.1);
}

.divider-text {
  background: rgba(255, 255, 255, 0.95);
  padding: 0 16px;
  color: #7f8c8d;
  font-size: 12px;
  font-weight: 500;
}

/* 快速登录 */
.quick-login {
  margin-bottom: 0;
}

.quick-login-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-button {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  height: 44px;
  padding: 0 16px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 12px;
  font-weight: 500;
  gap: 12px;
}

.quick-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.quick-button.admin {
  border-color: #ff6b6b;
  color: #ff6b6b;
}

.quick-button.admin:hover {
  background: #ff6b6b;
  color: white;
}

.quick-button.teacher {
  border-color: #4ecdc4;
  color: #4ecdc4;
}

.quick-button.teacher:hover {
  background: #4ecdc4;
  color: white;
}

.quick-button.student {
  border-color: #45b7d1;
  color: #45b7d1;
}

.quick-button.student:hover {
  background: #45b7d1;
  color: white;
}

.quick-button-icon {
  font-size: 18px;
  margin-bottom: 0;
}

.quick-button-text {
  font-size: 12px;
  font-weight: 500;
}

/* 注册链接 */
.register-link {
  text-align: left;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  color: #7f8c8d;
  font-size: 12px;
}

.link-button {
  background: none;
  border: none;
  color: #667eea;
  font-weight: 600;
  cursor: pointer;
  text-decoration: none;
  margin-left: 4px;
  transition: all 0.3s ease;
}

.link-button:hover {
  color: #764ba2;
  text-decoration: underline;
}

/* 注册对话框 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.dialog-content {
  background: white;
  border-radius: 24px;
  width: 90%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 24px 24px 0 0;
}

.dialog-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.close-button {
  background: none;
  border: none;
  color: white;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.2);
}

.register-form {
  padding: 32px;
}

.dialog-footer {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.cancel-button {
  flex: 1;
  height: 48px;
  border: 2px solid #e0e0e0;
  background: white;
  color: #666;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.cancel-button:hover {
  border-color: #ccc;
  background: #f5f5f5;
}

.register-button {
  flex: 2;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  color: white;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.register-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.3);
}

.register-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-content {
    padding: 0 16px;
  }
  
  .login-card {
    padding: 32px 24px;
  }
  
  .login-title {
    font-size: 28px;
  }
  
  .quick-login-buttons {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .quick-button {
    height: 60px;
    flex-direction: row;
    gap: 8px;
  }
  
  .quick-button-icon {
    font-size: 20px;
    margin-bottom: 0;
  }
  
  .dialog-content {
    margin: 20px;
    width: calc(100% - 40px);
  }
  
  .register-form {
    padding: 24px;
  }
  
  .dialog-header {
    padding: 20px 24px;
  }
}
</style>