<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="学员姓名">
        <el-input v-model="searchInfo.userName" placeholder="搜索学员姓名" clearable style="width: 200px" />
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
        <div class="gva-btn-list">
            <el-tag type="info" size="large" style="padding: 10px 20px;">
              <el-icon><InfoFilled /></el-icon>
              课时历史记录（仅查询，新增课时请在"报名管理"中操作）
            </el-tag>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        >
        <el-table-column align="left" label="操作日期" width="120">
            <template #default="scope">{{ formatDate(scope.row.useDate) }}</template>
        </el-table-column>
        <el-table-column align="left" label="学员姓名" prop="userName" width="120" />
        <el-table-column align="left" label="授课老师" prop="teacherName" width="120" />
        <el-table-column align="left" label="课程名称" prop="courseName" width="150" />
        <el-table-column align="left" label="操作类型" width="100">
            <template #default="scope">
                <el-tag :type="scope.row.action === 'add' ? 'success' : 'warning'">
                    {{ scope.row.action === 'add' ? '加课' : '消课' }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column align="left" label="课时数量" width="100">
            <template #default="scope">
                <span :style="{ color: scope.row.action === 'add' ? '#67C23A' : '#E6A23C', fontWeight: 'bold' }">
                    {{ scope.row.action === 'add' ? '+' : '-' }}{{ scope.row.numSessions }}
                </span>
            </template>
        </el-table-column>
        <el-table-column align="left" label="操作原因" prop="reason" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="记录时间" width="180">
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
  getEduClassSessionList
} from '@/api/eduClassSession'

// 全量引入格式化工具 请按需保留
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

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

// 初始化：检查 URL 参数
onMounted(() => {
  if (route.query.enrollmentId) {
    searchInfo.value.enrollmentId = Number(route.query.enrollmentId)
  }
  getTableData()
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
</script>

<style>
</style>
