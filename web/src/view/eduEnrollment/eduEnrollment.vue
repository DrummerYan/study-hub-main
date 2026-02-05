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
        <el-table-column v-if="isSuperAdmin" align="left" label="课时单价" width="110">
            <template #default="scope">{{ formatMoney(scope.row.pricePerSession) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="left" label="优惠" width="90">
            <template #default="scope">{{ formatMoney(scope.row.discountAmount) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="left" label="应收总额" width="120">
            <template #default="scope">{{ formatMoney(scope.row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="left" label="已收" width="90">
            <template #default="scope">{{ formatMoney(scope.row.paidAmount) }}</template>
        </el-table-column>
        <el-table-column v-if="isSuperAdmin" align="left" label="应收余额" width="120">
            <template #default="scope">{{ formatMoney(scope.row.balanceAmount) }}</template>
        </el-table-column>
        <el-table-column align="left" label="报名日期" width="120">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="300" fixed="right">
            <template #default="scope">
            <el-button type="warning" link icon="remove" @click="openConsumeDialog(scope.row)">消课</el-button>
            <el-button type="success" link icon="plus" @click="openAddDialog(scope.row)">加课</el-button>
            <el-button type="info" link icon="tickets" @click="viewHistory(scope.row)">历史</el-button>
            <el-button v-if="isSuperAdmin" type="info" link icon="wallet" @click="openPaymentDialog(scope.row)">收款记录</el-button>
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
        <template v-if="isSuperAdmin">
          <el-form-item label="课时单价:" prop="pricePerSession">
            <el-input-number v-model="formData.pricePerSession" :min="0" :precision="2" :step="1" placeholder="单节课价格" style="width: 100%" @change="recalcEnrollmentAmounts" />
          </el-form-item>
          <el-form-item label="优惠金额:" prop="discountAmount">
            <el-input-number v-model="formData.discountAmount" :min="0" :precision="2" :step="1" placeholder="优惠金额" style="width: 100%" @change="recalcEnrollmentAmounts" />
          </el-form-item>
          <el-form-item label="应收总额:">
            <el-input-number v-model="formData.totalAmount" :precision="2" :controls="false" :disabled="true" style="width: 100%" />
          </el-form-item>
          <el-form-item label="已收金额:">
            <el-input-number v-model="formData.paidAmount" :precision="2" :controls="false" :disabled="true" style="width: 100%" />
          </el-form-item>
          <el-form-item label="应收余额:">
            <el-input-number v-model="formData.balanceAmount" :precision="2" :controls="false" :disabled="true" style="width: 100%" />
          </el-form-item>
        </template>
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
        <el-form-item label="授课老师:" prop="teacherId">
          <el-select v-model="consumeForm.teacherId" placeholder="选择老师" filterable clearable style="width: 100%">
            <el-option
              v-for="teacher in teacherList"
              :key="teacher.ID"
              :label="`${teacher.nickName} (${teacher.phone || teacher.userName})`"
              :value="teacher.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="剩余课时:">
          <el-tag type="success" size="large">{{ consumeForm.remainingSessions }} 课时</el-tag>
        </el-form-item>
        <el-form-item label="消课数量:" prop="sessionsToConsume">
          <el-input-number v-model="consumeForm.sessionsToConsume" :min="1" :max="consumeForm.remainingSessions" :step="1" />
        </el-form-item>
        <template v-if="isSuperAdmin">
          <el-form-item label="课时单价:">
            <el-input-number v-model="consumeForm.unitPrice" :precision="2" :controls="false" :disabled="true" style="width: 100%" />
          </el-form-item>
          <el-form-item label="是否计费:">
            <el-switch v-model="consumeForm.chargeable" />
          </el-form-item>
          <el-form-item label="本次金额:">
            <el-input-number v-model="consumeForm.amount" :precision="2" :controls="false" :disabled="true" style="width: 100%" />
          </el-form-item>
        </template>
        <el-form-item label="上课时间:" prop="useDate">
          <el-date-picker v-model="consumeForm.useDate" type="datetime" placeholder="选择时间" value-format="YYYY-MM-DD HH:mm" style="width: 100%" />
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
        <el-form-item label="授课老师:" prop="teacherId">
          <el-select v-model="addForm.teacherId" placeholder="选择老师" filterable clearable style="width: 100%">
            <el-option
              v-for="teacher in teacherList"
              :key="teacher.ID"
              :label="`${teacher.nickName} (${teacher.phone || teacher.userName})`"
              :value="teacher.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="当前剩余:">
          <el-tag type="info" size="large">{{ addForm.remainingSessions }} 课时</el-tag>
        </el-form-item>
        <el-form-item label="增加数量:" prop="sessionsToAdd">
          <el-input-number v-model="addForm.sessionsToAdd" :min="1" :max="999" :step="1" />
        </el-form-item>
        <el-form-item label="日期:" prop="useDate">
          <el-date-picker v-model="addForm.useDate" type="datetime" placeholder="选择时间" value-format="YYYY-MM-DD HH:mm" style="width: 100%" />
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

    <!-- 收款记录弹窗（仅超管） -->
    <el-dialog v-if="isSuperAdmin" v-model="paymentDialogVisible" :before-close="closePaymentDialog" title="收款记录" width="700px">
      <el-form :model="paymentForm" label-position="right" ref="paymentFormRef" :rules="paymentRules" label-width="90px">
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="收款金额" prop="amount">
              <el-input-number v-model="paymentForm.amount" :min="0.01" :precision="2" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="收款时间" prop="payTime">
              <el-date-picker v-model="paymentForm.payTime" type="datetime" value-format="YYYY-MM-DD HH:mm" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="备注">
              <el-input v-model="paymentForm.remark" placeholder="可选" />
            </el-form-item>
          </el-col>
        </el-row>
        <div style="text-align: right; margin-bottom: 12px;">
          <el-button type="primary" @click="confirmPayment">记一笔收款</el-button>
        </div>
      </el-form>
      <el-table :data="paymentList" style="width: 100%">
        <el-table-column label="收款时间" min-width="160">
          <template #default="scope">{{ formatDateTime(scope.row.payTime) }}</template>
        </el-table-column>
        <el-table-column label="金额" prop="amount" width="120">
          <template #default="scope">{{ formatMoney(scope.row.amount) }}</template>
        </el-table-column>
        <el-table-column label="操作人" prop="operatorName" width="120" />
        <el-table-column label="备注" prop="remark" />
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="paymentPage"
          :page-size="paymentPageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="paymentTotal"
          @current-change="handlePaymentPageChange"
          @size-change="handlePaymentSizeChange"
        />
      </div>
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
import { createEduPayment, getEduPaymentList } from '@/api/eduPayment'
import { getEduCourseList } from '@/api/eduCourse'
import { getUserList } from '@/api/user'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'

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

