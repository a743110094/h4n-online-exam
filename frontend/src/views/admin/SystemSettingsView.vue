<template>
  <div class="system-settings-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <p class="page-subtitle">管理系统配置和参数</p>
    </div>

    <!-- 设置选项卡 -->
    <el-tabs v-model="activeTab" class="settings-tabs">
      <!-- 基本设置 -->
      <el-tab-pane label="基本设置" name="basic">
        <div class="settings-section dopamine-card">
          <h3>系统信息</h3>
          <el-form :model="basicSettings" label-width="120px">
            <el-form-item label="系统名称">
              <el-input v-model="basicSettings.systemName" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统版本">
              <el-input v-model="basicSettings.version" readonly />
            </el-form-item>
            <el-form-item label="系统描述">
              <el-input
                v-model="basicSettings.description"
                type="textarea"
                :rows="3"
                placeholder="请输入系统描述"
              />
            </el-form-item>
            <el-form-item label="管理员邮箱">
              <el-input v-model="basicSettings.adminEmail" placeholder="请输入管理员邮箱" />
            </el-form-item>
            <el-form-item label="系统Logo">
              <el-upload
                class="logo-uploader"
                action="#"
                :show-file-list="false"
                :before-upload="beforeLogoUpload"
              >
                <img v-if="basicSettings.logo" :src="basicSettings.logo" class="logo" />
                <el-icon v-else class="logo-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 安全设置 -->
      <el-tab-pane label="安全设置" name="security">
        <div class="settings-section dopamine-card">
          <h3>密码策略</h3>
          <el-form :model="securitySettings" label-width="150px">
            <el-form-item label="最小密码长度">
              <el-input-number v-model="securitySettings.minPasswordLength" :min="6" :max="20" />
            </el-form-item>
            <el-form-item label="密码复杂度">
              <el-checkbox-group v-model="securitySettings.passwordComplexity">
                <el-checkbox label="uppercase">包含大写字母</el-checkbox>
                <el-checkbox label="lowercase">包含小写字母</el-checkbox>
                <el-checkbox label="number">包含数字</el-checkbox>
                <el-checkbox label="special">包含特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="密码有效期">
              <el-input-number v-model="securitySettings.passwordExpireDays" :min="0" :max="365" />
              <span style="margin-left: 8px;">天（0表示永不过期）</span>
            </el-form-item>
            <el-form-item label="登录失败锁定">
              <el-switch v-model="securitySettings.loginLockEnabled" />
            </el-form-item>
            <el-form-item label="最大失败次数" v-if="securitySettings.loginLockEnabled">
              <el-input-number v-model="securitySettings.maxLoginAttempts" :min="3" :max="10" />
            </el-form-item>
            <el-form-item label="锁定时间" v-if="securitySettings.loginLockEnabled">
              <el-input-number v-model="securitySettings.lockDuration" :min="5" :max="60" />
              <span style="margin-left: 8px;">分钟</span>
            </el-form-item>
          </el-form>
        </div>

        <div class="settings-section dopamine-card">
          <h3>会话管理</h3>
          <el-form :model="securitySettings" label-width="150px">
            <el-form-item label="会话超时时间">
              <el-input-number v-model="securitySettings.sessionTimeout" :min="30" :max="480" />
              <span style="margin-left: 8px;">分钟</span>
            </el-form-item>
            <el-form-item label="强制单点登录">
              <el-switch v-model="securitySettings.singleSignOn" />
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 考试设置 -->
      <el-tab-pane label="考试设置" name="exam">
        <div class="settings-section dopamine-card">
          <h3>考试规则</h3>
          <el-form :model="examSettings" label-width="150px">
            <el-form-item label="默认考试时长">
              <el-input-number v-model="examSettings.defaultDuration" :min="30" :max="300" />
              <span style="margin-left: 8px;">分钟</span>
            </el-form-item>
            <el-form-item label="允许暂停考试">
              <el-switch v-model="examSettings.allowPause" />
            </el-form-item>
            <el-form-item label="自动提交">
              <el-switch v-model="examSettings.autoSubmit" />
            </el-form-item>
            <el-form-item label="防作弊模式">
              <el-switch v-model="examSettings.antiCheat" />
            </el-form-item>
            <el-form-item label="随机题目顺序">
              <el-switch v-model="examSettings.randomOrder" />
            </el-form-item>
            <el-form-item label="显示题目编号">
              <el-switch v-model="examSettings.showQuestionNumber" />
            </el-form-item>
          </el-form>
        </div>

        <div class="settings-section dopamine-card">
          <h3>成绩设置</h3>
          <el-form :model="examSettings" label-width="150px">
            <el-form-item label="及格分数">
              <el-input-number v-model="examSettings.passingScore" :min="0" :max="100" />
              <span style="margin-left: 8px;">分</span>
            </el-form-item>
            <el-form-item label="立即显示成绩">
              <el-switch v-model="examSettings.showScoreImmediately" />
            </el-form-item>
            <el-form-item label="显示正确答案">
              <el-switch v-model="examSettings.showCorrectAnswer" />
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 邮件设置 -->
      <el-tab-pane label="邮件设置" name="email">
        <div class="settings-section dopamine-card">
          <h3>SMTP配置</h3>
          <el-form :model="emailSettings" label-width="120px">
            <el-form-item label="SMTP服务器">
              <el-input v-model="emailSettings.smtpHost" placeholder="请输入SMTP服务器地址" />
            </el-form-item>
            <el-form-item label="端口">
              <el-input-number v-model="emailSettings.smtpPort" :min="1" :max="65535" />
            </el-form-item>
            <el-form-item label="加密方式">
              <el-select v-model="emailSettings.encryption">
                <el-option label="无" value="none" />
                <el-option label="SSL" value="ssl" />
                <el-option label="TLS" value="tls" />
              </el-select>
            </el-form-item>
            <el-form-item label="用户名">
              <el-input v-model="emailSettings.username" placeholder="请输入邮箱用户名" />
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="emailSettings.password" type="password" placeholder="请输入邮箱密码" />
            </el-form-item>
            <el-form-item label="发件人名称">
              <el-input v-model="emailSettings.fromName" placeholder="请输入发件人名称" />
            </el-form-item>
            <el-form-item label="发件人邮箱">
              <el-input v-model="emailSettings.fromEmail" placeholder="请输入发件人邮箱" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="testEmail">测试邮件发送</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-tab-pane>

      <!-- 系统维护 -->
      <el-tab-pane label="系统维护" name="maintenance">
        <div class="settings-section dopamine-card">
          <h3>数据备份</h3>
          <div class="maintenance-actions">
            <div class="action-item">
              <div class="action-info">
                <h4>数据库备份</h4>
                <p>备份系统数据库，建议定期执行</p>
              </div>
              <el-button type="primary" @click="backupDatabase">立即备份</el-button>
            </div>
            <div class="action-item">
              <div class="action-info">
                <h4>文件备份</h4>
                <p>备份系统文件和用户上传的文件</p>
              </div>
              <el-button type="primary" @click="backupFiles">立即备份</el-button>
            </div>
          </div>
        </div>

        <div class="settings-section dopamine-card">
          <h3>系统清理</h3>
          <div class="maintenance-actions">
            <div class="action-item">
              <div class="action-info">
                <h4>清理临时文件</h4>
                <p>清理系统产生的临时文件和缓存</p>
              </div>
              <el-button type="warning" @click="cleanTempFiles">清理文件</el-button>
            </div>
            <div class="action-item">
              <div class="action-info">
                <h4>清理日志文件</h4>
                <p>清理过期的系统日志文件</p>
              </div>
              <el-button type="warning" @click="cleanLogs">清理日志</el-button>
            </div>
          </div>
        </div>

        <div class="settings-section dopamine-card">
          <h3>系统重启</h3>
          <div class="maintenance-actions">
            <div class="action-item">
              <div class="action-info">
                <h4>重启系统服务</h4>
                <p class="warning-text">⚠️ 重启会中断所有用户的操作，请谨慎执行</p>
              </div>
              <el-button type="danger" @click="restartSystem">重启系统</el-button>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 保存按钮 -->
    <div class="save-actions">
      <el-button type="primary" size="large" @click="saveSettings">
        保存设置
      </el-button>
      <el-button size="large" @click="resetSettings">
        重置设置
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 响应式数据
const activeTab = ref('basic')

