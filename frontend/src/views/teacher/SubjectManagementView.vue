<template>
  <div class="subject-management-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">科目分类管理</h1>
      <p class="page-subtitle">管理学科分类和知识点体系</p>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar dopamine-card">
      <div class="toolbar-left">
        <el-button type="primary" @click="showAddSubjectDialog = true">
          <el-icon><Plus /></el-icon>
          新增科目
        </el-button>
        <el-button @click="expandAll">展开全部</el-button>
        <el-button @click="collapseAll">收起全部</el-button>
      </div>
      <div class="toolbar-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索科目或知识点"
          style="width: 300px"
          clearable
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>
    </div>

    <!-- 科目树形结构 -->
    <div class="subject-tree dopamine-card">
      <el-tree
        ref="treeRef"
        :data="subjectTree"
        :props="treeProps"
        :filter-node-method="filterNode"
        node-key="id"
        default-expand-all
        :expand-on-click-node="false"
        class="subject-tree-component"
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <div class="node-content">
              <el-icon v-if="data.type === 'subject'" class="node-icon subject-icon">
                <Folder />
              </el-icon>
              <el-icon v-else class="node-icon knowledge-icon">
                <Document />
              </el-icon>
              <span class="node-label">{{ data.name }}</span>
              <el-tag v-if="data.type === 'subject'" size="small" type="primary">
                {{ data.children?.length || 0 }}个知识点
              </el-tag>
              <el-tag v-else size="small" type="info">
                {{ data.questionCount || 0 }}道题目
              </el-tag>
            </div>
            <div class="node-actions">
              <el-button
                v-if="data.type === 'subject'"
                size="small"
                type="primary"
                text
                @click="addKnowledgePoint(data)"
              >
                添加知识点
              </el-button>
              <el-button size="small" type="primary" text @click="editNode(data)">
                编辑
              </el-button>
              <el-button size="small" type="danger" text @click="deleteNode(data)">
                删除
              </el-button>
            </div>
          </div>
        </template>
      </el-tree>
    </div>

    <!-- 新增科目对话框 -->
    <el-dialog v-model="showAddSubjectDialog" title="新增科目" width="500px">
      <el-form :model="subjectForm" :rules="subjectRules" ref="subjectFormRef" label-width="80px">
        <el-form-item label="科目名称" prop="name">
          <el-input v-model="subjectForm.name" placeholder="请输入科目名称" />
        </el-form-item>
        <el-form-item label="科目代码" prop="code">
          <el-input v-model="subjectForm.code" placeholder="请输入科目代码" />
        </el-form-item>
        <el-form-item label="科目描述">
          <el-input
            v-model="subjectForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入科目描述"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="subjectForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddSubjectDialog = false">取消</el-button>
        <el-button type="primary" @click="saveSubject">确定</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑知识点对话框 -->
    <el-dialog v-model="showKnowledgeDialog" :title="knowledgeDialogTitle" width="500px">
      <el-form :model="knowledgeForm" :rules="knowledgeRules" ref="knowledgeFormRef" label-width="80px">
        <el-form-item label="所属科目">
          <el-select v-model="knowledgeForm.subjectId" placeholder="请选择科目" disabled>
            <el-option
              v-for="subject in subjects"
              :key="subject.id"
              :label="subject.name"
              :value="subject.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="知识点名称" prop="name">
          <el-input v-model="knowledgeForm.name" placeholder="请输入知识点名称" />
        </el-form-item>
        <el-form-item label="知识点代码" prop="code">
          <el-input v-model="knowledgeForm.code" placeholder="请输入知识点代码" />
        </el-form-item>
        <el-form-item label="难度级别">
          <el-rate v-model="knowledgeForm.difficulty" show-text />
        </el-form-item>
        <el-form-item label="知识点描述">
          <el-input
            v-model="knowledgeForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入知识点描述"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="knowledgeForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showKnowledgeDialog = false">取消</el-button>
        <el-button type="primary" @click="saveKnowledge">确定</el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="showEditDialog" :title="editDialogTitle" width="500px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="代码" prop="code">
          <el-input v-model="editForm.code" placeholder="请输入代码" />
        </el-form-item>
        <el-form-item v-if="editForm.type === 'knowledge'" label="难度级别">
          <el-rate v-model="editForm.difficulty" show-text />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="editForm.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Folder, Document } from '@element-plus/icons-vue'

