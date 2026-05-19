<template>
  <div class="user-manage">
    <div class="page-header">
      <h1>用户管理</h1>
    </div>
    
    <div class="card">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="openid" label="OpenID" width="180" show-overflow-tooltip />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="avatar_url" label="头像" width="80">
          <template #default="{ row }">
            <el-avatar v-if="row.avatar_url" :size="40" :src="row.avatar_url" />
            <el-avatar v-else :size="40">{{ (row.nickname || 'U')[0] }}</el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column label="已购题库" width="120">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleViewBanks(row)">
              {{ row.bankCount || 0 }} 个题库
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleToggleStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="loadData"
        />
      </div>
    </div>

    <el-dialog v-model="bankDialogVisible" title="用户已购题库" width="500px">
      <el-empty v-if="!userBanks.length" description="该用户暂未购买任何题库" />
      <div v-else class="bank-list">
        <div v-for="bank in userBanks" :key="bank.id" class="bank-item">
          <span class="bank-name">{{ bank.name }}</span>
          <span class="bank-price">¥{{ bank.price }}</span>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getUsers, updateUserStatus, getUserDetail } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const bankDialogVisible = ref(false)
const userBanks = ref([])

const loadData = async () => {
  loading.value = true
  try {
    const data = await getUsers({ page: page.value, page_size: pageSize.value })
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleViewBanks = async (row) => {
  try {
    const data = await getUserDetail(row.id)
    userBanks.value = data.banks || []
    bankDialogVisible.value = true
  } catch (e) {
    console.error(e)
  }
}

const handleToggleStatus = async (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'
  try {
    await ElMessageBox.confirm(`确定${action}用户"${row.nickname || row.openid}"吗？`, '提示', { type: 'warning' })
    await updateUserStatus(row.id, { status: newStatus })
    ElMessage.success(`${action}成功`)
    loadData()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const formatTime = (time) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'

onMounted(loadData)
</script>

<style lang="scss" scoped>
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.bank-list {
  .bank-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid #eee;
    
    &:last-child {
      border-bottom: none;
    }
    
    .bank-name {
      font-weight: 500;
    }
    
    .bank-price {
      color: #f59e0b;
      font-weight: 600;
    }
  }
}
</style>
