<template>
  <div class="audit-log">
    <div class="page-header">
      <h1>审计日志</h1>
    </div>

    <div class="card">
      <div class="filter-bar">
        <el-input
          v-model="filters.adminName"
          placeholder="操作人"
          clearable
          style="width: 180px"
          @keyup.enter="handleSearch"
        />
        <el-select
          v-model="filters.action"
          placeholder="操作类型"
          clearable
          style="width: 160px"
        >
          <el-option label="删除" value="delete" />
          <el-option label="更新状态" value="update_status" />
          <el-option label="授予权限" value="grant" />
          <el-option label="撤销权限" value="revoke" />
          <el-option label="导入" value="import" />
        </el-select>
        <el-date-picker
          v-model="filters.dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 260px"
        />
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="admin_name" label="操作人" width="140" />
        <el-table-column prop="action" label="操作类型" width="120">
          <template #default="{ row }">
            <el-tag :type="actionTagType(row.action)" size="small">
              {{ actionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源类型" width="160">
          <template #default="{ row }">
            {{ resourceLabel(row.resource) }}
          </template>
        </el-table-column>
        <el-table-column prop="resource_id" label="资源ID" width="100" />
        <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getAuditLogs } from '@/api'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filters = reactive({
  adminName: '',
  action: '',
  dateRange: null
})

const actionMap = {
  delete: '删除',
  update_status: '更新状态',
  grant: '授予权限',
  revoke: '撤销权限',
  import: '导入'
}

const actionTagMap = {
  delete: 'danger',
  update_status: 'warning',
  grant: 'success',
  revoke: 'info',
  import: 'primary'
}

const resourceMap = {
  question_bank: '题库',
  question: '题目',
  user: '用户',
  user_bank_access: '用户题库权限'
}

function actionLabel(action) {
  return actionMap[action] || action
}

function actionTagType(action) {
  return actionTagMap[action] || 'info'
}

function resourceLabel(resource) {
  return resourceMap[resource] || resource
}

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filters.adminName) params.admin_name = filters.adminName
    if (filters.action) params.action = filters.action
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    const res = await getAuditLogs(params)
    list.value = res.list || []
    total.value = res.total || 0
  } catch {
    // error handled by request interceptor
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handleReset() {
  filters.adminName = ''
  filters.action = ''
  filters.dateRange = null
  page.value = 1
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.audit-log {
  .page-header {
    margin-bottom: 24px;

    h1 {
      font-size: 20px;
      font-weight: 600;
      color: #1f2937;
      margin: 0;
    }
  }

  .card {
    background: #fff;
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  }

  .filter-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
