<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="学员姓名">
        <el-input v-model="searchInfo.userName" placeholder="搜索学员姓名" clearable style="width: 200px" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="searchInfo.userPhone" placeholder="搜索手机号" clearable style="width: 200px" />
      </el-form-item>
      <el-form-item label="课程">
        <el-select v-model="searchInfo.courseId" placeholder="选择课程" clearable style="width: 200px">
          <el-option
            v-for="course in courseList"
            :key="course.ID"
            :label="course.courseName"
            :value="course.ID"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增报名</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="confirmBatchDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="学员姓名" prop="userName" width="120" />
        <el-table-column align="left" label="手机号" prop="userPhone" width="130" />
        <el-table-column align="left" label="课程名称" prop="eduCourse.courseName" width="180" />
        <el-table-column align="left" label="总课时" prop="totalSessions" width="90" />
        <el-table-column align="left" label="剩余课时" prop="remainingSessions" width="100">
            <template #default="scope">
                <el-tag :type="scope.row.remainingSessions <= 5 ? 'danger' : 'success'">
                    {{ scope.row.remainingSessions }}
                </el-tag>
            </template>
        </el-table-column>
        <el-table-column align="left" label="报名日期" width="120">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="300" fixed="right">
            <template #default="scope">
            <el-button type="warning" link icon="remove" @click="openConsumeDialog(scope.row)">消课</el-button>
            <el-button type="success" link icon="plus" @click="openAddDialog(scope.row)">加课</el-button>
            <el-button type="info" link icon="tickets" @click="viewHistory(scope.row)">历史</el-button>
            <el-button type="primary" link icon="edit" @click="updateEduEnrollmentFunc(scope.row)">编辑</el-button>
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="新增报名">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="100px">
        <el-form-item label="选择学员:"  prop="userId" >
          <el-select v-model="formData.userId" placeholder="请选择学员" filterable clearable style="width: 100%">
            <el-option
              v-for="user in userList"
              :key="user.ID"
              :label="`${user.nickName} (${user.phone})`"
              :value="user.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择课程:"  prop="courseId" >
          <el-select v-model="formData.courseId" placeholder="请选择课程" filterable clearable style="width: 100%">
            <el-option
              v-for="course in courseList"
              :key="course.ID"
              :label="course.courseName"
              :value="course.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="总课时数:"  prop="totalSessions" >
          <el-input-number v-model="formData.totalSessions" :min="1" :max="1000" placeholder="请输入总课时数" style="width: 100%" @change="onTotalSessionsChange" />
        </el-form-item>
        <el-form-item label="剩余课时:" prop="remainingSessions">
          <el-input-number v-model="formData.remainingSessions" :min="0" :max="formData.totalSessions || 1000" placeholder="剩余课时（默认等于总课时）" style="width: 100%" />
          <div style="color: #999; font-size: 12px; margin-top: 4px;">提示：新建时默认等于总课时</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 消课弹窗 -->
    <el-dialog v-model="consumeDialogVisible" :before-close="closeConsumeDialog" title="消课" width="500px">
      <el-form :model="consumeForm" label-position="right" ref="consumeFormRef" :rules="consumeRules" label-width="100px">
        <el-form-item label="学员:">
          <el-input v-model="consumeForm.userName" :disabled="true" />
        </el-form-item>
        <el-form-item label="课程:">
          <el-input v-model="consumeForm.courseName" :disabled="true" />
        </el-form-item>
        <el-form-item label="剩余课时:">
          <el-tag type="success" size="large">{{ consumeForm.remainingSessions }} 课时</el-tag>
        </el-form-item>
        <el-form-item label="消课数量:" prop="sessionsToConsume">
          <el-input-number v-model="consumeForm.sessionsToConsume" :min="1" :max="consumeForm.remainingSessions" :step="1" />
        </el-form-item>
        <el-form-item label="上课日期:" prop="useDate">
          <el-date-picker v-model="consumeForm.useDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注原因:" prop="reasonType">
          <el-select v-model="consumeForm.reasonType" placeholder="请选择原因" style="width: 100%" @change="onConsumeReasonChange">
            <el-option label="正常消课" value="正常消课" />
            <el-option label="请假扣课" value="请假扣课" />
            <el-option label="补课消课" value="补课消课" />
            <el-option label="试听消课" value="试听消课" />
            <el-option label="自定义" value="自定义" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="consumeForm.reasonType === '自定义'" label="详细说明:" prop="customReason">
          <el-input v-model="consumeForm.customReason" type="textarea" :rows="3" placeholder="请输入详细原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeConsumeDialog">取 消</el-button>
          <el-button type="warning" @click="confirmConsume">确认消课</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 加课弹窗 -->
    <el-dialog v-model="addDialogVisible" :before-close="closeAddDialog" title="加课" width="500px">
      <el-form :model="addForm" label-position="right" ref="addFormRef" :rules="addRules" label-width="100px">
        <el-form-item label="学员:">
          <el-input v-model="addForm.userName" :disabled="true" />
        </el-form-item>
        <el-form-item label="课程:">
          <el-input v-model="addForm.courseName" :disabled="true" />
        </el-form-item>
        <el-form-item label="当前剩余:">
          <el-tag type="info" size="large">{{ addForm.remainingSessions }} 课时</el-tag>
        </el-form-item>
        <el-form-item label="增加数量:" prop="sessionsToAdd">
          <el-input-number v-model="addForm.sessionsToAdd" :min="1" :max="999" :step="1" />
        </el-form-item>
        <el-form-item label="日期:" prop="useDate">
          <el-date-picker v-model="addForm.useDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注原因:" prop="reasonType">
          <el-select v-model="addForm.reasonType" placeholder="请选择原因" style="width: 100%" @change="onAddReasonChange">
            <el-option label="学员购课" value="学员购课" />
            <el-option label="赠送课时" value="赠送课时" />
            <el-option label="补录课时" value="补录课时" />
            <el-option label="活动赠课" value="活动赠课" />
            <el-option label="退课返还" value="退课返还" />
            <el-option label="自定义" value="自定义" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="addForm.reasonType === '自定义'" label="详细说明:" prop="customReason">
          <el-input v-model="addForm.customReason" type="textarea" :rows="3" placeholder="请输入详细原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeAddDialog">取 消</el-button>
          <el-button type="success" @click="confirmAdd">确认加课</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'EduEnrollment'
}
</script>