// 基本设置
const basicSettings = reactive({
  systemName: '在线考试系统',
  version: 'v1.0.0',
  description: '基于Vue3和Element Plus的在线考试管理系统',
  adminEmail: 'admin@example.com',
  logo: ''
})

// 安全设置
const securitySettings = reactive({
  minPasswordLength: 8,
  passwordComplexity: ['lowercase', 'number'],
  passwordExpireDays: 90,
  loginLockEnabled: true,
  maxLoginAttempts: 5,
  lockDuration: 15,
  sessionTimeout: 120,
  singleSignOn: false
})

// 考试设置
const examSettings = reactive({
  defaultDuration: 120,
  allowPause: true,
  autoSubmit: true,
  antiCheat: false,
  randomOrder: false,
  showQuestionNumber: true,
  passingScore: 60,
  showScoreImmediately: true,
  showCorrectAnswer: false
})

// 邮件设置
const emailSettings = reactive({
  smtpHost: 'smtp.example.com',
  smtpPort: 587,
  encryption: 'tls',
  username: '',
  password: '',
  fromName: '在线考试系统',
  fromEmail: 'noreply@example.com'
})

// 方法
const beforeLogoUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('Logo只能是 JPG/PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('Logo大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

const testEmail = () => {
  ElMessage.success('邮件测试功能开发中')
}

const backupDatabase = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要备份数据库吗？此操作可能需要一些时间。',
      '确认备份',
      { type: 'info' }
    )
    ElMessage.success('数据库备份已开始，请稍候...')
  } catch {
    // 用户取消
  }
}