// 辅助函数：格式化日期时间为 YYYY-MM-DD HH:mm
const formatDateTime = (date) => {
  const value = typeof date === 'string' ? date.replace(' ', 'T') : date
  const d = new Date(value)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        userId: 0,
        courseId: 0,
        totalSessions: 0,
        remainingSessions: 0,
        pricePerSession: 0,
        discountAmount: 0,
        totalAmount: 0,
        paidAmount: 0,
        balanceAmount: 0
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
  ],
  pricePerSession: [
    { type: 'number', min: 0, message: '课时单价不能为负数', trigger: 'blur' }
  ],
  discountAmount: [
    { type: 'number', min: 0, message: '优惠金额不能为负数', trigger: 'blur' }
  ]
})

const elFormRef = ref()

const recalcEnrollmentAmounts = () => {
  const totalSessions = Number(formData.value.totalSessions || 0)
  const price = Number(formData.value.pricePerSession || 0)
  const discount = Number(formData.value.discountAmount || 0)
  let total = totalSessions * price - discount
  if (total < 0) total = 0
  formData.value.totalAmount = Number(total.toFixed(2))
  const paid = Number(formData.value.paidAmount || 0)
  formData.value.balanceAmount = Number((formData.value.totalAmount - paid).toFixed(2))
}

watch(
  () => [formData.value.totalSessions, formData.value.pricePerSession, formData.value.discountAmount, formData.value.paidAmount],
  () => recalcEnrollmentAmounts()
)


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

// 用户列表、教师列表和课程列表
const userList = ref([])
const teacherList = ref([])
const courseList = ref([])