<script setup>
import {
  createEduEnrollment,
  deleteEduEnrollment,
  deleteEduEnrollmentByIds,
  updateEduEnrollment,
  findEduEnrollment,
  getEduEnrollmentList,
  consumeSession,
  addSession
} from '@/api/eduEnrollment'
import { getEduCourseList } from '@/api/eduCourse'
import { getUserList } from '@/api/user'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 辅助函数：格式化日期为 YYYY-MM-DD 格式
const formatDateOnly = (date) => {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        userId: 0,
        courseId: 0,
        totalSessions: 0,
        remainingSessions: 0,
        })

// 验证规则
const rule = reactive({
  userId: [{ required: true, message: '请选择学员', trigger: 'change' }],
  courseId: [{ required: true, message: '请选择课程', trigger: 'change' }],
  totalSessions: [
    { required: true, message: '请输入总课时数', trigger: 'blur' },
    { type: 'number', min: 1, message: '总课时数必须大于0', trigger: 'blur' }
  ],
  remainingSessions: [
    { required: true, message: '请输入剩余课时数', trigger: 'blur' },
    { type: 'number', min: 0, message: '剩余课时数不能为负数', trigger: 'blur' }
  ]
})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
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
  const table = await getEduEnrollmentList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 用户列表和课程列表
const userList = ref([])
const courseList = ref([])