// 搜索关键词
const searchKeyword = ref('')

// 树形组件引用
const treeRef = ref()

// 树形结构配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 科目数据
const subjects = ref([
  { id: 1, name: '数据结构', code: 'DS' },
  { id: 2, name: '算法设计', code: 'AD' },
  { id: 3, name: '操作系统', code: 'OS' },
  { id: 4, name: '计算机网络', code: 'CN' },
  { id: 5, name: '数据库原理', code: 'DB' }
])

// 科目树形数据
const subjectTree = ref([
  {
    id: 1,
    name: '数据结构',
    code: 'DS',
    type: 'subject',
    description: '数据结构与算法基础',
    sort: 1,
    children: [
      {
        id: 101,
        name: '线性表',
        code: 'DS_LINEAR',
        type: 'knowledge',
        subjectId: 1,
        difficulty: 3,
        description: '线性表的基本概念和操作',
        questionCount: 25,
        sort: 1
      },
      {
        id: 102,
        name: '栈和队列',
        code: 'DS_STACK_QUEUE',
        type: 'knowledge',
        subjectId: 1,
        difficulty: 3,
        description: '栈和队列的实现与应用',
        questionCount: 18,
        sort: 2
      },
      {
        id: 103,
        name: '树和二叉树',
        code: 'DS_TREE',
        type: 'knowledge',
        subjectId: 1,
        difficulty: 4,
        description: '树的基本概念和二叉树操作',
        questionCount: 32,
        sort: 3
      }
    ]
  },
  {
    id: 2,
    name: '算法设计',
    code: 'AD',
    type: 'subject',
    description: '算法设计与分析',
    sort: 2,
    children: [
      {
        id: 201,
        name: '排序算法',
        code: 'AD_SORT',
        type: 'knowledge',
        subjectId: 2,
        difficulty: 3,
        description: '各种排序算法的实现和分析',
        questionCount: 28,
        sort: 1
      },
      {
        id: 202,
        name: '查找算法',
        code: 'AD_SEARCH',
        type: 'knowledge',
        subjectId: 2,
        difficulty: 3,
        description: '各种查找算法的实现和分析',
        questionCount: 22,
        sort: 2
      }
    ]
  }
])

// 对话框显示状态
const showAddSubjectDialog = ref(false)
const showKnowledgeDialog = ref(false)
const showEditDialog = ref(false)

// 科目表单
const subjectForm = reactive({
  name: '',
  code: '',
  description: '',
  sort: 0
})

// 知识点表单
const knowledgeForm = reactive({
  subjectId: null,
  name: '',
  code: '',
  difficulty: 3,
  description: '',
  sort: 0
})

// 编辑表单
const editForm = reactive({
  id: null,
  name: '',
  code: '',
  type: '',
  difficulty: 3,
  description: '',
  sort: 0
})

// 表单验证规则
const subjectRules = {
  name: [{ required: true, message: '请输入科目名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入科目代码', trigger: 'blur' }]
}

const knowledgeRules = {
  name: [{ required: true, message: '请输入知识点名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入知识点代码', trigger: 'blur' }]
}

const editRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入代码', trigger: 'blur' }]
}

// 表单引用
const subjectFormRef = ref()
const knowledgeFormRef = ref()
const editFormRef = ref()

// 计算属性
const knowledgeDialogTitle = computed(() => {
  return knowledgeForm.id ? '编辑知识点' : '新增知识点'
})

const editDialogTitle = computed(() => {
  return editForm.type === 'subject' ? '编辑科目' : '编辑知识点'
})

// 监听搜索关键词
watch(searchKeyword, (val) => {
  treeRef.value?.filter(val)
})

// 过滤节点
const filterNode = (value: string, data: any) => {
  if (!value) return true
  return data.name.includes(value) || data.code?.includes(value)
}

// 展开全部
const expandAll = () => {
  const nodes = treeRef.value?.store.nodesMap
  for (let node in nodes) {
    nodes[node].expanded = true
  }
}

// 收起全部
const collapseAll = () => {
  const nodes = treeRef.value?.store.nodesMap
  for (let node in nodes) {
    nodes[node].expanded = false
  }
}

