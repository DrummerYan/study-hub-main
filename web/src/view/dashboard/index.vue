<template>
  <div class="page">
    <div class="gva-card-box">
      <div class="gva-card gva-top-card">
        <div class="gva-top-card-left">
          <div class="gva-top-card-left-title">
            {{ greetingMessage }}
          </div>
          <div class="gva-top-card-left-dot">{{ weatherInfo }}</div>
          <div class="gva-top-card-left-rows">
            <el-row :gutter="20">
              <el-col
                v-for="(card, key) in toolCards"
                :key="key"
                :span="4"
                :xs="8"
                class="quick-entrance-items"
                @click="toTarget(card.name)"
              >
                <div class="quick-entrance-item">
                  <div
                    class="quick-entrance-item-icon"
                    :style="{ backgroundColor: card.bg }"
                  >
                    <el-icon>
                      <component
                        :is="card.icon"
                        :style="{ color: card.color }"
                      />
                    </el-icon>
                  </div>
                  <p>{{ card.label }}</p>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
        <img src="@/assets/dashboard.png" class="gva-top-card-right" alt>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWeatherInfo } from '@/view/dashboard/weather.js'
import { useUserStore } from '@/pinia/modules/user'
import { useRouterStore } from '@/pinia/modules/router'

const userStore = useUserStore()
const routerStore = useRouterStore()
const weatherInfo = useWeatherInfo()

// 根据时间和用户信息生成问候语
const greetingMessage = computed(() => {
  const hour = new Date().getHours()
  let timeGreeting = ''
  
  // 根据时间段显示不同的问候
  if (hour >= 5 && hour < 12) {
    timeGreeting = '早安'
  } else if (hour >= 12 && hour < 18) {
    timeGreeting = '午安'
  } else if (hour >= 18 && hour < 22) {
    timeGreeting = '晚上好'
  } else {
    timeGreeting = '夜深了'
  }
  
  // 获取用户昵称或角色名
  const roleName = userStore.userInfo?.authority?.authorityName || ''
  const nickName = userStore.userInfo?.nickName || ''
  // 当昵称是“超级管理员/管理员”等通用占位时，优先显示当前角色名
  const isGenericAdminName = nickName === '超级管理员' || nickName === '管理员'
  const displayName = (roleName && isGenericAdminName)
    ? roleName
    : (nickName || roleName || '用户')
  
  // 根据角色显示不同的祝福语
  let message = ''
  if (roleName.includes('学员') || roleName.includes('学生')) {
    message = `${timeGreeting}，${displayName}，欢迎回来！今天也要加油学习哦~`
  } else if (roleName.includes('教师') || roleName.includes('老师')) {
    const teacherName = (displayName.includes('教师') || displayName.includes('老师'))
      ? displayName
      : `${displayName}老师`
    message = `${timeGreeting}，${teacherName}，欢迎回来！`
  } else if (roleName.includes('管理员')) {
    message = `${timeGreeting}，${displayName}，请开始一天的工作吧`
  } else {
    message = `${timeGreeting}，${displayName}，欢迎回来！`
  }
  
  return message
})

const quickEntryDefs = [
  {
    label: '用户管理',
    icon: 'monitor',
    name: 'user',
    color: '#ff9c6e',
    bg: 'rgba(255, 156, 110,.3)',
  },
  {
    label: '角色管理',
    icon: 'setting',
    name: 'authority',
    color: '#69c0ff',
    bg: 'rgba(105, 192, 255,.3)',
  },
  {
    label: '菜单管理',
    icon: 'menu',
    name: 'menu',
    color: '#b37feb',
    bg: 'rgba(179, 127, 235,.3)',
  },
  {
    label: '组织机构管理',
    icon: 'cpu',
    name: 'eduOrganization',
    color: '#ffd666',
    bg: 'rgba(255, 214, 102,.3)',
  },
  {
    label: '课程管理',
    icon: 'document-checked',
    name: 'eduCourse',
    color: '#ff85c0',
    bg: 'rgba(255, 133, 192,.3)',
  },
]