// 获取用户列表和课程列表
const setOptions = async () =>{
  // 获取所有用户（不分页）
  const userRes = await getUserList({ page: 1, pageSize: 9999 })
  if (userRes.code === 0) {
    // 只显示学员角色的用户（角色ID为9002，或角色名包含"学员"/"学生"）
    const allUsers = userRes.data.list || []
    userList.value = allUsers.filter(user => {
      const userRoles = user.authorities || []
      return userRoles.some(auth => 
        auth.authorityId === 9002 ||  // 精确匹配学员角色ID
        auth.authorityName?.includes('学员') || 
        auth.authorityName?.includes('学生')
      )
    })
  }
  
  // 获取所有课程（不分页）
  const courseRes = await getEduCourseList({ page: 1, pageSize: 9999 })
  if (courseRes.code === 0) {
    courseList.value = courseRes.data.list || []
  }
}

// 初始化加载
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteEduEnrollmentFunc(row)
        })
    }


// 批量删除确认
const confirmBatchDelete = () => {
  if (multipleSelection.value.length === 0) {
    ElMessage({
      type: 'warning',
      message: '请选择要删除的数据'
    })
    return
  }
  
  ElMessageBox.confirm(
    `确定要删除选中的 ${multipleSelection.value.length} 条记录吗？`,
    '批量删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    onDelete()
  }).catch(() => {
    // 用户取消删除
  })
}

