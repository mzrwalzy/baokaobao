<template>
  <div class="admin-user-container">
    <div class="page-header">
      <h2>管理员管理</h2>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon> 新增管理员
      </el-button>
    </div>

    <el-table :data="adminList" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.role === 'superadmin'" type="danger">超级管理员</el-tag>
          <el-tag v-else type="info">运营人员</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login" label="最后登录" width="180">
        <template #default="{ row }">
          {{ row.last_login ? formatDate(row.last_login) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="320" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEditDialog(row)">编辑</el-button>
          <el-button size="small" type="warning" @click="openResetDialog(row)">重置密码</el-button>
          <el-button size="small" type="primary" @click="openBankDialog(row)">题库权限</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50]"
      layout="total, sizes, prev, pager, next"
      @change="fetchList"
      class="pagination"
    />

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="formVisible" :title="isEdit ? '编辑管理员' : '新增管理员'" width="500px">
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="超级管理员" value="superadmin" />
            <el-option label="运营人员" value="operator" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">正常</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码弹窗 -->
    <el-dialog v-model="resetVisible" title="重置密码" width="400px">
      <el-form :model="resetForm" :rules="resetRules" ref="resetRef" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="resetForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReset" :loading="resetting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 题库权限弹窗 -->
    <el-dialog v-model="bankVisible" title="分配题库权限" width="500px">
      <el-form label-width="0">
        <el-form-item>
          <el-transfer
            v-model="selectedBanks"
            :data="bankOptions"
            :titles="['所有题库', '已分配']"
            filterable
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="bankVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBankPermission" :loading="bankSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getAdminUsers,
  createAdminUser,
  updateAdminUser,
  resetAdminPassword,
  deleteAdminUser,
  getAdminUserBanks,
  grantAdminBankAccess,
  revokeAdminBankAccess,
  getQuestionBanks
} from '@/api'

const loading = ref(false)
const adminList = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const formVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()
const submitting = ref(false)
const form = reactive({
  id: null,
  username: '',
  password: '',
  nickname: '',
  role: 'operator',
  status: 1
})
const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const resetVisible = ref(false)
const resetRef = ref()
const resetting = ref(false)
const resetForm = reactive({
  id: null,
  password: ''
})
const resetRules = {
  password: [{ required: true, message: '请输入新密码', trigger: 'blur' }]
}

const bankVisible = ref(false)
const bankSubmitting = ref(false)
const currentAdmin = ref(null)
const bankOptions = ref([])
const selectedBanks = ref([])
const currentBanks = ref([])

const formatDate = (dateStr) => {
  const d = new Date(dateStr)
  return d.toLocaleString('zh-CN')
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getAdminUsers({ page: page.value, page_size: pageSize.value })
    adminList.value = res.list || []
    total.value = res.total || 0
  } catch (e) {
    ElMessage.error(e.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  isEdit.value = false
  Object.assign(form, { id: null, username: '', password: '', nickname: '', role: 'operator', status: 1 })
  formVisible.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  Object.assign(form, { id: row.id, username: row.username, password: '', nickname: row.nickname || '', role: row.role, status: row.status })
  formVisible.value = true
}

const submitForm = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateAdminUser(form.id, { nickname: form.nickname, role: form.role, status: form.status })
      ElMessage.success('更新成功')
    } else {
      await createAdminUser({
        username: form.username,
        password: form.password,
        nickname: form.nickname,
        role: form.role
      })
      ElMessage.success('创建成功')
    }
    formVisible.value = false
    fetchList()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const openResetDialog = (row) => {
  resetForm.id = row.id
  resetForm.password = ''
  resetVisible.value = true
}

const submitReset = async () => {
  const valid = await resetRef.value.validate().catch(() => false)
  if (!valid) return

  resetting.value = true
  try {
    await resetAdminPassword(resetForm.id, { password: resetForm.password })
    ElMessage.success('密码重置成功')
    resetVisible.value = false
  } catch (e) {
    ElMessage.error(e.message || '重置失败')
  } finally {
    resetting.value = false
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确认删除该管理员？', '提示', { type: 'warning' })
    await deleteAdminUser(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

const openBankDialog = async (row) => {
  currentAdmin.value = row
  bankVisible.value = true
  bankOptions.value = []
  selectedBanks.value = []
  currentBanks.value = []

  try {
    const [allBanksRes, adminBanksRes] = await Promise.all([
      getQuestionBanks({ page: 1, page_size: 9999 }),
      getAdminUserBanks(row.id)
    ])

    const allBanks = allBanksRes.list || []
    const perms = adminBanksRes.data || adminBanksRes || []

    bankOptions.value = allBanks.map(b => ({
      key: b.id,
      label: b.name,
      disabled: false
    }))

    currentBanks.value = perms.map(p => p.bank_id || p.BankID)
    selectedBanks.value = [...currentBanks.value]
  } catch (e) {
    ElMessage.error(e.message || '加载题库失败')
  }
}

const submitBankPermission = async () => {
  if (!currentAdmin.value) return
  bankSubmitting.value = true
  try {
    const adminID = currentAdmin.value.id
    const toAdd = selectedBanks.value.filter(id => !currentBanks.value.includes(id))
    const toRemove = currentBanks.value.filter(id => !selectedBanks.value.includes(id))

    await Promise.all([
      ...toAdd.map(bankID => grantAdminBankAccess(adminID, { bank_id: bankID })),
      ...toRemove.map(bankID => revokeAdminBankAccess(adminID, bankID))
    ])

    ElMessage.success('权限更新成功')
    bankVisible.value = false
  } catch (e) {
    ElMessage.error(e.message || '权限更新失败')
  } finally {
    bankSubmitting.value = false
  }
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.admin-user-container {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    
    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 600;
    }
  }
  
  .pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
