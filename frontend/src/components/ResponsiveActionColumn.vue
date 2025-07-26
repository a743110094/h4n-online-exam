<template>
  <div class="responsive-action-column" ref="actionColumnRef">
    <!-- 宽度充足时显示完整按钮 -->
    <div v-if="displayMode === 'full'" class="action-buttons-full">
      <slot name="actions" :mode="'full'"></slot>
    </div>
    
    <!-- 中等宽度时只显示图标 -->
    <div v-else-if="displayMode === 'icon'" class="action-buttons-icon">
      <slot name="actions" :mode="'icon'"></slot>
    </div>
    
    <!-- 窄宽度时显示三点菜单 -->
    <div v-else class="action-buttons-menu">
      <el-dropdown trigger="hover" placement="bottom-end">
        <el-button type="text" class="more-button">
          <el-icon><MoreFilled /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <slot name="menu-items"></slot>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { MoreFilled } from '@element-plus/icons-vue'

type DisplayMode = 'full' | 'icon' | 'menu'

const actionColumnRef = ref<HTMLElement>()
const displayMode = ref<DisplayMode>('full')

const updateDisplayMode = () => {
  if (!actionColumnRef.value) return
  
  const width = actionColumnRef.value.offsetWidth
  
  if (width >= 150) {
    displayMode.value = 'full'
  } else if (width >= 80) {
    displayMode.value = 'icon'
  } else {
    displayMode.value = 'menu'
  }
}

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  if (actionColumnRef.value) {
    updateDisplayMode()
    
    // 使用 ResizeObserver 监听宽度变化
    resizeObserver = new ResizeObserver(() => {
      updateDisplayMode()
    })
    resizeObserver.observe(actionColumnRef.value)
  }
})

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
})
</script>

<style scoped>
.responsive-action-column {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-buttons-full {
  display: flex;
  gap: 8px;
  flex-wrap: nowrap;
}

.action-buttons-icon {
  display: flex;
  gap: 4px;
  flex-wrap: nowrap;
}

.action-buttons-menu {
  display: flex;
  justify-content: center;
}

.more-button {
  padding: 4px 8px;
  min-width: auto;
}

.more-button:hover {
  background-color: var(--el-color-primary-light-9);
}
</style>