// 快捷入口：名称取当前菜单标题，权限不足则不显示
const toolCards = computed(() => {
  // 依赖异步路由加载标记，确保菜单加载后重新计算
  routerStore.asyncRouterFlag
  return quickEntryDefs
    .map((item) => {
      const routeInfo = routerStore.routeMap[item.name]
      if (!routeInfo || routeInfo.hidden || routeInfo.meta?.hidden) return null
      return {
        ...item,
        label: routeInfo.meta?.title || item.label,
      }
    })
    .filter(Boolean)
})

const router = useRouter()

const toTarget = (name) => {
  router.push({ name })
}
</script>
<script>
export default {
  name: 'Dashboard',
}
</script>

<style lang="scss" scoped>
@mixin flex-center {
  display: flex;
  align-items: center;
}
.page {
  background: #f0f2f5;
  padding: 0;
  .gva-card-box {
    padding: 12px 16px;
    & + .gva-card-box {
      padding-top: 0px;
    }
  }
  .gva-card {
    box-sizing: border-box;
    background-color: #fff;
    border-radius: 2px;
    height: auto;
    padding: 26px 30px;
    overflow: hidden;
    box-shadow: 0 0 7px 1px rgba(0, 0, 0, 0.03);
  }
  .gva-top-card {
    height: 260px;
    @include flex-center;
    justify-content: space-between;
    color: #777;
    &-left {
      height: 100%;
      display: flex;
      flex-direction: column;
      &-title {
        font-size: 22px;
        color: #343844;
      }
      &-dot {
        font-size: 16px;
        color: #6b7687;
        margin-top: 24px;
      }
      &-rows {
        // margin-top: 15px;
        margin-top: 18px;
        color: #6b7687;
        width: 600px;
        align-items: center;
      }
      &-item {
        + .gva-top-card-left-item {
          margin-top: 24px;
        }
        margin-top: 14px;
      }
    }
    &-right {
      height: 600px;
      width: 600px;
      margin-top: 28px;
    }
  }
  ::v-deep(.el-card__header) {
    padding: 0;
    border-bottom: none;
  }
  .card-header {
    padding-bottom: 20px;
    border-bottom: 1px solid #e8e8e8;
  }
  .quick-entrance-title {
    height: 30px;
    font-size: 22px;
    color: #333;
    width: 100%;
    border-bottom: 1px solid #eee;
  }
  .quick-entrance-items {
    @include flex-center;
    justify-content: center;
    text-align: center;
    color: #333;
    .quick-entrance-item {
      padding: 16px 28px;
      margin-top: -16px;
      margin-bottom: -16px;
      border-radius: 4px;
      transition: all 0.2s;
      &:hover {
        box-shadow: 0px 0px 7px 0px rgba(217, 217, 217, 0.55);
      }
      cursor: pointer;
      height: auto;
      text-align: center;
      // align-items: center;
      &-icon {
        width: 50px;
        height: 50px !important;
        border-radius: 8px;
        @include flex-center;
        justify-content: center;
        margin: 0 auto;
        i {
          font-size: 24px;
        }
      }
      p {
        margin-top: 10px;
      }
    }
  }
  .echart-box {
    padding: 14px;
  }
}
.dashboard-icon {
  font-size: 20px;
  color: rgb(85, 160, 248);
  width: 30px;
  height: 30px;
  margin-right: 10px;
  @include flex-center;
}
.flex-center {
  @include flex-center;
}

//小屏幕不显示右侧，将登录框居中
@media (max-width: 750px) {
  .gva-card {
    padding: 20px 10px !important;
    .gva-top-card {
      height: auto;
      &-left {
        &-title {
          font-size: 20px !important;
        }
        &-rows {
          margin-top: 15px;
          align-items: center;
        }
      }
      &-right {
        display: none;
      }
    }
    .gva-middle-card {
      &-item {
        line-height: 20px;
      }
    }
    .dashboard-icon {
      font-size: 18px;
    }
  }
}
</style>
