<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="类型">
          <el-select v-model="searchInfo.action" placeholder="全部" clearable style="width: 160px">
            <el-option label="退费" value="refund" />
            <el-option label="转课" value="transfer" />
          </el-select>
        </el-form-item>
        <el-form-item label="报名ID">
          <el-input v-model.number="searchInfo.enrollmentId" placeholder="按报名ID筛选" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="退费时间">
          <el-date-picker v-model="searchInfo.startRefundTime" type="datetime" placeholder="开始时间" />
          —
          <el-date-picker v-model="searchInfo.endRefundTime" type="datetime" placeholder="结束时间" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID" style="width: 100%">
        <el-table-column align="left" label="退费时间" width="120">
          <template #default="scope">
            <div class="dt-cell">
              <div class="dt-date">{{ getDatePart(scope.row.refundTime || scope.row.CreatedAt) }}</div>
              <div class="dt-time">{{ getTimePart(scope.row.refundTime || scope.row.CreatedAt) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="学员姓名" width="120">
          <template #default="scope">{{ getUserName(scope.row.userId) }}</template>
        </el-table-column>
        <el-table-column align="left" label="课程名称" width="150">
          <template #default="scope">{{ getCourseName(scope.row.courseId) }}</template>
        </el-table-column>
        <el-table-column align="left" label="类型" width="100">
          <template #default="scope">{{ scope.row.action === 'transfer' ? '转课' : '退费' }}</template>
        </el-table-column>
        <el-table-column align="left" label="课时" prop="refundSessions" width="90" />
        <el-table-column align="left" label="金额" width="110">
          <template #default="scope">{{ formatMoney(scope.row.refundAmount) }}</template>
        </el-table-column>
        <el-table-column align="left" label="原因" prop="reason" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="操作人" prop="operatorName" width="120" />
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
  </div>
</template>

<script>
export default {
  name: 'EduRefund'
}
</script>

<script setup>
import { ref, onMounted } from 'vue'
import { getEduRefundList } from '@/api/eduRefund'
import { getUserList } from '@/api/user'
import { getEduCourseList } from '@/api/eduCourse'
import { formatDate } from '@/utils/format'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const userMap = ref({})
const courseMap = ref({})

const formatMoney = (val) => {
  const num = Number(val || 0)
  return num.toFixed(2)
}

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

const getUserName = (userId) => {
  if (!userId) return '-'
  return userMap.value[userId] || '-'
}

const getCourseName = (courseId) => {
  if (!courseId) return '-'
  return courseMap.value[courseId] || '-'
}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getEduRefundList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list || []
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const loadUserMap = async () => {
  const res = await getUserList({ page: 1, pageSize: 9999 })
  if (res.code === 0) {
    const map = {}
    ;(res.data.list || []).forEach((u) => {
      map[u.ID] = u.nickName || u.userName || '-'
    })
    userMap.value = map
  }
}

const loadCourseMap = async () => {
  const res = await getEduCourseList({ page: 1, pageSize: 9999 })
  if (res.code === 0) {
    const map = {}
    ;(res.data.list || []).forEach((c) => {
      map[c.ID] = c.courseName || '-'
    })
    courseMap.value = map
  }
}

onMounted(() => {
  loadUserMap()
  loadCourseMap()
  getTableData()
})
</script>

<style>
.dt-cell {
  line-height: 1.2;
}
.dt-date,
.dt-time {
  white-space: nowrap;
}
.dt-time {
  margin-top: 4px;
  font-size: 12px;
  opacity: 0.75;
}
</style>
