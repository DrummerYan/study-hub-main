<template>
  <div>
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-table-box">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="教师管理" name="teacher">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="addUser">新增教师</el-button>
            <el-button type="primary" icon="plus" @click="importUser">导入教师</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane label="学员管理" name="student">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="addUser">新增学员</el-button>
            <el-button type="primary" icon="plus" @click="importUser">导入学员</el-button>
          </div>
        </el-tab-pane>
        <el-tab-pane label="所有用户" name="all">
          <div class="gva-btn-list">
            <el-alert 
              title="此标签显示所有非管理员用户，包括可能缺失角色的用户" 
              type="info" 
              show-icon 
              :closable="false"
              style="margin-bottom: 10px;"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
      <el-table :data="filteredTableData" row-key="ID">
        <el-table-column align="left" label="头像" min-width="75">
          <template #default="scope">
            <CustomPic style="margin-top:8px" :pic-src="scope.row.headerImg" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="ID" min-width="50" prop="ID" />
        <el-table-column align="left" label="用户名" min-width="150" prop="userName" />
        <el-table-column align="left" label="昵称" min-width="150" prop="nickName" />
        <el-table-column align="left" label="手机号" min-width="180" prop="phone" />
        <el-table-column align="left" label="邮箱" min-width="180" prop="email" />
        <el-table-column align="left" label="用户角色" min-width="200">
          <template #default="scope">
            <div style="display: flex; align-items: center; gap: 8px;">
              <el-cascader v-model="scope.row.authorityIds" :options="authOptions" :show-all-levels="false" collapse-tags
                :props="{ multiple: true, checkStrictly: true, label: 'authorityName', value: 'authorityId', disabled: 'disabled', emitPath: false }"
                :clearable="false" @visible-change="(flag) => { changeAuthority(scope.row, flag, 0) }"
                @remove-tag="(removeAuth) => { changeAuthority(scope.row, false, removeAuth) }" />
              <el-tooltip 
                v-if="!scope.row.authorityIds || scope.row.authorityIds.length === 0"
                content="⚠️ 该用户缺少角色，可能因为角色被删除导致"
                placement="top"
              >
                <el-icon color="#F56C6C" :size="20"><WarningFilled /></el-icon>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="所属组织" min-width="200">
          <template #default="scope">
            <el-select v-model="scope.row.eduOrganizationID" style="width:100%">
              <el-option v-for="item in organizationOptions" :key="item.ID" :label="item.name" :value="item.ID" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column align="left" label="启用" min-width="150">
          <template #default="scope">
            <el-switch v-model="scope.row.enable" inline-prompt :active-value="1" :inactive-value="2"
              @change="() => { switchEnable(scope.row) }" />
          </template>
        </el-table-column>

        <el-table-column label="操作" min-width="250" fixed="right">
          <template #default="scope">
            <el-popover v-model="scope.row.visible" placement="top" width="160" trigger="click">
              <p>确定要删除此用户吗</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="scope.row.visible = false">取消</el-button>
                <el-button type="primary" @click="deleteUserFunc(scope.row)">确定</el-button>
              </div>
              <template #reference>
                <el-button type="primary" link icon="delete">删除</el-button>
              </template>
            </el-popover>
            <el-button type="primary" link icon="edit" @click="openEdit(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="magic-stick" @click="resetPasswordFunc(scope.row)">重置密码</el-button>
          </template>
        </el-table-column>

      </el-table>
      <div class="gva-pagination">
        <el-pagination :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100]" :total="total"
          layout="total, sizes, prev, pager, next, jumper" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-dialog v-model="addUserDialog" custom-class="user-dialog" title="用户" :show-close="false"
      :close-on-press-escape="false" :close-on-click-modal="false">
      <div style="height:60vh;overflow:auto;padding:0 12px;">
        <el-form ref="userForm" :rules="rules" :model="userInfo" label-width="80px">
          <el-form-item v-if="dialogFlag === 'add'" label="用户名" prop="userName">
            <el-input v-model="userInfo.userName" />
          </el-form-item>
          <el-form-item v-if="dialogFlag === 'add'" label="密码" prop="password">
            <el-input v-model="userInfo.password" />
          </el-form-item>
          <el-form-item label="昵称" prop="nickName">
            <el-input v-model="userInfo.nickName" />
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="userInfo.phone" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="userInfo.email" />
          </el-form-item>
          <!-- 用户角色：根据当前标签自动确定 -->
          <el-form-item label="用户角色">
            <el-tag v-if="activeTab === 'teacher'" type="success" size="large">
              教师（自动设置）
            </el-tag>
            <el-tag v-else-if="activeTab === 'student'" type="primary" size="large">
              学员（自动设置）
            </el-tag>
            <div style="color: #999; font-size: 12px; margin-top: 4px;">
              根据当前标签自动设置为{{ activeTab === 'teacher' ? '教师' : '学员' }}角色
            </div>
          </el-form-item>
          <el-form-item label="所属组织" prop="eduOrganizationID">
            <el-select v-model="userInfo.eduOrganizationID" style="width:100%">
              <el-option v-for="item in organizationOptions" :key="item.ID" :label="item.name" :value="item.ID" />
            </el-select>

          </el-form-item>
          <el-form-item label="启用" prop="disabled">
            <el-switch v-model="userInfo.enable" inline-prompt :active-value="1" :inactive-value="2" />
          </el-form-item>
          <el-form-item label="头像" label-width="80px">
            <div style="display:inline-block" @click="openHeaderChange">
              <img v-if="userInfo.headerImg" alt="头像" class="header-img-box"
                :src="(userInfo.headerImg && userInfo.headerImg.slice(0, 4) !== 'http') ? path + userInfo.headerImg : userInfo.headerImg">
              <div v-else class="header-img-box">从媒体库选择</div>
            </div>
          </el-form-item>

        </el-form>

      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeAddUserDialog">取 消</el-button>
          <el-button type="primary" @click="enterAddUserDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="importUserDialog" custom-class="user-dialog" title="用户" :show-close="false"
      :close-on-press-escape="false" :close-on-click-modal="false">
      <div style="height:30vh;overflow:auto;padding:0 12px;">
        <el-form ref="importUserForm" :rules="rules" :model="importUserInfo" label-width="80px">
          <!-- 用户角色：根据当前标签自动确定 -->
          <el-form-item label="用户角色">
            <el-tag v-if="activeTab === 'teacher'" type="success" size="large">
              教师（自动设置）
            </el-tag>
            <el-tag v-else-if="activeTab === 'student'" type="primary" size="large">
              学员（自动设置）
            </el-tag>
          </el-form-item>
          <el-form-item label="所属组织" prop="eduOrganizationID">
            <el-select v-model="importUserInfo.eduOrganizationID" style="width:100%">
              <el-option v-for="item in organizationOptions" :key="item.ID" :label="item.name" :value="item.ID" />
            </el-select>
          </el-form-item>
          <el-form-item label="导入文件" prop="importFile">
            <el-upload class="upload-demo" drag :action="`${path}/fileUploadAndDownload/upload`"
              :headers="{ 'x-token': userStore.token }" :on-success="handleImageSuccess"
              :before-upload="beforeImageUpload" multiple:false>
              <i class="el-icon-upload"></i>
              <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
              <div class="el-upload__tip" slot="tip">只能上传xlsx文件，且不超过500kb</div>
            </el-upload>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeImportUser">取 消</el-button>
          <el-button type="primary" @click="enterImportUserDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    <ChooseImg ref="chooseImg" :target="userInfo" :target-key="`headerImg`" />
  </div>