// 批量删除
const onDelete = async() => {
      const ids = []
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          ids.push(item.ID)
        })
      const res = await deleteEduEnrollmentByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateEduEnrollmentFunc = async(row) => {
    const res = await findEduEnrollment({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
        formData.value = res.data.reeduEnrollment
        dialogFormVisible.value = true
    }
}


// 删除行
const deleteEduEnrollmentFunc = async (row) => {
    const res = await deleteEduEnrollment({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    // 重置表单数据，确保是数字类型
    formData.value = {
        userId: 0,
        courseId: 0,
        totalSessions: 10,  // 默认10课时
        remainingSessions: 10,  // 默认等于总课时
    }
    dialogFormVisible.value = true
}

// 总课时数改变时，自动同步剩余课时（仅新建时）
const onTotalSessionsChange = (val) => {
  if (type.value === 'create') {
    formData.value.remainingSessions = val || 0
  }
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        userId: 0,
        courseId: 0,
        totalSessions: 0,
        remainingSessions: 0,
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createEduEnrollment(formData.value)
                  break
                case 'update':
                  res = await updateEduEnrollment(formData.value)
                  break
                default:
                  res = await createEduEnrollment(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

// ============== 消课功能 ==============
const consumeDialogVisible = ref(false)
const consumeFormRef = ref()
const consumeForm = ref({
  userId: 0,
  courseId: 0,
  userName: '',
  courseName: '',
  remainingSessions: 0,
  sessionsToConsume: 1,
  useDate: formatDateOnly(new Date()),
  reasonType: '',  // 原因类型（下拉选择）
  customReason: ''  // 自定义原因（仅当reasonType为"自定义"时使用）
})

const consumeRules = reactive({
  sessionsToConsume: [
    { required: true, message: '请输入消课数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '消课数量至少为1', trigger: 'blur' }
  ],
  useDate: [
    { required: true, message: '请选择上课日期', trigger: 'change' }
  ],
  reasonType: [
    { required: true, message: '请选择原因', trigger: 'change' }
  ],
  customReason: [
    { required: true, message: '请输入详细原因', trigger: 'blur' }
  ]
})

// 消课原因类型改变时的处理
const onConsumeReasonChange = (value) => {
  // 如果不是自定义，清空自定义原因
  if (value !== '自定义') {
    consumeForm.value.customReason = ''
  }
}

// 打开消课弹窗
const openConsumeDialog = (row) => {
  consumeForm.value = {
    userId: row.userId,
    courseId: row.courseId,
    userName: row.userName,
    courseName: row.eduCourse?.courseName || '',
    remainingSessions: row.remainingSessions || 0,
    sessionsToConsume: 1,
    useDate: formatDateOnly(new Date()),
    reasonType: '',
    customReason: ''
  }
  consumeDialogVisible.value = true
}

// 关闭消课弹窗
const closeConsumeDialog = () => {
  consumeDialogVisible.value = false
  consumeFormRef.value?.resetFields()
}

// 确认消课
const confirmConsume = async () => {
  consumeFormRef.value?.validate(async (valid) => {
    if (!valid) return
    
    // 组合最终的原因：如果是自定义则使用customReason，否则使用reasonType
    const finalReason = consumeForm.value.reasonType === '自定义' 
      ? consumeForm.value.customReason 
      : consumeForm.value.reasonType
    
    const params = {
      userId: consumeForm.value.userId,
      courseId: consumeForm.value.courseId,
      sessionsToConsume: consumeForm.value.sessionsToConsume,
      useDate: consumeForm.value.useDate,
      reason: finalReason
    }
    
    const res = await consumeSession(params)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '消课成功'
      })
      closeConsumeDialog()
      getTableData()
    }
  })
}

// ============== 加课功能 ==============
const addDialogVisible = ref(false)
const addFormRef = ref()
const addForm = ref({
  userId: 0,
  courseId: 0,
  userName: '',
  courseName: '',
  remainingSessions: 0,
  sessionsToAdd: 1,
  useDate: formatDateOnly(new Date()),
  reasonType: '',  // 原因类型（下拉选择）
  customReason: ''  // 自定义原因（仅当reasonType为"自定义"时使用）
})

const addRules = reactive({
  sessionsToAdd: [
    { required: true, message: '请输入加课数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '加课数量至少为1', trigger: 'blur' }
  ],
  useDate: [
    { required: true, message: '请选择日期', trigger: 'change' }
  ],
  reasonType: [
    { required: true, message: '请选择原因', trigger: 'change' }
  ],
  customReason: [
    { required: true, message: '请输入详细原因', trigger: 'blur' }
  ]
})

// 加课原因类型改变时的处理
const onAddReasonChange = (value) => {
  // 如果不是自定义，清空自定义原因
  if (value !== '自定义') {
    addForm.value.customReason = ''
  }
}

// 打开加课弹窗
const openAddDialog = (row) => {
  addForm.value = {
    userId: row.userId,
    courseId: row.courseId,
    userName: row.userName,
    courseName: row.eduCourse?.courseName || '',
    remainingSessions: row.remainingSessions || 0,
    sessionsToAdd: 1,
    useDate: formatDateOnly(new Date()),
    reasonType: '',
    customReason: ''
  }
  addDialogVisible.value = true
}

// 关闭加课弹窗
const closeAddDialog = () => {
  addDialogVisible.value = false
  addFormRef.value?.resetFields()
}

// 确认加课
const confirmAdd = async () => {
  addFormRef.value?.validate(async (valid) => {
    if (!valid) return
    
    // 组合最终的原因：如果是自定义则使用customReason，否则使用reasonType
    const finalReason = addForm.value.reasonType === '自定义' 
      ? addForm.value.customReason 
      : addForm.value.reasonType
    
    const params = {
      userId: addForm.value.userId,
      courseId: addForm.value.courseId,
      sessionsToAdd: addForm.value.sessionsToAdd,
      useDate: addForm.value.useDate,
      reason: finalReason
    }
    
    const res = await addSession(params)
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '加课成功'
      })
      closeAddDialog()
      getTableData()
    }
  })
}

// ============== 查看历史 ==============
const viewHistory = (row) => {
  // 跳转到课时记录页面，传递 enrollmentId 参数
  // 使用 name 而不是 path，因为这是动态路由
  router.push({
    name: 'eduClassSession',
    query: { enrollmentId: row.ID }
  })
}
</script>

<style>
</style>