// 添加知识点
const addKnowledgePoint = (subject: any) => {
  Object.assign(knowledgeForm, {
    subjectId: subject.id,
    name: '',
    code: '',
    difficulty: 3,
    description: '',
    sort: 0
  })
  showKnowledgeDialog.value = true
}

// 编辑节点
const editNode = (data: any) => {
  Object.assign(editForm, {
    id: data.id,
    name: data.name,
    code: data.code,
    type: data.type,
    difficulty: data.difficulty || 3,
    description: data.description || '',
    sort: data.sort || 0
  })
  showEditDialog.value = true
}

// 删除节点
const deleteNode = async (data: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除${data.type === 'subject' ? '科目' : '知识点'} "${data.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用删除API
    ElMessage.success('删除成功')
    
    // 从树形数据中移除
    if (data.type === 'subject') {
      const index = subjectTree.value.findIndex(item => item.id === data.id)
      if (index > -1) {
        subjectTree.value.splice(index, 1)
      }
    } else {
      // 删除知识点
      subjectTree.value.forEach(subject => {
        if (subject.children) {
          const index = subject.children.findIndex(item => item.id === data.id)
          if (index > -1) {
            subject.children.splice(index, 1)
          }
        }
      })
    }
  } catch {
    // 用户取消删除
  }
}

// 保存科目
const saveSubject = async () => {
  try {
    await subjectFormRef.value?.validate()
    
    // 这里应该调用保存API
    const newSubject = {
      id: Date.now(),
      ...subjectForm,
      type: 'subject',
      children: []
    }
    
    subjectTree.value.push(newSubject)
    subjects.value.push({ id: newSubject.id, name: newSubject.name, code: newSubject.code })
    
    ElMessage.success('科目创建成功')
    showAddSubjectDialog.value = false
    
    // 重置表单
    Object.assign(subjectForm, {
      name: '',
      code: '',
      description: '',
      sort: 0
    })
  } catch {
    // 验证失败
  }
}

// 保存知识点
const saveKnowledge = async () => {
  try {
    await knowledgeFormRef.value?.validate()
    
    // 这里应该调用保存API
    const newKnowledge = {
      id: Date.now(),
      ...knowledgeForm,
      type: 'knowledge',
      questionCount: 0
    }
    
    // 添加到对应科目下
    const subject = subjectTree.value.find(item => item.id === knowledgeForm.subjectId)
    if (subject) {
      if (!subject.children) {
        subject.children = []
      }
      subject.children.push(newKnowledge)
    }
    
    ElMessage.success('知识点创建成功')
    showKnowledgeDialog.value = false
    
    // 重置表单
    Object.assign(knowledgeForm, {
      subjectId: null,
      name: '',
      code: '',
      difficulty: 3,
      description: '',
      sort: 0
    })
  } catch {
    // 验证失败
  }
}

// 保存编辑
const saveEdit = async () => {
  try {
    await editFormRef.value?.validate()
    
    // 这里应该调用更新API
    ElMessage.success('更新成功')
    showEditDialog.value = false
    
    // 更新树形数据
    if (editForm.type === 'subject') {
      const subject = subjectTree.value.find(item => item.id === editForm.id)
      if (subject) {
        Object.assign(subject, editForm)
      }
    } else {
      // 更新知识点
      subjectTree.value.forEach(subject => {
        if (subject.children) {
          const knowledge = subject.children.find(item => item.id === editForm.id)
          if (knowledge) {
            Object.assign(knowledge, editForm)
          }
        }
      })
    }
  } catch {
    // 验证失败
  }
}
</script>

<style scoped>
.subject-management-view {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.page-subtitle {
  color: #666;
  margin: 0;
}

.dopamine-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toolbar-left {
  display: flex;
  gap: 12px;
}

.subject-tree-component {
  margin-top: 16px;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  padding: 4px 0;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.node-icon {
  font-size: 16px;
}

.subject-icon {
  color: #409eff;
}

.knowledge-icon {
  color: #67c23a;
}

.node-label {
  font-weight: 500;
  color: #1a1a1a;
}

.node-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.tree-node:hover .node-actions {
  opacity: 1;
}
</style>