</template>

<script>
export default {
  name: 'User',
}
</script>

<script setup>

import {
  getUserList,
  setUserAuthorities,
  register,
  deleteUser,
  importExcelUser
} from '@/api/user'

import { getAuthorityList } from '@/api/authority'
import { getEduOrganizationList } from '@/api/eduOrganization'
import CustomPic from '@/components/customPic/index.vue'
import ChooseImg from '@/components/chooseImg/index.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { setUserInfo, resetPassword } from '@/api/user.js'
import { useUserStore } from '@/pinia/modules/user'
import { WarningFilled } from '@element-plus/icons-vue'

import { nextTick, ref, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
const path = ref(import.meta.env.VITE_BASE_API + '/')

const userStore = useUserStore()

// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
    AuthorityData.forEach(item => {
      if (item.children && item.children.length) {
        const option = {
          authorityId: item.authorityId,
          authorityName: item.authorityName,
          children: []
        }
        setAuthorityOptions(item.children, option.children)
        optionsData.push(option)
      } else {
        const option = {
          authorityId: item.authorityId,
          authorityName: item.authorityName
        }
        optionsData.push(option)
      }
    })
}

const handleImageSuccess = (res) => {
  const { data } = res
  if (data.file) {
    importUserInfo.value.importFile = data.file.url
  }
}