const backupFiles = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要备份文件吗？此操作可能需要一些时间。',
      '确认备份',
      { type: 'info' }
    )
    ElMessage.success('文件备份已开始，请稍候...')
  } catch {
    // 用户取消
  }
}

const cleanTempFiles = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理临时文件吗？此操作不可恢复。',
      '确认清理',
      { type: 'warning' }
    )
    ElMessage.success('临时文件清理完成')
  } catch {
    // 用户取消
  }
}

const cleanLogs = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理日志文件吗？此操作不可恢复。',
      '确认清理',
      { type: 'warning' }
    )
    ElMessage.success('日志文件清理完成')
  } catch {
    // 用户取消
  }
}

const restartSystem = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重启系统吗？这将中断所有用户的操作！',
      '确认重启',
      { type: 'error' }
    )
    ElMessage.success('系统重启指令已发送')
  } catch {
    // 用户取消
  }
}

const saveSettings = () => {
  ElMessage.success('设置保存成功')
}

const resetSettings = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置所有设置吗？此操作将恢复默认配置。',
      '确认重置',
      { type: 'warning' }
    )
    ElMessage.success('设置已重置')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.system-settings-view {
  padding: var(--spacing-lg);
}

.page-header {
  margin-bottom: var(--spacing-xl);
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 var(--spacing-sm) 0;
}

.page-subtitle {
  color: var(--text-secondary);
  margin: 0;
}

.settings-tabs {
  margin-bottom: var(--spacing-xl);
}

.settings-section {
  padding: var(--spacing-xl);
  margin-bottom: var(--spacing-lg);
}

.settings-section h3 {
  margin: 0 0 var(--spacing-lg) 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.logo-uploader {
  border: 1px dashed var(--border-color);
  border-radius: var(--radius-md);
  width: 120px;
  height: 120px;
  text-align: center;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.logo-uploader:hover {
  border-color: var(--color-primary);
}

.logo {
  width: 120px;
  height: 120px;
  display: block;
  object-fit: cover;
}

.logo-uploader-icon {
  font-size: 28px;
  color: var(--text-placeholder);
  width: 120px;
  height: 120px;
  line-height: 120px;
}

.maintenance-actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.action-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-lg);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
}

.action-info h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.action-info p {
  margin: 0;
  color: var(--text-secondary);
  font-size: 14px;
}

.warning-text {
  color: var(--color-warning) !important;
}

.save-actions {
  display: flex;
  justify-content: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-xl) 0;
}
</style>