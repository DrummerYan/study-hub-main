<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="学员姓名">
        <el-input v-model="searchInfo.userName" placeholder="搜索学员姓名" clearable style="width: 200px" />
      </el-form-item>
      <el-form-item label="授课老师">
        <el-select v-model="searchInfo.teacherId" placeholder="选择老师" clearable filterable style="width: 200px">
          <el-option
            v-for="teacher in teacherList"
            :key="teacher.ID"
            :label="`${teacher.nickName} (${teacher.phone || teacher.userName})`"
            :value="teacher.ID"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="课程">
        <el-select v-model="searchInfo.courseId" placeholder="选择课程" clearable filterable style="width: 200px">
          <el-option
            v-for="course in courseList"
            :key="course.ID"
            :label="course.courseName"
            :value="course.ID"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="报名ID">
        <el-input v-model.number="searchInfo.enrollmentId" placeholder="按报名ID筛选" clearable style="width: 200px" />
      </el-form-item>
      <el-form-item label="操作时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button icon="back" @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div v-if="isSuperAdmin" class="gva-search-box" style="margin-bottom: 12px;">
          <el-form :inline="true">
            <el-form-item label="统计月份">
              <el-date-picker v-model="summaryMonth" type="month" value-format="YYYY-MM" placeholder="选择月份" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="refresh" @click="loadMonthlySummary">刷新</el-button>
            </el-form-item>
            <el-form-item>
              <el-tag type="success">计费节数：{{ monthlySummary.chargeableSessions }}</el-tag>
            </el-form-item>
            <el-form-item>
              <el-tag type="warning">非计费节数：{{ monthlySummary.nonChargeableSessions }}</el-tag>
            </el-form-item>
            <el-form-item>
              <el-tag type="info">计费金额：{{ formatMoney(monthlySummary.chargeableAmount) }}</el-tag>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-btn-list">
            <el-tag type="info" size="large" style="padding: 10px 20px;">
              <el-icon><InfoFilled /></el-icon>
              课时历史记录（仅查询，新增课时请在"报名管理"中操作）
            </el-tag>
            <el-button v-if="canViewLowSessions" type="warning" icon="bell" style="margin-left: 12px" @click="openLowSessions">
              课时不足提醒
            </el-button>
        </div>
        <el-table
        ref="multipleTable"
        class="class-session-table"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        >
        <el-table-column align="left" label="操作日期" width="120">
            <template #default="scope">
              <div class="dt-cell">
                <div class="dt-date">{{ getDatePart(scope.row.useDate) }}</div>
                <div class="dt-time">{{ getTimePart(scope.row.useDate) }}</div>
              </div>
            </template>
        </el-table-column>
        <el-table-column align="left" label="学员姓名" prop="userName" width="120" />
        <el-table-column align="left" label="授课老师" prop="teacherName" width="120" />
        <el-table-column align="left" label="课程名称" prop="courseName" width="150" />
        <el-table-column align="left" label="课时情况" width="130">
            <template #default="scope">
              <div class="session-cell">
                <el-tooltip
                  :content="`付${getEnrollmentPaidRemaining(scope.row)}/赠${getEnrollmentGiftRemaining(scope.row)}`"
                  placement="top"
                  :open-delay="0"
                >
                  <div class="session-line">
                    <span class="session-key">剩余</span>
                    <span class="session-value">{{ getEnrollmentRemaining(scope.row) }}</span>
                  </div>
                </el-tooltip>
              </div>
            </template>
        </el-table-column>
        <el-table-column align="left" label="操作类型" width="100">
            <template #default="scope">
                <el-tag :type="scope.row.action === 'add' ? 'success' : 'warning'">
                    {{ scope.row.action === 'add' ? '加课' : '消课' }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column
          align="center"
          header-align="center"
          label="课时数量"
          width="120"
          class-name="qty-col"
          header-class-name="qty-col"
        >
            <template #default="scope">
                <div class="qty-cell">
                  <span :style="{ color: scope.row.action === 'add' ? '#67C23A' : '#E6A23C', fontWeight: 'bold' }">
                    {{ scope.row.action === 'add' ? '+' : '-' }}{{ scope.row.numSessions }}
                  </span>
                </div>
            </template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="center" label="单价" width="120">
            <template #default="scope">{{ formatMoney(scope.row.unitPrice) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="center" label="金额" width="140">
            <template #default="scope">{{ formatMoney(scope.row.amount) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="center" label="计费" width="120">
            <template #default="scope">
                <el-tag :type="scope.row.chargeable ? 'success' : 'info'">
                    {{ scope.row.chargeable ? '计费' : '不计费' }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column align="left" label="操作原因" prop="reason" min-width="240" show-overflow-tooltip />
        <el-table-column align="center" label="记录时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="150" fixed="right">
            <template #default="scope">
            <el-button type="danger" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog v-model="lowSessionsVisible" title="剩余课时不足学员" width="700px">
      <el-table :data="lowSessions" style="width: 100%">
        <el-table-column align="left" label="学员姓名" prop="nickName" width="140" />
        <el-table-column align="left" label="手机号" prop="phone" width="140" />
        <el-table-column align="left" label="剩余课时" prop="remainingSessions" width="120" />
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="lowSessionsVisible = false">关 闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'EduClassSession'
}
</script>

<script setup>
import {
  deleteEduClassSession,
  getEduClassSessionList,
  getStudentsWithLessThanFiveSessions,
  getMonthlyChargeSummary
} from '@/api/eduClassSession'
import { getEduCourseList } from '@/api/eduCourse'
import { getUserList } from '@/api/user'

// 全量引入格式化工具 请按需保留
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

// ===== 日期/时间两行显示（YYYY-MM-DD / HH:mm:ss） =====
const normalizeDateTime = (val) => {
  if (!val) return ''
  const date = new Date(val)
  if (!Number.isNaN(date.getTime())) {
    return formatDate(date)
  }
  let s = String(val).replace('T', ' ').split('.')[0].trim()
  const parts = s.split(' ')
  if (parts.length === 1) return parts[0]
  let time = parts[1]
  const tparts = time.split(':')
  if (tparts.length === 2) {
    time = `${tparts[0].padStart(2, '0')}:${tparts[1].padStart(2, '0')}:00`
  } else if (tparts.length === 3) {
    time = `${tparts[0].padStart(2, '0')}:${tparts[1].padStart(2, '0')}:${tparts[2].padStart(2, '0')}`
  }
  return `${parts[0]} ${time}`.trim()
}

const getDatePart = (val) => {
  const s = normalizeDateTime(val)
  return s.split(' ')[0] || '-'
}

const getTimePart = (val) => {
  const s = normalizeDateTime(val)
  return s.split(' ')[1] || ''
}

const getEnrollmentTotal = (row) => {
  const total = row?.eduEnrollment?.totalSessions
  return total === null || total === undefined ? '-' : total
}

const getEnrollmentRemaining = (row) => {
  const remaining = row?.eduEnrollment?.remainingSessions
  return remaining === null || remaining === undefined ? '-' : remaining
}
const getEnrollmentPaidTotal = (row) => {
  const val = row?.eduEnrollment?.paidSessions
  return val === null || val === undefined ? '-' : val
}
const getEnrollmentGiftTotal = (row) => {
  const val = row?.eduEnrollment?.giftSessions
  return val === null || val === undefined ? '-' : val
}
const getEnrollmentPaidRemaining = (row) => {
  const val = row?.eduEnrollment?.remainingPaidSessions
  return val === null || val === undefined ? '-' : val
}
const getEnrollmentGiftRemaining = (row) => {
  const val = row?.eduEnrollment?.remainingGiftSessions
  return val === null || val === undefined ? '-' : val
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isSuperAdmin = computed(() => {
  const roleId = userStore.userInfo?.authority?.authorityId || userStore.userInfo?.authorityId || userStore.userInfo?.authority_id
  return roleId === 888
})
const formatMoney = (val) => {
  const num = Number(val || 0)
  return num.toFixed(2)
}
const canViewLowSessions = computed(() => {
  const roleId = userStore.userInfo?.authority?.authorityId || userStore.userInfo?.authorityId
  return roleId !== 9002 // 学员不显示
})

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const teacherList = ref([])
const courseList = ref([])
const lowSessionsVisible = ref(false)
const lowSessions = ref([])
const now = new Date()
const summaryMonth = ref(`${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`)
const monthlySummary = ref({
  chargeableSessions: 0,
  nonChargeableSessions: 0,
  chargeableAmount: 0
})

// 重置
const onReset = () => {
  searchInfo.value = {}
  // 如果是从报名列表跳转过来的，保留 enrollmentId
  if (route.query.enrollmentId) {
    searchInfo.value.enrollmentId = Number(route.query.enrollmentId)
  }
  getTableData()
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 初始化下拉选项（老师/课程）
const setOptions = async () => {
  const userRes = await getUserList({ page: 1, pageSize: 9999 })
  if (userRes.code === 0) {
    const allUsers = userRes.data.list || []
    // 只显示教师角色
    teacherList.value = allUsers.filter(user => {
      const userRoles = user.authorities || []
      if (user.authorityId === 9001 || user.authority_id === 9001) return true
      return userRoles.some(auth =>
        auth.authorityId === 9001 ||
        auth.authorityName?.includes('教师') ||
        auth.authorityName?.includes('老师')
      )
    })
  }
  const courseRes = await getEduCourseList({ page: 1, pageSize: 9999 })
  if (courseRes.code === 0) {
    courseList.value = courseRes.data.list || []
  }
}

// 查询
const getTableData = async() => {
  const table = await getEduClassSessionList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const loadMonthlySummary = async () => {
  if (!isSuperAdmin.value) return
  const res = await getMonthlyChargeSummary({ month: summaryMonth.value })
  if (res.code === 0) {
    monthlySummary.value = res.data.summary || {
      chargeableSessions: 0,
      nonChargeableSessions: 0,
      chargeableAmount: 0
    }
  }
}

watch(summaryMonth, () => loadMonthlySummary())

// 初始化：检查 URL 参数
onMounted(() => {
  if (route.query.enrollmentId) {
    searchInfo.value.enrollmentId = Number(route.query.enrollmentId)
  }
  setOptions()
  getTableData()
  loadMonthlySummary()
})

// ============== 表格控制部分结束 ===============

// 删除行（直接删除，权限由后端Casbin控制）
const deleteRow = (row) => {
    ElMessageBox.confirm(
      '确定要删除这条课时记录吗？',
      '确认删除',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(() => {
      deleteEduClassSessionFunc(row)
    }).catch(() => {
      // 用户取消删除
    })
}

// 删除记录
const deleteEduClassSessionFunc = async (row) => {
    const res = await deleteEduClassSession({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
            type: 'success',
            message: '删除成功'
        })
        if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    } else {
        // 显示后端返回的错误信息（包括权限不足）
        ElMessage({
            type: 'error',
            message: res.msg || '删除失败'
        })
    }
}

const openLowSessions = async () => {
  const res = await getStudentsWithLessThanFiveSessions()
  if (res.code === 0) {
    lowSessions.value = res.data.list || []
    lowSessionsVisible.value = true
  } else {
    ElMessage({
      type: 'error',
      message: res.msg || '获取失败'
    })
  }
}
</script>

<style>
/* 日期/时间两行显示（YYYY-MM-DD / HH:mm:ss）单元格 */
.dt-cell {
  line-height: 1.25;
}

.dt-date,
.dt-time {
  white-space: nowrap; /* 防止被挤成“2026-02-05 1 / 1:44:00”这种断行 */
}

.dt-time {
  margin-top: 6px;
  font-size: 12px;
  opacity: 0.75;
}

.class-session-table .el-table__cell {
  vertical-align: middle;
  padding: 10px 0; /* 轻微增加行高，避免拥挤 */
}

/* 表头稍微松一点 */
.class-session-table .el-table__header-wrapper th.el-table__cell {
  padding: 12px 0;
}

/* 单元格内容左右留一点空隙（不影响对齐逻辑） */
.class-session-table .el-table__cell .cell {
  padding-left: 14px;
  padding-right: 14px;
}
.class-session-table .dt-cell {
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;
}
.class-session-table .session-cell {
  display: flex;
  align-items: center;
  height: 100%;
}


.qty-cell {
  width: 100%;
  text-align: center;
}

/* 强制“课时数量”这一列（表头+单元格）居中，避免被全局 .el-table .cell 样式覆盖 */
.class-session-table .qty-col .cell {
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center !important;
}

.class-session-table .qty-col .cell .qty-cell {
  width: 100%;
  text-align: center;
}

.session-line {
  display: grid;
  grid-template-columns: 32px auto;
  align-items: center;
  column-gap: 6px;
  white-space: nowrap;
}

/* 让 tag/number 列不显得太拥挤 */
.class-session-table .el-tag {
  line-height: 20px;
}
.session-key {
  color: #606266;
  text-align: right;
}
.session-value {
  color: #303133;
  font-weight: 600;
}
</style>