const beforeImageUpload = (file) => {
  const isXlsx = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
  if (!isXlsx) {
    ElMessage.error('只能是 xlsx 格式!')
    return false
  }
}


const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const activeTab = ref('teacher') // 默认显示教师管理

// 根据标签过滤表格数据
const filteredTableData = computed(() => {
  // 过滤掉系统管理员
  const nonAdminUsers = tableData.value.filter(user => {
    const userRoles = user.authorities || []
    // 排除超级管理员、管理员等系统角色
    const isAdmin = userRoles.some(auth => 
      auth.authorityName?.includes('超级管理员') ||
      auth.authorityName?.includes('管理员') ||
      auth.authorityId === 888 ||
      auth.authorityId === 8881
    )
    return !isAdmin
  })
  
  if (activeTab.value === 'teacher') {
    // 教师标签：显示角色ID为9001的用户，或角色名包含"教师"
    return nonAdminUsers.filter(user => {
      const userRoles = user.authorities || []
      return userRoles.some(auth => 
        auth.authorityId === 9001 ||  // 精确匹配教师角色ID
        auth.authorityName?.includes('教师') || 
        auth.authorityName?.includes('老师')
      )
    })
  } else if (activeTab.value === 'student') {
    // 学员标签：显示角色ID为9002的用户，或角色名包含"学员"
    return nonAdminUsers.filter(user => {
      const userRoles = user.authorities || []
      return userRoles.some(auth => 
        auth.authorityId === 9002 ||  // 精确匹配学员角色ID
        auth.authorityName?.includes('学员') || 
        auth.authorityName?.includes('学生')
      )
    })
  } else if (activeTab.value === 'all') {
    // 所有用户标签：显示所有非管理员用户
    return nonAdminUsers
  }
  
  return nonAdminUsers
})