// 获取用户列表和课程列表
const setOptions = async () =>{
  // 获取所有用户（不分页）
  const userRes = await getUserList({ page: 1, pageSize: 9999 })
  if (userRes.code === 0) {
    const allUsers = userRes.data.list || []
    // 先排除管理员/系统角色，避免出现在学员下拉中
    const nonAdminUsers = allUsers.filter(user => {
      const userRoles = user.authorities || []
      const isAdminRole = userRoles.some(auth =>
        auth.authorityName?.includes('超级管理员') ||
        auth.authorityName?.includes('管理员') ||
        auth.authorityId === 888 ||
        auth.authorityId === 8881
      )
      const isAdminDefault = user.authorityId === 888 || user.authorityId === 8881 || user.authority_id === 888 || user.authority_id === 8881
      return !isAdminRole && !isAdminDefault
    })

    // 只显示学员角色的用户（角色ID为9002，或角色名包含"学员"/"学生"）
    userList.value = nonAdminUsers.filter(user => {
      if (user.authorityId === 9002 || user.authority_id === 9002) return true
      const userRoles = user.authorities || []
      return userRoles.some(auth => 
        auth.authorityId === 9002 ||  // 精确匹配学员角色ID
        auth.authorityName?.includes('学员') || 
        auth.authorityName?.includes('学生')
      )
    })

    // 只显示教师角色的用户（角色ID为9001，或角色名包含"教师"/"老师"）
    teacherList.value = nonAdminUsers.filter(user => {
      if (user.authorityId === 9001 || user.authority_id === 9001) return true
      const userRoles = user.authorities || []
      return userRoles.some(auth =>
        auth.authorityId === 9001 ||  // 精确匹配教师角色ID
        auth.authorityName?.includes('教师') ||
        auth.authorityName?.includes('老师')
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

const getTeacherNameById = (id) => {
  if (!id) return ''
  const teacher = teacherList.value.find(t => t.ID === id || t.id === id)
  return teacher?.nickName || teacher?.userName || ''
}


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
        recalcEnrollmentAmounts()
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
        pricePerSession: 0,
        discountAmount: 0,
        totalAmount: 0,
        paidAmount: 0,
        balanceAmount: 0
    }
    recalcEnrollmentAmounts()
    dialogFormVisible.value = true
}

// 总课时数改变时，自动同步剩余课时（仅新建时）
const onTotalSessionsChange = (val) => {
  if (type.value === 'create') {
    formData.value.remainingSessions = val || 0
  }
  recalcEnrollmentAmounts()
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        userId: 0,
        courseId: 0,
        totalSessions: 0,
        remainingSessions: 0,
        pricePerSession: 0,
        discountAmount: 0,
        totalAmount: 0,
        paidAmount: 0,
        balanceAmount: 0
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
  teacherId: 0,
  teacherName: '',
  remainingSessions: 0,
  sessionsToConsume: 1,
  unitPrice: 0,
  chargeable: true,
  amount: 0,
  useDate: formatDateTime(new Date()),
  reasonType: '',  // 原因类型（下拉选择）
  customReason: ''  // 自定义原因（仅当reasonType为"自定义"时使用）
})

const consumeRules = reactive({
  sessionsToConsume: [
    { required: true, message: '请输入消课数量', trigger: 'blur' },
    { type: 'number', min: 1, message: '消课数量至少为1', trigger: 'blur' }
  ],
  useDate: [
    { required: true, message: '请选择上课时间', trigger: 'change' }
  ],
  reasonType: [
    { required: true, message: '请选择原因', trigger: 'change' }
  ],
  customReason: [
    { required: true, message: '请输入详细原因', trigger: 'blur' }
  ]
})

const isNonChargeableReason = (reasonType) => {
  if (!reasonType) return false
  if (reasonType === '自定义') return true
  return ['试听', '赠课', '补课'].some((key) => reasonType.includes(key))
}

const recalcConsumeAmount = () => {
  const sessions = Number(consumeForm.value.sessionsToConsume || 0)
  const unitPrice = Number(consumeForm.value.unitPrice || 0)
  if (!consumeForm.value.chargeable) {
    consumeForm.value.amount = 0
    return
  }
  consumeForm.value.amount = Number((sessions * unitPrice).toFixed(2))
}

// 消课原因类型改变时的处理
const onConsumeReasonChange = (value) => {
  // 如果不是自定义，清空自定义原因
  if (value !== '自定义') {
    consumeForm.value.customReason = ''
  }
  consumeForm.value.chargeable = !isNonChargeableReason(value)
  recalcConsumeAmount()
}

watch(
  () => [consumeForm.value.sessionsToConsume, consumeForm.value.unitPrice, consumeForm.value.chargeable],
  () => recalcConsumeAmount()
)

// 打开消课弹窗
const openConsumeDialog = (row) => {
  consumeForm.value = {
    userId: row.userId,
    courseId: row.courseId,
    userName: row.userName,
    courseName: row.eduCourse?.courseName || '',
    teacherId: userStore.userInfo?.ID || 0,
    teacherName: userStore.userInfo?.nickName || '',
    remainingSessions: row.remainingSessions || 0,
    sessionsToConsume: 1,
    unitPrice: row.pricePerSession || 0,
    chargeable: true,
    amount: 0,
    useDate: formatDateTime(new Date()),
    reasonType: '',
    customReason: ''
  }
  recalcConsumeAmount()
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
    
    const teacherName = getTeacherNameById(consumeForm.value.teacherId) || consumeForm.value.teacherName || ''
    const params = {
      userId: consumeForm.value.userId,
      courseId: consumeForm.value.courseId,
      sessionsToConsume: consumeForm.value.sessionsToConsume,
      useDate: consumeForm.value.useDate,
      reason: finalReason,
      teacherId: consumeForm.value.teacherId || 0,
      teacherName,
      chargeable: consumeForm.value.chargeable
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
  teacherId: 0,
  teacherName: '',
  remainingSessions: 0,
  sessionsToAdd: 1,
  useDate: formatDateTime(new Date()),
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
    teacherId: userStore.userInfo?.ID || 0,
    teacherName: userStore.userInfo?.nickName || '',
    remainingSessions: row.remainingSessions || 0,
    sessionsToAdd: 1,
    useDate: formatDateTime(new Date()),
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
    
    const teacherName = getTeacherNameById(addForm.value.teacherId) || addForm.value.teacherName || ''
    const params = {
      userId: addForm.value.userId,
      courseId: addForm.value.courseId,
      sessionsToAdd: addForm.value.sessionsToAdd,
      useDate: addForm.value.useDate,
      reason: finalReason,
      teacherId: addForm.value.teacherId || 0,
      teacherName
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

// ============== 收款记录 ==============
const paymentDialogVisible = ref(false)
const paymentFormRef = ref()
const paymentForm = ref({
  amount: 0,
  payTime: formatDateTime(new Date()),
  remark: ''
})
const paymentRules = reactive({
  amount: [
    { required: true, message: '请输入收款金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '收款金额必须大于0', trigger: 'blur' }
  ],
  payTime: [
    { required: true, message: '请选择收款时间', trigger: 'change' }
  ]
})
const paymentList = ref([])
const paymentTotal = ref(0)
const paymentPage = ref(1)
const paymentPageSize = ref(10)
const paymentTarget = ref(null)

const openPaymentDialog = (row) => {
  paymentTarget.value = row
  paymentPage.value = 1
  paymentPageSize.value = 10
  paymentForm.value = {
    amount: 0,
    payTime: formatDateTime(new Date()),
    remark: ''
  }
  paymentDialogVisible.value = true
  loadPaymentList()
}

const closePaymentDialog = () => {
  paymentDialogVisible.value = false
  paymentList.value = []
  paymentTarget.value = null
}

const loadPaymentList = async () => {
  if (!paymentTarget.value?.ID) return
  const res = await getEduPaymentList({
    enrollmentId: paymentTarget.value.ID,
    page: paymentPage.value,
    pageSize: paymentPageSize.value
  })
  if (res.code === 0) {
    paymentList.value = res.data.list || []
    paymentTotal.value = res.data.total || 0
  }
}

const confirmPayment = async () => {
  paymentFormRef.value?.validate(async (valid) => {
    if (!valid || !paymentTarget.value?.ID) return
    const res = await createEduPayment({
      enrollmentId: paymentTarget.value.ID,
      amount: paymentForm.value.amount,
      payTime: paymentForm.value.payTime,
      remark: paymentForm.value.remark
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '收款记录已保存'
      })
      paymentForm.value.amount = 0
      paymentForm.value.remark = ''
      paymentForm.value.payTime = formatDateTime(new Date())
      loadPaymentList()
      getTableData()
    }
  })
}

const handlePaymentPageChange = (val) => {
  paymentPage.value = val
  loadPaymentList()
}

const handlePaymentSizeChange = (val) => {
  paymentPageSize.value = val
  loadPaymentList()
}
</script>

<style>
</style>
