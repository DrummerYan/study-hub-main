<template>
  <div class="authority">
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="addAuthority(0)">新增角色</el-button>
        <el-button type="success" icon="data-analysis" @click="showOptimizationDialog">角色优化建议</el-button>
      </div>
      <el-table
        :data="tableData"
        :tree-props="{children: 'children', hasChildren: 'hasChildren'}"
        row-key="authorityId"
        style="width: 100%"
      >
        <el-table-column label="角色ID" min-width="180" prop="authorityId" />
        <el-table-column align="left" label="角色名称" min-width="180" prop="authorityName" />
        <el-table-column align="left" label="操作" width="560">
          <template #default="scope">
            <el-button
              icon="setting"
              type="primary"
              link
              @click="opdendrawer(scope.row)"
            >设置权限</el-button>
            <el-button
              icon="plus"
              type="primary"
              link
              @click="addAuthority(scope.row.authorityId)"
            >新增子角色</el-button>
            <el-button
              icon="view"
              type="primary"
              link
              @click="viewUsageInfo(scope.row)"
            >使用情况</el-button>
            <el-button
              icon="copy-document"
              type="primary"
              link
              @click="copyAuthorityFunc(scope.row)"
            >拷贝</el-button>
            <el-button
              icon="edit"
              type="primary"
              link
              @click="editAuthority(scope.row)"
            >编辑</el-button>
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteAuth(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    <!-- 新增角色弹窗 -->
    <el-dialog v-model="dialogFormVisible" :title="dialogTitle">
      <el-form ref="authorityForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="父级角色" prop="parentId">
          <el-cascader
            v-model="form.parentId"
            style="width:100%"
            :disabled="dialogType==='add'"
            :options="AuthorityOption"
            :props="{ checkStrictly: true,label:'authorityName',value:'authorityId',disabled:'disabled',emitPath:false}"
            :show-all-levels="false"
            filterable
          />
        </el-form-item>
        <el-form-item label="角色ID" prop="authorityId">
          <el-input v-model="form.authorityId" :disabled="dialogType==='edit'" autocomplete="off" />
        </el-form-item>
        <el-form-item label="角色姓名" prop="authorityName">
          <el-input v-model="form.authorityName" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-drawer v-if="drawer" v-model="drawer" custom-class="auth-drawer" :with-header="false" size="40%" title="角色配置">
      <el-tabs :before-leave="autoEnter" type="border-card">
        <el-tab-pane label="角色菜单">
          <Menus ref="menus" :row="activeRow" @changeRow="changeRow" />
        </el-tab-pane>
        <el-tab-pane label="角色api">
          <Apis ref="apis" :row="activeRow" @changeRow="changeRow" />
        </el-tab-pane>
        <el-tab-pane label="资源权限">
          <Datas ref="datas" :authority="tableData" :row="activeRow" @changeRow="changeRow" />
        </el-tab-pane>
      </el-tabs>
    </el-drawer>

    <!-- 角色使用情况对话框 -->
    <el-dialog v-model="usageDialogVisible" title="角色使用情况" width="700px">
      <div v-loading="usageLoading">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="角色ID">
            {{ currentUsageInfo.authorityInfo?.authorityId }}
          </el-descriptions-item>
          <el-descriptions-item label="角色名称">
            {{ currentUsageInfo.authorityInfo?.authorityName }}
          </el-descriptions-item>
          <el-descriptions-item label="使用该角色的用户数" :span="2">
            <el-tag :type="currentUsageInfo.multiUserCount > 0 ? 'danger' : 'success'">
              {{ currentUsageInfo.multiUserCount || 0 }} 人
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="以此为默认角色的用户" :span="2">
            <el-tag :type="currentUsageInfo.defaultUserCount > 0 ? 'warning' : 'success'">
              {{ currentUsageInfo.defaultUserCount || 0 }} 人
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="子角色数量" :span="2">
            <el-tag :type="currentUsageInfo.childrenCount > 0 ? 'info' : 'success'">
              {{ currentUsageInfo.childrenCount || 0 }} 个
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="是否可删除" :span="2">
            <el-tag :type="currentUsageInfo.canDelete ? 'success' : 'danger'">
              {{ currentUsageInfo.canDelete ? '✅ 可以删除' : '❌ 不可删除' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="!currentUsageInfo.canDelete" style="margin-top: 20px;">
          <el-alert type="error" :closable="false">
            <template #title>
              <strong>无法删除的原因：</strong>
            </template>
            <ul style="margin: 10px 0; padding-left: 20px;">
              <li v-for="(reason, index) in currentUsageInfo.blockReasons" :key="index">
                {{ reason }}
              </li>
            </ul>
          </el-alert>
        </div>

        <div v-if="currentUsageInfo.users && currentUsageInfo.users.length > 0" style="margin-top: 20px;">
          <el-divider content-position="left">使用该角色的用户（前10个）</el-divider>
          <el-table :data="currentUsageInfo.users" border stripe>
            <el-table-column prop="id" label="用户ID" width="80" />
            <el-table-column prop="username" label="用户名" min-width="120" />
            <el-table-column prop="nickName" label="昵称" min-width="120" />
            <el-table-column prop="phone" label="手机号" min-width="130" />
          </el-table>
        </div>
      </div>
      <template #footer>
        <el-button @click="usageDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 角色优化建议对话框 -->
    <el-dialog v-model="optimizationDialogVisible" title="🚀 角色优化建议" width="900px">
      <div style="max-height: 600px; overflow-y: auto;">
        <el-alert 
          type="info" 
          :closable="false" 
          style="margin-bottom: 20px;"
          title="基于您的系统分析，以下是角色优化建议"
        />

        <!-- 当前角色统计 -->
        <el-card shadow="never" style="margin-bottom: 20px;">
          <template #header>
            <strong>📊 当前角色统计</strong>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-statistic title="总角色数" :value="roleStats.totalRoles">
                <template #suffix>个</template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic title="核心角色" :value="roleStats.coreRoles">
                <template #suffix>个</template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic title="可清理角色" :value="roleStats.unusedRoles">
                <template #suffix>个</template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>

        <!-- 优化建议 -->
        <el-card shadow="never" style="margin-bottom: 20px;">
          <template #header>
            <strong>💡 优化建议</strong>
          </template>
          <el-timeline>
            <el-timeline-item 
              v-for="(suggestion, index) in optimizationSuggestions" 
              :key="index"
              :type="suggestion.type"
              :icon="suggestion.icon"
            >
              <div style="padding-bottom: 10px;">
                <strong>{{ suggestion.title }}</strong>
                <p style="margin: 8px 0; color: #606266;">{{ suggestion.description }}</p>
                <el-button 
                  v-if="suggestion.action"
                  size="small" 
                  :type="suggestion.actionType"
                  @click="handleOptimizationAction(suggestion.action, suggestion.data)"
                >
                  {{ suggestion.actionText }}
                </el-button>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- 推荐的扁平化角色结构 -->
        <el-card shadow="never">
          <template #header>
            <strong>🎯 推荐的扁平化角色结构</strong>
          </template>
          <el-table :data="recommendedRoles" border>
            <el-table-column prop="authorityId" label="角色ID" width="100" />
            <el-table-column prop="authorityName" label="角色名称" min-width="150" />
            <el-table-column prop="description" label="说明" min-width="200" />
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'exists' ? 'success' : 'info'">
                  {{ scope.row.status === 'exists' ? '已存在' : '建议新增' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
      <template #footer>
        <el-button @click="optimizationDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="exportOptimizationReport">导出优化报告</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getAuthorityList,
  deleteAuthority,
  createAuthority,
  updateAuthority,
  copyAuthority,
  getAuthorityUsageInfo
} from '@/api/authority'

import Menus from '@/view/superAdmin/authority/components/menus.vue'
import Apis from '@/view/superAdmin/authority/components/apis.vue'
import Datas from '@/view/superAdmin/authority/components/datas.vue'
import WarningBar from '@/components/warningBar/warningBar.vue'

import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Promotion, Edit } from '@element-plus/icons-vue'

const mustUint = (rule, value, callback) => {
  if (!/^[0-9]*[1-9][0-9]*$/.test(value)) {
    return callback(new Error('请输入正整数'))
  }
  return callback()
}

// 检查是否会造成循环引用
const checkCircularReference = (currentId, parentId, allData) => {
  // 如果父级是自己，直接返回 true（循环）
  if (currentId === parentId) {
    return true
  }
  
  // 递归检查父级链，看是否会回到自己
  const visited = new Set()
  let checkId = parentId
  
  while (checkId && checkId !== 0) {
    // 如果遇到自己，说明有循环
    if (checkId === currentId) {
      return true
    }
    
    // 如果已经访问过这个节点，说明有循环
    if (visited.has(checkId)) {
      return true
    }
    
    visited.add(checkId)
    
    // 查找这个角色的父级
    const role = findRoleInTree(allData, checkId)
    if (!role) break
    
    checkId = role.parentId
  }
  
  return false
}

// 在树形结构中查找角色
const findRoleInTree = (data, targetId) => {
  for (const item of data) {
    if (item.authorityId === targetId) {
      return item
    }
    if (item.children && item.children.length > 0) {
      const found = findRoleInTree(item.children, targetId)
      if (found) return found
    }
  }
  return null
}

const AuthorityOption = ref([
  {
    authorityId: 0,
    authorityName: '根角色'
  }
])
const drawer = ref(false)
const dialogType = ref('add')
const activeRow = ref({})

const dialogTitle = ref('新增角色')
const dialogFormVisible = ref(false)
const apiDialogFlag = ref(false)
const copyForm = ref({})

const form = ref({
  authorityId: 0,
  authorityName: '',
  parentId: 0
})
const rules = ref({
  authorityId: [
    { required: true, message: '请输入角色ID', trigger: 'blur' },
    { validator: mustUint, trigger: 'blur', message: '必须为正整数' }
  ],
  authorityName: [
    { required: true, message: '请输入角色名', trigger: 'blur' }
  ],
  parentId: [
    { required: true, message: '请选择父角色', trigger: 'blur' },
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(999)
const tableData = ref([])
const searchInfo = ref({})

// 查询
const getTableData = async() => {
  const table = await getAuthorityList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const changeRow = (key, value) => {
  activeRow.value[key] = value
}
const menus = ref(null)
const apis = ref(null)
const datas = ref(null)
const autoEnter = (activeName, oldActiveName) => {
  const paneArr = [menus, apis, datas]
  if (oldActiveName) {
    if (paneArr[oldActiveName].value.needConfirm) {
      paneArr[oldActiveName].value.enterAndNext()
      paneArr[oldActiveName].value.needConfirm = false
    }
  }
}
// 拷贝角色
const copyAuthorityFunc = (row) => {
  setOptions()
  dialogTitle.value = '拷贝角色'
  dialogType.value = 'copy'
  for (const k in form.value) {
    form.value[k] = row[k]
  }
  copyForm.value = row
  dialogFormVisible.value = true
}
const opdendrawer = (row) => {
  drawer.value = true
  activeRow.value = row
}
// 删除角色
const deleteAuth = (row) => {
  // 检查是否为核心角色（只保护888超级管理员和9001教师角色）
  const coreRoles = [888, 9001]
  if (coreRoles.includes(row.authorityId)) {
    ElMessageBox.alert(
      `角色 "${row.authorityName}" (ID: ${row.authorityId}) 是系统核心角色，禁止删除！`,
      '⚠️ 无法删除核心角色',
      {
        type: 'error',
        confirmButtonText: '知道了'
      }
    )
    return
  }
  
  // 检查是否有子角色
  const hasChildren = row.children && row.children.length > 0
  const childrenNames = hasChildren ? row.children.map(c => c.authorityName).join('、') : ''
  
  // 构建警告消息
  let warningMessage = `<div style="line-height: 1.8;">
    <p><strong>即将删除角色：</strong>${row.authorityName} (ID: ${row.authorityId})</p>
  `
  
  if (hasChildren) {
    warningMessage += `
    <p style="color: #E6A23C;">
      <strong>⚠️ 警告：</strong>此角色下有 ${row.children.length} 个子角色<br/>
      子角色：${childrenNames}<br/>
      <strong>必须先删除所有子角色才能删除此角色！</strong>
    </p>
    `
  }
  
  warningMessage += `
    <p style="color: #F56C6C;">
      <strong>⚠️ 重要提示：</strong><br/>
      1. 如果有用户正在使用此角色，删除将会失败<br/>
      2. 删除后无法恢复，请谨慎操作<br/>
      3. 建议先在"师生管理"中检查该角色的用户
    </p>
  </div>
  `
  
  ElMessageBox.confirm(warningMessage, '⚠️ 确认删除角色', {
    confirmButtonText: hasChildren ? '无法删除' : '确定删除',
    cancelButtonText: '取消',
    type: 'warning',
    dangerouslyUseHTMLString: true,
    distinguishCancelAndClose: true,
    confirmButtonClass: hasChildren ? 'is-disabled' : '',
    beforeClose: (action, instance, done) => {
      if (action === 'confirm' && hasChildren) {
        ElMessage({
          type: 'warning',
          message: '请先删除所有子角色！'
        })
        done()
        return false
      }
      done()
    }
  })
    .then(async() => {
      if (hasChildren) {
        return // 有子角色时不执行删除
      }
      
      const res = await deleteAuthority({ authorityId: row.authorityId })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
    .catch(() => {
      // 用户取消删除，不显示提示
    })
}
// 初始化表单
const authorityForm = ref(null)
const initForm = () => {
  if (authorityForm.value) {
    authorityForm.value.resetFields()
  }
  form.value = {
    authorityId: 0,
    authorityName: '',
    parentId: 0
  }
}
// 关闭窗口
const closeDialog = () => {
  initForm()
  dialogFormVisible.value = false
  apiDialogFlag.value = false
}
// 确定弹窗

const enterDialog = () => {
  form.value.authorityId = Number(form.value.authorityId)
  if (form.value.authorityId === 0) {
    ElMessage({
      type: 'error',
      message: '角色id不能为0'
    })
    return false
  }
  
  // 防止循环引用：检查完整的父级链
  if (checkCircularReference(form.value.authorityId, form.value.parentId, tableData.value)) {
    ElMessage({
      type: 'error',
      message: '⚠️ 不能选择自己或自己的子孙角色作为父级！这会导致循环引用，角色将无法显示。',
      duration: 5000
    })
    return false
  }
  
  authorityForm.value.validate(async valid => {
    if (valid) {
      switch (dialogType.value) {
        case 'add':
          {
            const res = await createAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '添加成功!'
              })
              getTableData()
              closeDialog()
            }
          }
          break
        case 'edit':
          {
            const res = await updateAuthority(form.value)
            if (res.code === 0) {
              ElMessage({
                type: 'success',
                message: '添加成功!'
              })
              getTableData()
              closeDialog()
            }
          }
          break
        case 'copy': {
          const data = {
            authority: {
              authorityId: 0,
              authorityName: '',
              datauthorityId: [],
              parentId: 0
            },
            oldAuthorityId: 0
          }
          data.authority.authorityId = form.value.authorityId
          data.authority.authorityName = form.value.authorityName
          data.authority.parentId = form.value.parentId
          data.authority.dataAuthorityId = copyForm.value.dataAuthorityId
          data.oldAuthorityId = copyForm.value.authorityId
          const res = await copyAuthority(data)
          if (res.code === 0) {
            ElMessage({
              type: 'success',
              message: '复制成功！'
            })
            getTableData()
          }
        }
      }

      initForm()
      dialogFormVisible.value = false
    }
  })
}
const setOptions = () => {
  AuthorityOption.value = [
    {
      authorityId: 0,
      authorityName: '根角色'
    }
  ]
  setAuthorityOptions(tableData.value, AuthorityOption.value, false)
}

// ========== 角色使用情况功能 ==========
const usageDialogVisible = ref(false)
const usageLoading = ref(false)
const currentUsageInfo = ref({})

const viewUsageInfo = async (row) => {
  usageDialogVisible.value = true
  usageLoading.value = true
  currentUsageInfo.value = {}
  
  try {
    const res = await getAuthorityUsageInfo({ authorityId: row.authorityId })
    if (res.code === 0) {
      currentUsageInfo.value = res.data
    } else {
      ElMessage.error(res.msg || '获取角色使用情况失败')
    }
  } catch (error) {
    ElMessage.error('获取角色使用情况失败：' + error.message)
  } finally {
    usageLoading.value = false
  }
}

// ========== 角色优化建议功能 ==========
const optimizationDialogVisible = ref(false)
const roleStats = ref({
  totalRoles: 0,
  coreRoles: 0,
  unusedRoles: 0
})
const optimizationSuggestions = ref([])
const recommendedRoles = ref([])

// 分析角色数据
const analyzeRoles = async () => {
  const allRoles = []
  const flattenRoles = (roles) => {
    roles.forEach(role => {
      allRoles.push(role)
      if (role.children && role.children.length > 0) {
        flattenRoles(role.children)
      }
    })
  }
  flattenRoles(tableData.value)
  
  const coreRoleIds = [888, 9001]
  roleStats.value.totalRoles = allRoles.length
  roleStats.value.coreRoles = allRoles.filter(r => coreRoleIds.includes(r.authorityId)).length
  
  // 检查每个角色的使用情况
  const unusedRoles = []
  for (const role of allRoles) {
    if (!coreRoleIds.includes(role.authorityId)) {
      try {
        const res = await getAuthorityUsageInfo({ authorityId: role.authorityId })
        if (res.code === 0 && res.data.canDelete) {
          unusedRoles.push(role)
        }
      } catch (error) {
        console.error('检查角色失败:', error)
      }
    }
  }
  
  roleStats.value.unusedRoles = unusedRoles.length
  
  // 生成优化建议
  optimizationSuggestions.value = []
  
  // 建议1: 清理未使用的角色
  if (unusedRoles.length > 0) {
    optimizationSuggestions.value.push({
      type: 'warning',
      icon: 'Delete',
      title: '清理未使用的角色',
      description: `发现 ${unusedRoles.length} 个未使用的角色可以删除：${unusedRoles.map(r => r.authorityName).join('、')}`,
      action: 'cleanUnused',
      actionType: 'danger',
      actionText: '批量删除',
      data: unusedRoles
    })
  }
  
  // 建议2: 扁平化角色结构
  const hasNestedRoles = allRoles.some(r => r.children && r.children.length > 0)
  if (hasNestedRoles) {
    optimizationSuggestions.value.push({
      type: 'primary',
      icon: 'Promotion',
      title: '扁平化角色结构',
      description: '当前角色存在多层嵌套结构，建议采用扁平化设计，将所有角色都设为根角色的直接子级，简化管理。',
      action: 'flatten',
      actionType: 'primary',
      actionText: '查看扁平化方案',
      data: null
    })
  }
  
  // 建议3: 角色命名规范
  optimizationSuggestions.value.push({
    type: 'success',
    icon: 'Edit',
    title: '角色命名规范建议',
    description: '建议采用统一的角色命名规范：如"教师"、"学员"、"管理员"等，避免使用ID或代号命名。',
    action: null,
    actionType: 'info',
    actionText: '了解更多',
    data: null
  })
  
  // 推荐的角色结构
  recommendedRoles.value = [
    {
      authorityId: 888,
      authorityName: '超级管理员',
      description: '系统最高权限，负责系统配置和用户管理',
      status: 'exists'
    },
    {
      authorityId: 9001,
      authorityName: '教师',
      description: '教学人员，负责课程管理、学员管理、课时记录等',
      status: 'exists'
    },
    {
      authorityId: 9002,
      authorityName: '学员',
      description: '学习人员，可查看自己的课程和课时记录',
      status: allRoles.some(r => r.authorityId === 9002) ? 'exists' : 'suggest'
    },
    {
      authorityId: 9003,
      authorityName: '财务',
      description: '财务人员，负责费用管理和报表统计（如需要）',
      status: allRoles.some(r => r.authorityId === 9003) ? 'exists' : 'suggest'
    }
  ]
}

const showOptimizationDialog = async () => {
  optimizationDialogVisible.value = true
  await analyzeRoles()
}

const handleOptimizationAction = async (action, data) => {
  if (action === 'cleanUnused') {
    ElMessageBox.confirm(
      `确定要批量删除 ${data.length} 个未使用的角色吗？删除后无法恢复！`,
      '⚠️ 批量删除确认',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        distinguishCancelAndClose: true
      }
    ).then(async () => {
      let successCount = 0
      let failCount = 0
      
      for (const role of data) {
        try {
          const res = await deleteAuthority({ authorityId: role.authorityId })
          if (res.code === 0) {
            successCount++
          } else {
            failCount++
          }
        } catch (error) {
          failCount++
        }
      }
      
      ElMessage.success(`批量删除完成：成功 ${successCount} 个，失败 ${failCount} 个`)
      await getTableData()
      await analyzeRoles()
    }).catch(() => {
      // 用户取消
    })
  } else if (action === 'flatten') {
    ElMessage.info('扁平化功能需要手动调整：将子角色的父级都改为0（根角色）')
  }
}

const exportOptimizationReport = () => {
  const report = {
    生成时间: new Date().toLocaleString(),
    角色统计: roleStats.value,
    优化建议: optimizationSuggestions.value.map(s => ({
      标题: s.title,
      描述: s.description
    })),
    推荐角色结构: recommendedRoles.value
  }
  
  const blob = new Blob([JSON.stringify(report, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `角色优化报告_${Date.now()}.json`
  a.click()
  URL.revokeObjectURL(url)
  
  ElMessage.success('优化报告已导出')
}
const setAuthorityOptions = (AuthorityData, optionsData, disabled) => {
  form.value.authorityId = String(form.value.authorityId)
  AuthorityData &&
        AuthorityData.forEach(item => {
          if (item.children && item.children.length) {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId,
              children: []
            }
            setAuthorityOptions(
              item.children,
              option.children,
              disabled || item.authorityId === form.value.authorityId
            )
            optionsData.push(option)
          } else {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              disabled: disabled || item.authorityId === form.value.authorityId
            }
            optionsData.push(option)
          }
        })
}
// 增加角色
const addAuthority = (parentId) => {
  initForm()
  dialogTitle.value = '新增角色'
  dialogType.value = 'add'
  form.value.parentId = parentId
  setOptions()
  dialogFormVisible.value = true
}
// 编辑角色
const editAuthority = (row) => {
  setOptions()
  dialogTitle.value = '编辑角色'
  dialogType.value = 'edit'
  for (const key in form.value) {
    form.value[key] = row[key]
  }
  setOptions()
  dialogFormVisible.value = true
}

</script>

<script>

export default {
  name: 'Authority'
}
</script>

<style lang="scss">
.authority {
  .el-input-number {
    margin-left: 15px;
    span {
      display: none;
    }
  }
}
.tree-content{
  margin-top: 10px;
  height: calc(100vh - 148px);
  overflow: auto;
}

</style>