// 标签切换处理
const handleTabChange = (tab) => {
  activeTab.value = tab
  // 切换标签时可以重新加载数据或只是前端过滤
  // getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getUserList({ page: page.value, pageSize: pageSize.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

watch(() => tableData.value, () => {
  setAuthorityIds()
})

const initPage = async () => {
  getTableData()
  const res = await getAuthorityList({ page: 1, pageSize: 999 })
  setOptions(res.data.list)

  const organization = await getEduOrganizationList({ page: 1, pageSize: 999 })
  setOrganization(organization.data.list)
}

initPage()

const resetPasswordFunc = (row) => {
  ElMessageBox.confirm(
    '是否将此用户密码重置为123456?',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    const res = await resetPassword({
      ID: row.ID,
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}
const setAuthorityIds = () => {
  tableData.value && tableData.value.forEach((user) => {
    user.authorityIds = user.authorities && user.authorities.map(i => {
      return i.authorityId
    })
  })
}

const chooseImg = ref(null)
const openHeaderChange = () => {
  chooseImg.value.open()
}

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  setAuthorityOptions(authData, authOptions.value)
}

const organizationOptions = ref([])
const setOrganization = (organizationData) => {
  organizationOptions.value = organizationData
  console.log();
}

const deleteUserFunc = async (row) => {
  const res = await deleteUser({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    row.visible = false
    await getTableData()
  }
}

// 弹窗相关
const userInfo = ref({
  username: '',
  password: '',
  nickName: '',
  headerImg: '',
  authorityId: '',
  authorityIds: [],
  enable: 1,
})

const importUserInfo = ref({
  authorityId: '',
  authorityIds: [],
  eduOrganizationID: '',
  importFile: '',
})

const rules = ref({
  userName: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 5, message: '最低5位字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入用户密码', trigger: 'blur' },
    { min: 6, message: '最低6位字符', trigger: 'blur' }
  ],
  nickName: [
    { required: true, message: '请输入用户昵称', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1([38][0-9]|4[014-9]|[59][0-35-9]|6[2567]|7[0-8])\d{8}$/, message: '请输入合法手机号', trigger: 'blur' },
  ],
  email: [
    { pattern: /^([0-9A-Za-z\-_.]+)@([0-9a-z]+\.[a-z]{2,3}(\.[a-z]{2})?)$/g, message: '请输入正确的邮箱', trigger: 'blur' },
  ],
  authorityId: [
    { required: true, message: '请选择用户角色', trigger: 'blur' }
  ],
  eduOrganizationID: [
    { required: true, message: '请选择用户组织', trigger: 'blur' }
  ],
  importFile: [
    { required: true, message: '请选择文件', trigger: 'blur' }
  ]
})
const userForm = ref(null)
const enterAddUserDialog = async () => {
  userInfo.value.authorityId = userInfo.value.authorityIds[0]
  userForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...userInfo.value
      }
      if (dialogFlag.value === 'add') {
        const res = await register(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '创建成功' })
          await getTableData()
          closeAddUserDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await setUserInfo(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '编辑成功' })
          await getTableData()
          closeAddUserDialog()
        }
      }
    }
  })
}

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  userForm.value.resetFields()
  userInfo.value.headerImg = ''
  userInfo.value.authorityIds = []
  addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
  dialogFlag.value = 'add'
  
  // 根据当前标签自动设置角色
  if (activeTab.value === 'teacher') {
    // 教师管理标签：自动设置为教师角色（9001）
    userInfo.value.authorityIds = [9001]
  } else if (activeTab.value === 'student') {
    // 学员管理标签：自动设置为学员角色（9002）
    userInfo.value.authorityIds = [9002]
  }
  
  addUserDialog.value = true
}

const importUserDialog = ref(false)
const importUser = () => {
  dialogFlag.value = 'import'
  
  // 根据当前标签自动设置角色
  if (activeTab.value === 'teacher') {
    importUserInfo.value.authorityIds = [9001]
  } else if (activeTab.value === 'student') {
    importUserInfo.value.authorityIds = [9002]
  }
  
  importUserDialog.value = true
}
const closeImportUser = () => {
  importUserDialog.value = false
}
//导入用户
const importUserForm = ref(null)
const enterImportUserDialog = () => {
  importUserInfo.value.authorityId = importUserInfo.value.authorityIds[0]
  importUserForm.value.validate(async valid => {
    if (valid) {
      const req = {
        ...importUserInfo.value
      }
      const res = await importExcelUser(req)
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '创建成功' })
        await getTableData()
        closeImportUser()
      }
    }
  })
}

const tempAuth = {}
const changeAuthority = async (row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds]
    }
    return
  }
  await nextTick()
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds
  })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '角色设置成功' })
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds]
    }
  }
}

const openEdit = (row) => {
  dialogFlag.value = 'edit'
  userInfo.value = JSON.parse(JSON.stringify(row))
  addUserDialog.value = true
}

const switchEnable = async (row) => {
  userInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...userInfo.value
  }
  const res = await setUserInfo(req)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: `${req.enable === 2 ? '禁用' : '启用'}成功` })
    await getTableData()
    userInfo.value.headerImg = ''
    userInfo.value.authorityIds = []
  }
}

</script>

<style lang="scss">
.user-dialog {
  .header-img-box {
    width: 200px;
    height: 200px;
    border: 1px dashed #ccc;
    border-radius: 20px;
    text-align: center;
    line-height: 200px;
    cursor: pointer;
  }

  .avatar-uploader .el-upload:hover {
    border-color: #409eff;
  }

  .avatar-uploader-icon {
    border: 1px dashed #d9d9d9 !important;
    border-radius: 6px;
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }

  .avatar {
    width: 178px;
    height: 178px;
    display: block;
  }
}

.nickName {
  display: flex;
  justify-content: flex-start;
  align-items: center;
}

.pointer {
  cursor: pointer;
  font-size: 16px;
  margin-left: 2px;
}
</style>
